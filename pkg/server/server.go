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
			Addr:         "0.0.0.0:8000",
			Handler:      S.NewRouter(),
			ReadTimeout:  time.Duration(15) * time.Second,
			WriteTimeout: time.Duration(15) * time.Second,
			IdleTimeout:  time.Duration(15) * time.Second,
		}
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

	fmt.Printf("Running server on %s\n", S.HTTP.Addr)
	go func() {
		if err := S.HTTP.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-runChan

	fmt.Printf("Shutting server down due to %s", sig.String())
	if err := S.HTTP.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func NewWithOptions(opts ...ServerOption) *Server {
	S := new(Server)

	for _, opt := range opts {
		opt(S)
	}

	return S
}
