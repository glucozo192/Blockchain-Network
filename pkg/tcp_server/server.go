package tcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/blockchain-network/configs"
	"github.com/blockchain-network/internal/models"
)

type TcpServer struct {
	Addr     configs.ConnectionAddr
	listener net.Listener
	Handlers map[models.Event]func(ctx context.Context, req *models.Request) (*models.Response, error)
	Logger   *log.Logger
}

func (s *TcpServer) Init(ctx context.Context) error {
	s.Logger.Printf("tcp listen on port: %v", s.Addr.Port)
	l, err := net.Listen("tcp", s.Addr.GetConnectionString())
	if err != nil {
		return fmt.Errorf("unable to connect tcp address: %s", s.Addr.GetConnectionString())
	}
	s.listener = l
	return nil
}

func (s *TcpServer) Start(ctx context.Context) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.throwError(conn, err)
			continue
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			s.throwError(conn, err)
			continue
		}
		data := buf[:n]
		var (
			req models.Request
		)
		if err := json.Unmarshal(data, &req); err != nil {
			s.throwError(conn, err)
			continue
		}
		res, err := s.Handlers[req.Event](ctx, &req)
		if err != nil {
			s.throwError(conn, err)
			continue
		}
		b, _ := json.Marshal(res)
		if _, err := conn.Write(b); err != nil {
			s.throwError(conn, err)
			continue
		}
		conn.Close()
	}
}

func (s *TcpServer) Stop(ctx context.Context) error {
	return s.listener.Close()
}

func (s *TcpServer) throwError(conn net.Conn, err error) {
	s.Logger.Println(err)
	conn.Close()
}
