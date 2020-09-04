package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ServerOption func(*Server)

func WithDefaults() ServerOption {
	return func(S *Server) {
		S.HTTP = &http.Server{
			Addr:              "0.0.0.0:8000",
			Handler:           S.NewRouter(),
			ReadHeaderTimeout: time.Duration(10) * time.Second,
			ReadTimeout:       time.Duration(15) * time.Second,
			WriteTimeout:      time.Duration(15) * time.Second,
			IdleTimeout:       time.Duration(15) * time.Second,
		}
	}
}

func WithHTTPServerAddress(ip string, port int) ServerOption {
	return func(S *Server) {
		S.HTTP.Addr = fmt.Sprintf("%s:%d", ip, port)
	}
}

func WithHTTPServerHandler(handler http.Handler) ServerOption {
	return func(S *Server) {
		S.HTTP.Handler = handler
	}
}

func WithHTTPServerReadHeaderTimeout(t time.Duration) ServerOption {
	return func(S *Server) {
		S.HTTP.ReadHeaderTimeout = t
	}
}

func WithHTTPServerReadTimeout(t time.Duration) ServerOption {
	return func(S *Server) {
		S.HTTP.ReadTimeout = t
	}
}

func WithHTTPServerWriteTimeout(t time.Duration) ServerOption {
	return func(S *Server) {
		S.HTTP.WriteTimeout = t
	}
}

func WithHTTPServerIdleTimeout(t time.Duration) ServerOption {
	return func(S *Server) {
		S.HTTP.IdleTimeout = t
	}
}

type Server struct {
	HTTP *http.Server
	DB   *sql.DB
}

func (S *Server) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	runChan := make(chan os.Signal, 1)

	signal.Notify(runChan, os.Interrupt)

	log.Printf("Running server on %s\n", S.HTTP.Addr)
	go func() {
		if err := S.HTTP.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-runChan

	log.Printf("\rShutting server down due to %s\n", sig.String())
	if err := S.HTTP.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func NewWithOptions(opts ...ServerOption) *Server {
	S := new(Server)

	// Establish the defaults
	S.HTTP = &http.Server{
		Addr:              "0.0.0.0:8000",
		Handler:           S.NewRouter(),
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		ReadTimeout:       time.Duration(15) * time.Second,
		WriteTimeout:      time.Duration(15) * time.Second,
		IdleTimeout:       time.Duration(15) * time.Second,
	}

	for _, opt := range opts {
		opt(S)
	}

	return S
}
