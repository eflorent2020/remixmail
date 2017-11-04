package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	// public
	r.HandleFunc("/api/entreprise", getEntrepriseInfo).Methods(http.MethodGet)

	// Warning: order does matter

	// authentified by validation key
	r.HandleFunc("/api/alias/validate/{validationKey}", validateAlias).Methods(http.MethodGet)
	r.HandleFunc("/api/alias/validate/{validationKey}", updateAlias).Methods(http.MethodPut)
	r.HandleFunc("/api/alias/validate/{validationKey}", deleteAlias).Methods(http.MethodDelete)

	// authentified by api key
	r.HandleFunc("/api/alias/{email}/{fullname}", extPutAlias).Methods(http.MethodPut)
	r.HandleFunc("/api/alias/{email}", extDeleteAlias).Methods(http.MethodDelete)

	// authentified by admin
	r.HandleFunc("/api/keys/{email}", handleOption).Methods(http.MethodOptions)
	r.HandleFunc("/api/keys", listApiKey).Methods(http.MethodGet)
	r.HandleFunc("/api/keys/{email}", putApiKey).Methods(http.MethodPut)
	r.HandleFunc("/api/keys/{email}", deleteApiKey).Methods(http.MethodDelete)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./remixmail/dist/")))

	return r
}
