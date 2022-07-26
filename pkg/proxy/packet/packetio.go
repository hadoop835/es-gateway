package packet

import (
	"bufio"
	"github.com/pingcap/errors"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/terror"
	"io"
	"time"
)

const (
	defaultWriterSize          = 16 * 1024
	DefMaxAllowedPacket uint64 = 67108864
)

// packetIO is a helper to read and write data in packet format.
// MySQL Packets: https://dev.mysql.com/doc/internals/en/mysql-packet.html
type PacketIO struct {
	bufReadConn *BufferedReadConn
	bufWriter   *bufio.Writer
	sequence    uint8
	readTimeout time.Duration
	// maxAllowedPacket is the maximum size of one packet in readPacket.
	maxAllowedPacket uint64
	// accumulatedLength count the length of totally received 'payload' in readPacket.
	accumulatedLength uint64
}

func NewPacketIO(bufReadConn *BufferedReadConn) *PacketIO {
	p := &PacketIO{sequence: 0}
	p.SetBufferedReadConn(bufReadConn)
	p.setMaxAllowedPacket(DefMaxAllowedPacket)
	return p
}

func (p *PacketIO) SetBufferedReadConn(bufReadConn *BufferedReadConn) {
	p.bufReadConn = bufReadConn
	p.bufWriter = bufio.NewWriterSize(bufReadConn, defaultWriterSize)
}

func (p *PacketIO) setReadTimeout(timeout time.Duration) {
	p.readTimeout = timeout
}

func (p *PacketIO) readOnePacket() ([]byte, error) {
	var header [4]byte
	if p.readTimeout > 0 {
		if err := p.bufReadConn.SetReadDeadline(time.Now().Add(p.readTimeout)); err != nil {
			return nil, err
		}
	}
	if _, err := io.ReadFull(p.bufReadConn, header[:]); err != nil {
		return nil, errors.Trace(err)
	}

	sequence := header[3]
	if sequence != p.sequence {
		return nil, errors.Errorf("invalid sequence %d != %d", sequence, p.sequence)
	}

	p.sequence++

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)

	// Accumulated payload length exceeds the limit.
	if p.accumulatedLength += uint64(length); p.accumulatedLength > p.maxAllowedPacket {
		return nil, errors.Errorf("err net Packet too large")
	}

	data := make([]byte, length)
	if p.readTimeout > 0 {
		if err := p.bufReadConn.SetReadDeadline(time.Now().Add(p.readTimeout)); err != nil {
			return nil, err
		}
	}
	if _, err := io.ReadFull(p.bufReadConn, data); err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

func (p *PacketIO) setMaxAllowedPacket(maxAllowedPacket uint64) {
	p.maxAllowedPacket = maxAllowedPacket
}

func (p *PacketIO) readPacket() ([]byte, error) {
	p.accumulatedLength = 0
	if p.readTimeout == 0 {
		if err := p.bufReadConn.SetReadDeadline(time.Time{}); err != nil {
			return nil, errors.Trace(err)
		}
	}
	data, err := p.readOnePacket()
	if err != nil {
		return nil, errors.Trace(err)
	}

	if len(data) < mysql.MaxPayloadLen {

		return data, nil
	}

	// handle multi-packet
	for {
		buf, err := p.readOnePacket()
		if err != nil {
			return nil, errors.Trace(err)
		}

		data = append(data, buf...)

		if len(buf) < mysql.MaxPayloadLen {
			break
		}
	}

	return data, nil
}

// writePacket writes data that already have header
func (p *PacketIO) writePacket(data []byte) error {
	length := len(data) - 4
	for length >= mysql.MaxPayloadLen {
		data[3] = p.sequence
		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff

		if n, err := p.bufWriter.Write(data[:4+mysql.MaxPayloadLen]); err != nil {
			return errors.Trace(mysql.ErrBadConn)
		} else if n != (4 + mysql.MaxPayloadLen) {
			return errors.Trace(mysql.ErrBadConn)
		} else {
			p.sequence++
			length -= mysql.MaxPayloadLen
			data = data[mysql.MaxPayloadLen:]
		}
	}
	data[3] = p.sequence
	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)

	if n, err := p.bufWriter.Write(data); err != nil {
		terror.Log(errors.Trace(err))
		return errors.Trace(mysql.ErrBadConn)
	} else if n != len(data) {
		return errors.Trace(mysql.ErrBadConn)
	} else {
		p.sequence++
		return nil
	}
}

func (p *PacketIO) flush() error {
	err := p.bufWriter.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return err
}
