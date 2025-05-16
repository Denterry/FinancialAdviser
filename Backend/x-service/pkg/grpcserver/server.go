package grpcserver

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	_defaultAddr     = "0.0.0.0:9090"
	_shutdownTimeout = 5 * time.Second
)

// Server wraps a gRPC.Server plus lifecycle channels.
type Server struct {
	address string
	opts    []grpc.ServerOption

	tlsCert, tlsKey string

	grpcServer *grpc.Server
	notify     chan error
}

// New builds a Server configured by opts.
func New(opts ...Option) *Server {
	s := &Server{
		address: _defaultAddr,
		notify:  make(chan error, 1),
	}

	// apply options
	for _, o := range opts {
		o(s)
	}

	return s
}

// Serve registers svc via register and then begins listening
func (s *Server) Serve(register func(*grpc.Server)) {
	// set up TLS if requested
	if s.tlsCert != "" && s.tlsKey != "" {
		creds := credentials.NewServerTLSFromCert(&tls.Certificate{
			Certificate: [][]byte{}, // loaded below
		})
		// load actual cert
		cert, err := tls.LoadX509KeyPair(s.tlsCert, s.tlsKey)
		if err != nil {
			s.notify <- fmt.Errorf("tls.LoadX509KeyPair: %w", err)
			close(s.notify)
			return
		}
		creds = credentials.NewServerTLSFromCert(&cert)
		s.opts = append(s.opts, grpc.Creds(creds))
	}

	// construct grpc server
	s.grpcServer = grpc.NewServer(s.opts...)

	// register user‐provided services
	register(s.grpcServer)

	// listen
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Printf("net.Listen(%s): %w", s.address, err)
		s.notify <- fmt.Errorf("net.Listen(%s): %w", s.address, err)
		close(s.notify)
		return
	}

	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil {
			log.Printf("Serve: %w", err)
			s.notify <- fmt.Errorf("Serve: %w", err)
		}
		close(s.notify)
	}()
}

// Notify returns a channel that will receive any Serve errors
func (s *Server) Notify() <-chan error {
	return s.notify
}

// GracefulStop stops the server, waiting up to timeout for inflight RPCs
func (s *Server) GracefulStop(timeout time.Duration) {
	ch := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(ch)
	}()
	select {
	case <-ch:
	case <-time.After(timeout):
		s.grpcServer.Stop()
	}
}

// GRPC returns the underlying gRPC server
func (s *Server) GRPC() *grpc.Server {
	return s.grpcServer
}
