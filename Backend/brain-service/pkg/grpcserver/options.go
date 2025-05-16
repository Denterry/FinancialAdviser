package grpcserver

import "google.golang.org/grpc"

// Option configures the gRPC server.
type Option func(*Server)

// Port on which to listen, default ":9090".
func Port(port string) Option {
	return func(s *Server) {
		s.address = port
	}
}

// TLS enables TLS with cert/key. If empty, serves plaintext
func TLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.tlsCert = certFile
		s.tlsKey = keyFile
	}
}

// MaxStreams sets MaxConcurrentStreams
func MaxStreams(n uint32) Option {
	return func(s *Server) {
		s.opts = append(s.opts, grpc.MaxConcurrentStreams(n))
	}
}
