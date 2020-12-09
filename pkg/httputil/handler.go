package httputil

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	Register(router *mux.Router)
	Handle(w http.ResponseWriter, r *http.Request)
}
