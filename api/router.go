package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func New() http.Handler {
	r := mux.NewRouter().PathPrefix("/accounts").Subrouter()
	r.StrictSlash(true)
	r.HandleFunc("/create", create).Methods("PUT")
	r.HandleFunc("/{name}", delete).Methods("DELETE")
	r.HandleFunc("/{name}/token", updateToken).Methods("PATCH").Queries("token", "")
	r.HandleFunc("/{name}/token", getToken).Methods("GET")
	return RecoverMiddleware(r)
}