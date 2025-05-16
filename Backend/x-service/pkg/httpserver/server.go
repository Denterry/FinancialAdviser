package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	_defaultAddr            = ":8080"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server represents HTTP server
type Server struct {
	engine          *gin.Engine
	server          *http.Server
	notify          chan error
	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
	tlsCert, tlsKey string
	corsConfig      interface{}
}

// New creates new HTTP server instance
func New(opts ...Option) *Server {
	s := &Server{
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	// create Gin engine
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	// set CORS config
	if cfg, ok := s.corsConfig.(cors.Config); ok {
		engine.Use(cors.New(cfg))
	}

	s.engine = engine

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      engine,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	// TLS config hook
	if s.tlsCert != "" && s.tlsKey != "" {
		s.server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return s
}

// Start starts the HTTP server
func (s *Server) Start() {
	go func() {
		var err error
		if s.tlsCert != "" && s.tlsKey != "" {
			err = s.server.ListenAndServeTLS(s.tlsCert, s.tlsKey)
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.notify <- fmt.Errorf("listen error: %w", err)
		}
		close(s.notify)
	}()
}

// Notify returns error channel
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// GetEngine returns the Gin engine instance
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}
