package codec

import (
	"encoding/binary"

	"github.com/cloudwego/netpoll"
)

type Message struct {
	Message string
}

// Encode .
func Encode(writer netpoll.Writer, msg *Message) (err error) {
	header, _ := writer.Malloc(4)
	binary.BigEndian.PutUint32(header, uint32(len(msg.Message)))

	writer.WriteString(msg.Message)
	err = writer.Flush()
	return err
}

// Decode .
func Decode(reader netpoll.Reader, msg *Message) (err error) {
	bLen, err := reader.Next(4)
	if err != nil {
		return err
	}
	l := int(binary.BigEndian.Uint32(bLen))

	msg.Message, err = reader.ReadString(l)
	if err != nil {
		return err
	}
	err = reader.Release()
	return err
}
