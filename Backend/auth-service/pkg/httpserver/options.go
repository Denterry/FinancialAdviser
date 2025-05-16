package httpserver

import (
	"net"
	"time"
)

// Option represents a function that configures the server
type Option func(*Server)

// Port sets the server port
func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

// ReadTimeout sets the read timeout
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WriteTimeout sets the write timeout
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// ShutdownTimeout sets the shutdown timeout
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

// EnableTLS tells the server to serve over HTTPS using the provided cert and key files
func EnableTLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.tlsCert = certFile
		s.tlsKey = keyFile
	}
}

// CORSConfig lets supply a Gin-compatible CORS configuration
func CORSConfig(cfg interface{}) Option {
	return func(s *Server) {
		s.corsConfig = cfg
	}
}
