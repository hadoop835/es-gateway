package proxy

import (
	"encoding/binary"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/cloudwego/netpoll/mux"
	"github.com/es-gateway/pkg/apiserver/config"
	"github.com/es-gateway/pkg/log"
	"github.com/es-gateway/pkg/proxy/codec"
	"github.com/es-gateway/pkg/proxy/options"
	"github.com/pingcap/errors"
	"github.com/pingcap/parser/mysql"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"math/rand"
	"net"
	"sync"
	"time"
	"unsafe"
)

const defaultCapability = mysql.ClientLongPassword | mysql.ClientLongFlag |
	mysql.ClientConnectWithDB | mysql.ClientProtocol41 |
	mysql.ClientTransactions | mysql.ClientSecureConnection | mysql.ClientFoundRows |
	mysql.ClientMultiStatements | mysql.ClientMultiResults | mysql.ClientLocalFiles |
	mysql.ClientConnectAtts | mysql.ClientPluginAuth | mysql.ClientInteractive

// Server is the MySQL protocol server
type Server struct {
	cfg              *config.Config
	tlsConfig        unsafe.Pointer // *tls.Config
	listener         net.Listener
	clients          map[uint64]*clientConn
	rwlock           sync.RWMutex
	capability       uint32
	inShutdownMode   bool
	sessionMapMutex  sync.Mutex
	internalSessions map[interface{}]struct{}
}

func NewProxy(options *options.ProxyRunConfig) (*Server, error) {
	s := &Server{
		clients:          make(map[uint64]*clientConn),
		internalSessions: make(map[interface{}]struct{}, 100),
	}
	s.capability = defaultCapability
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", options.BindAddress, options.InsecurePort))
	if err != nil {
		return nil, errors.Trace(err)
	}
	s.listener = listener

	rand.Seed(time.Now().UTC().UnixNano())

	return s, nil
}

//
func (s *Server) Run() error {
	errChan := make(chan error)
	go s.startNetworkListener(s.listener, false, errChan)
	err := <-errChan
	if err != nil {
		return err
	}
	return <-errChan
}

func (s *Server) startNetworkListener(listener net.Listener, isUnixSocket bool, errChan chan error) {
	if listener == nil {
		errChan <- nil
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Err.Error() == "use of closed network connection" {
					if s.inShutdownMode {
						errChan <- nil
					} else {
						errChan <- err
					}
					return
				}
			}
			log.Log().Error("accept failed", zap.Error(err))
			errChan <- err
			return
		}

		clientConn := s.newConn(conn)

		if err != nil {
			continue
		}

		go s.onConn(clientConn)
	}
}

func (s *Server) onConn(conn *clientConn) {
	ctx := log.WithConnID(context.Background(), conn.connectionID)
	conn.Run(ctx)
}

// ConnectionCount gets current connection count.
func (s *Server) ConnectionCount() int {
	s.rwlock.RLock()
	cnt := len(s.clients)
	s.rwlock.RUnlock()
	return cnt
}

// newConn creates a new *clientConn from a net.Conn.
// It allocates a connection ID and random salt data for authentication.
func (s *Server) newConn(conn net.Conn) *clientConn {
	cc := newClientConn(s)
	if _, ok := conn.(*net.TCPConn); ok {
		//if err := tcpConn.SetKeepAlive(s.cfg.Performance.TCPKeepAlive); err != nil {
		//logutil.BgLogger().Error("failed to set tcp keep alive option", zap.Error(err))
		//}
		//if err := tcpConn.SetNoDelay(s.cfg.Performance.TCPNoDelay); err != nil {
		//logutil.BgLogger().Error("failed to set tcp no delay option", zap.Error(err))
		//}
	}
	//cc.setConn(conn)
	//cc.salt = fastrand.Buf(20)
	return cc
}

func handle(ctx context.Context, conn netpoll.Connection) (err error) {
	log.Log().Info("111", zap.String("", ""))
	mc := ctx.Value(ctxkey).(*svrMuxConn)
	reader := conn.Reader()

	bLen, err := reader.Peek(4)
	if err != nil {
		return err
	}
	length := int(binary.BigEndian.Uint32(bLen)) + 4

	r2, err := reader.Slice(length)
	if err != nil {
		return err
	}

	// handler must use another goroutine
	go func() {
		req := &codec.Message{}
		err = codec.Decode(r2, req)
		if err != nil {
			panic(fmt.Errorf("netpoll decode failed: %s", err.Error()))
		}

		// handler
		resp := req

		// encode
		writer := netpoll.NewLinkBuffer()
		err = codec.Encode(writer, resp)
		if err != nil {
			panic(fmt.Errorf("netpoll encode failed: %s", err.Error()))
		}
		mc.Put(func() (buf netpoll.Writer, isNil bool) {
			return writer, false
		})
	}()
	return nil
}

type connkey struct{}

var ctxkey connkey

func prepare1(connection netpoll.Connection) context.Context {
	return context.Background()
}

func prepare(conn netpoll.Connection) context.Context {
	mc := newSvrMuxConn(conn)
	ctx := context.WithValue(context.Background(), ctxkey, mc)
	return ctx
}

func newSvrMuxConn(conn netpoll.Connection) *svrMuxConn {
	mc := &svrMuxConn{}
	mc.conn = conn
	mc.wqueue = mux.NewShardQueue(mux.ShardSize, conn)
	return mc
}

type svrMuxConn struct {
	conn   netpoll.Connection
	wqueue *mux.ShardQueue // use for write
}

// Put puts the buffer getter back to the queue.
func (c *svrMuxConn) Put(gt mux.WriterGetter) {
	c.wqueue.Add(gt)
}
