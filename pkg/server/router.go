package server

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewRouter is a function that will generate a new mux.Router instance
// with appropriate routes injected.
//
// Arguments:
//     None
//
// Returns:
//     (*mux.Router): A pointer to the new instance of mux.Router
func (S *Server) NewRouter() *mux.Router {
	R := mux.NewRouter()

	R.Handle("/webhook", handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(S.NewWebhookHandler)))

	return R
}
