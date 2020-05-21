package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Server provides an http.Server
type Server struct {
	l *zap.Logger
	*http.Server
}

// New creates and configures a server serving all application routes.
//
// The server implements a graceful shutdown and utilizes zap.Logger for logging purposes.
func New(listenAddr string, logger *zap.Logger, mux *chi.Mux) (*Server, error) {
	errorLog, _ := zap.NewStdLogAt(logger, zap.ErrorLevel)
	return NewWithHTTPServer(&http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ErrorLog:     errorLog,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}, logger)
}

// NewWithHTTPServer creates and configures a server serving all application routes.
func NewWithHTTPServer(srv *http.Server, logger *zap.Logger) (*Server, error) {
	return &Server{logger, srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *Server) Start() {
	srv.l.Info("Starting server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			srv.l.Fatal("Could not listen on", zap.String("addr", srv.Addr), zap.Error(err))
		}
	}()
	srv.l.Info("Server is ready to handle requests", zap.String("addr", srv.Addr))
	srv.gracefulShutdown()
}

func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	srv.l.Info("Server is shutting down", zap.String("reason", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		srv.l.Fatal("Could not gracefully shutdown the server", zap.Error(err))
	}
	srv.l.Info("Server stopped")
}
