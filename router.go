package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/entreprise", getEntrepriseInfo).Methods(http.MethodGet)
	r.HandleFunc("/api/alias/validate/{validationKey}", validateAlias).Methods(http.MethodGet)
	r.HandleFunc("/api/alias/validate/{validationKey}", updateAlias).Methods(http.MethodPut)
	r.HandleFunc("/api/alias/validate/{validationKey}", deleteAlias).Methods(http.MethodDelete)
	r.HandleFunc("/api/alias/{email}/{fullname}", putAlias).Methods(http.MethodPut)
	// r.HandleFunc("/api/key/{email}", putApiKey).Methods(http.MethodGet)
	r.HandleFunc("/api/keys", listApiKey).Methods(http.MethodGet)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./remixmail/dist/")))
	return r
}
