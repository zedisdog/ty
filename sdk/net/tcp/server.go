package tcp

import (
	"context"
	"fmt"
	"log"
	"net"
)

func NewServerWithContext(ctx context.Context, address string, handler func(conn *net.TCPConn)) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		address:     address,
		connHandler: handler,
		ctx:         ctx,
		stopFunc:    cancel,
	}
}

func NewServer(address string, handler func(*net.TCPConn)) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		address:     address,
		connHandler: handler,
		ctx:         ctx,
		stopFunc:    cancel,
	}
}

type Server struct {
	//address tcp监听地址
	address string

	//tcpListener tcp监听器
	tcpListener net.Listener

	//connHandler tcp链接处理器
	connHandler func(*net.TCPConn)

	ctx context.Context

	stopFunc func()
}

func (s *Server) Start() (err error) {
	add, _ := net.ResolveTCPAddr("tcp", s.address)
	s.tcpListener, err = net.ListenTCP("tcp", add)
	if err != nil {
		err = fmt.Errorf("start tcp server failed: %w", err)
		return
	}
	defer func() {
		_ = s.tcpListener.Close()
	}()

	connChan := make(chan net.Conn)
	go func(connChan chan<- net.Conn) {
		for {
			conn, err := s.tcpListener.Accept()
			if err != nil {
				log.Printf("logger: tcp.server, msg: tcp server stoped, error: %s\n", err.Error())
				//err = fmt.Errorf("tcp server stoped: %w", err)
				close(connChan)
				break
			}
			connChan <- conn
		}
	}(connChan)

	for {
		select {
		case conn, ok := <-connChan:
			if !ok {
				return
			}
			go s.connHandler(conn.(*net.TCPConn))
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Server) Stop() {
	s.stopFunc()
}
