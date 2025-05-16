package grpcserver

import "google.golang.org/grpc"

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.address = port
	}
}

// TLS -.
func TLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.tlsCert = certFile
		s.tlsKey = keyFile
	}
}

// MaxStreams -.
func MaxStreams(n uint32) Option {
	return func(s *Server) {
		s.opts = append(s.opts, grpc.MaxConcurrentStreams(n))
	}
}
