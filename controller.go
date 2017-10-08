package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func getTranslaterFromReq(r *http.Request) i18n.TranslateFunc {
	acceptLang := r.Header.Get("Accept-Language")
	defaultLang := "en-US" // known valid language
	T, err := i18n.Tfunc(acceptLang, acceptLang, defaultLang)
	if err != nil {
		println("something went wring with i18n")
	}
	return T
}

func getTranslater(acceptLang string) i18n.TranslateFunc {
	defaultLang := "en-US" // known valid language
	T, err := i18n.Tfunc(acceptLang, acceptLang, defaultLang)
	if err != nil {
		println("something went wring with i18n")
	}
	return T
}

func putApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)

	if u == nil {
		respondWithError(w, http.StatusForbidden, "must be logged in")
		return
	}
	if u.Admin == true {
		log.Errorf(ctx, "invalid acl for "+u.Email)
		respondWithError(w, http.StatusForbidden, "must be logged in")
		return
	}

	vars := mux.Vars(r)
	email := strings.TrimSpace(vars["email"])
	T := getTranslaterFromReq(r)
	_, err := dsPutAPiKey(ctx, T, email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "mail sent")
}

func putAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	email := strings.TrimSpace(vars["email"])
	fullname := strings.TrimSpace(vars["fullname"])
	alias := new(Alias)
	T := getTranslaterFromReq(r)
	alias, err := dsPutAlias(ctx, T, email, fullname)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, alias)
}

func validateAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	validationKey := strings.TrimSpace(vars["validationKey"])
	alias, err := dsValidateAlias(ctx, validationKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, alias)
}

func getAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	params := mux.Vars(r)
	email := params["email"]
	aliases, err := dsGetAliases(ctx, email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(aliases) > 0 {
		respondWithJSON(w, http.StatusFound, aliases)
	} else {
		respondWithError(w, http.StatusFound, "Unknown client use PUT to register")
	}
}

func deleteAliases(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	err := dsDeleteAliases(r, email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusFound, email+"aliases deleted")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
