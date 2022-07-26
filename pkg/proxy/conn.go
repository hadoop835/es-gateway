package proxy

import (
	"crypto/tls"
	"github.com/es-gateway/pkg/proxy/packet"
	"github.com/pingcap/parser/mysql"
	"golang.org/x/net/context"
	"net"
	"sync"
	"time"
)

// clientConn represents a connection between server and client, it maintains connection specific state,
// handles client query.
type clientConn struct {
	pkt          *packet.PacketIO         // a helper to read and write data in packet format.
	bufReadConn  *packet.BufferedReadConn // a buffered-read net.Conn or buffered-read tls.Conn.
	tlsConn      *tls.Conn                // TLS connection, nil if not TLS.
	server       *Server                  // a reference of server instance.
	capability   uint32                   // client capability affects the way server handles client request.
	connectionID uint64                   // atomically allocated by a global variable, unique in process scope.
	user         string                   // user of the client.
	dbname       string                   // default database name.
	salt         []byte                   // random bytes used for authentication.
	//alloc        arena.Allocator   // an memory allocator for reducing memory allocation.
	//chunkAlloc   chunk.Allocator
	lastPacket []byte // latest sql query string, currently used for logging error.
	// ShowProcess() and mysql.ComChangeUser both visit this field, ShowProcess() read information through
	// the TiDBContext and mysql.ComChangeUser re-create it, so a lock is required here.
	ctx struct {
		sync.RWMutex
		//*TiDBContext // an interface to execute sql statements.
	}
	attrs        map[string]string // attributes parsed from client handshake response, not used for now.
	peerHost     string            // peer host
	peerPort     string            // peer port
	status       int32             // dispatching/reading/shutdown/waitshutdown
	lastCode     uint16            // last error code
	collation    uint8             // collation used by client, may be different from the collation used by database.
	lastActive   time.Time         // last active time
	authPlugin   string            // default authentication plugin
	isUnixSocket bool              // connection is Unix Socket file
	//rsEncoder     *resultEncoder    // rsEncoder is used to encode the string result to different charsets.
	//inputDecoder  *inputDecoder     // inputDecoder is used to decode the different charsets of incoming strings to utf-8.
	socketCredUID uint32 // UID from the other end of the Unix Socket
	// mu is used for cancelling the execution of current transaction.
	mu struct {
		sync.RWMutex
		//cancelFunc context.CancelFunc
	}
}

// newClientConn creates a *clientConn object.
func newClientConn(s *Server) *clientConn {
	return &clientConn{
		server: s,
		//connectionID: s.globalConnID.NextID(),
		collation: mysql.DefaultCollationID,
		//alloc:        arena.NewAllocator(32 * 1024),
		//chunkAlloc:   chunk.NewAllocator(),
		//status:       connStatusDispatching,
		lastActive: time.Now(),
		//authPlugin:   mysql.AuthNativePassword,
	}
}

func (cc *clientConn) handshake(ctx context.Context) error {

	return nil
}

func (cc *clientConn) Run(ctx context.Context) {
	print(ctx)
}

func (cc *clientConn) setConn(conn net.Conn) {
	cc.bufReadConn = packet.NewBufferedReadConn(conn)
	if cc.pkt == nil {
		cc.pkt = packet.NewPacketIO(cc.bufReadConn)
	} else {
		// Preserve current sequence number.
		cc.pkt.SetBufferedReadConn(cc.bufReadConn)
	}
}
