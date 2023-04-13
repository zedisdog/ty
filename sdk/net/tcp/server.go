package tcp

import (
	"bufio"
	"fmt"
	bytex "github.com/zedisdog/ty/bytes"
	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/log"
	"net"
	"sync"
	"time"
)

var errClose = errx.New("close")

func WithPkgSize(size int) func(*Server) {
	return func(server *Server) {
		server.pkgSize = size
	}
}

func WithLogger(logger log.ILog) func(*Server) {
	return func(server *Server) {
		server.logger = logger
	}
}

func NewServer(
	address string,
	head []byte,
	foot []byte,
	onMsg func(msg []byte, index int) (replay []byte, close bool),
	opts ...func(*Server),
) (s *Server) {
	s = &Server{
		address: address,
		head:    head,
		foot:    foot,
		onMsg:   onMsg,

		conns:    make([]*net.TCPConn, 0, 100),
		connLock: new(sync.RWMutex),
		pkgSize:  1024,
	}

	for _, opt := range opts {
		opt(s)
	}

	return
}

type Server struct {
	//address tcp监听地址
	address string
	//tcpListener tcp监听器
	tcpListener net.Listener

	//conns tcp连接
	conns    []*net.TCPConn
	connLock *sync.RWMutex

	logger log.ILog

	head    []byte
	foot    []byte
	pkgSize int

	onFirstMsg func(msg []byte, index int) (replay []byte, close bool)
	onMsg      func(msg []byte, index int) (replay []byte, close bool)
}

func (s *Server) Log(msg string, level log.Level, fields ...*log.Field) {
	if s.logger != nil {
		fmt.Printf("[tcp server] %s, %#v", msg, fields)
	} else {
		s.logger.Log(msg, level, fields...)
	}
}

func (s *Server) Start() (err error) {
	add, err := net.ResolveTCPAddr("tcp", s.address)
	if err != nil {
		return
	}
	s.tcpListener, err = net.ListenTCP("tcp", add)
	if err != nil {
		return
	}
	defer func() {
		_ = s.tcpListener.Close()
	}()

	for {
		var conn net.Conn
		conn, err = s.tcpListener.Accept()
		if err != nil {
			err = errx.Wrap(err, "[tcp server] get conn failed")
			return
		}
		go func() {
			err := s.receiveConn(conn.(*net.TCPConn))
			if err != nil {
				s.Log("receive conn failed", log.Warn, log.NewField("error", err))
			}
		}()
	}
}

func (s *Server) Stop() {
	_ = s.tcpListener.Close()
}

func (s *Server) receiveConn(conn *net.TCPConn) (err error) {
	index := s.findFirstEmptySeat()
	if index == -1 {
		index = len(s.conns)
		s.connLock.Lock()
		s.conns = append(s.conns, conn)
		s.connLock.Unlock()
	} else {
		s.connLock.Lock()
		s.conns[index] = conn
		s.connLock.Unlock()
	}

	scanner := s.newScanner(conn)
	if s.onFirstMsg != nil {
		err = s.process(scanner, index, conn, s.onFirstMsg)
		if err != nil {
			s.connLock.Lock()
			_ = s.conns[index].Close()
			s.conns[index] = nil
			s.connLock.Unlock()
			return
		}
	}

	go s.watch(conn, scanner, index)

	return
}

func (s *Server) watch(conn *net.TCPConn, scanner *bufio.Scanner, index int) {
	for {
		err := s.process(scanner, index, conn, s.onMsg)
		if err != nil {
			s.connLock.Lock()
			_ = s.conns[index].Close()
			s.conns[index] = nil
			s.connLock.Unlock()
			s.Log("close", log.Warn, log.NewField("error", err))
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *Server) process(scanner *bufio.Scanner, index int, conn *net.TCPConn, callback func(msg []byte, index int) ([]byte, bool)) (err error) {
	if !scanner.Scan() {
		err = errx.Wrap(scanner.Err(), "process msg failed")
		return
	}

	replay, clos := callback(scanner.Bytes(), index)
	if replay != nil {
		_, err = conn.Write(replay)
		if err != nil {
			err = errx.Wrap(err, "write msg failed")
			return
		}
	}

	if clos {
		return errClose
	}

	return
}

func (s *Server) newScanner(conn *net.TCPConn) (scanner *bufio.Scanner) {
	buff := bufio.NewReader(conn)
	scanner = bufio.NewScanner(buff)
	scanner.Split(bytex.SplitByHeadAndFoot(s.head, s.foot))
	scanner.Buffer(make([]byte, s.pkgSize), s.pkgSize)
	return
}

func (s *Server) findFirstEmptySeat() int {
	s.connLock.RLock()
	defer s.connLock.RUnlock()

	for index, conn := range s.conns {
		if conn == nil {
			return index
		}
	}

	return -1
}
