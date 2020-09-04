package server

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (S *Server) NewRouter() *mux.Router {
	R := mux.NewRouter()

	R.Handle("/webhook", handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(S.NewWebhookHandler)))

	return R
}
