package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func makeRouter() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/api/apikey/{email}", putApiKey).Methods(http.MethodGet)
	r.HandleFunc("/api/alias/{email}/{fullname}", putAlias).Methods(http.MethodPut)
	r.HandleFunc("/api/alias/validate/{validationKey}", validateAlias).Methods(http.MethodGet)
	return r
}
