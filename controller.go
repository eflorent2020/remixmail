package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

// simple controller to get sitename, tagline and service email
func getEntrepriseInfo(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	/* if need seed for dev
	tr := getTranslaterFromReq(r)
	dsPutAliasSendValidationLink(ctx, tr, "me@privacy.net", "plop")
	*/
	if appengine.IsDevAppServer() {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	response := make(map[string]string)
	response["APPNAME"] = APP_NAME
	response["SERVICE_MAIL"] = SERVICE_MAIL
	response["TAGLINE"] = TAGLINE
	loginUrl, _ := user.LoginURL(ctx, "#/admin")
	logoutUrl, _ := user.LogoutURL(ctx, "/")
	u := user.Current(ctx)
	loggedIn := (u != nil)
	response["LOGGED"] = strconv.FormatBool(loggedIn)
	response["LOGIN"] = loginUrl
	response["LOGOUT"] = logoutUrl
	respondWithJSON(w, http.StatusOK, response)
}

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

/*
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
*/
func putAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	email := strings.TrimSpace(vars["email"])
	fullname := strings.TrimSpace(vars["fullname"])
	alias := new(Alias)
	T := getTranslaterFromReq(r)
	alias, err := dsPutAliasSendValidationLink(ctx, T, email, fullname)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, alias)
}

func updateAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	q := datastore.NewQuery("Alias").Filter("validation_key = ", strings.TrimSpace(vars["validationKey"]))
	var aliases []Alias
	keys, err := q.GetAll(ctx, &aliases)
	if appengine.IsDevAppServer() {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(aliases) == 1 {
		decoder := json.NewDecoder(r.Body)
		var t Alias
		err := decoder.Decode(&t)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(t.Fullname)
		aliases[0].Fullname = t.Fullname
		datastore.Put(ctx, keys[0], &aliases[0])
		respondWithJSON(w, http.StatusOK, aliases[0])
		return
	}
	respondWithError(w, http.StatusInternalServerError, "something went wrong")
}

func deleteAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	q := datastore.NewQuery("Alias").Filter("validation_key = ", strings.TrimSpace(vars["validationKey"]))
	var aliases []Alias
	keys, err := q.GetAll(ctx, &aliases)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(aliases) == 1 {
		decoder := json.NewDecoder(r.Body)
		var t Alias
		err := decoder.Decode(&t)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(t.Fullname)
		aliases[0].Fullname = t.Fullname
		datastore.Delete(ctx, keys[0])
		respondWithJSON(w, http.StatusOK, aliases[0])
		return
	}
	respondWithError(w, http.StatusInternalServerError, "something went wrong")

}

func validateAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	if appengine.IsDevAppServer() {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	validationKey := strings.TrimSpace(vars["validationKey"])
	alias, err := dsValidateAlias(ctx, validationKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	alias.Domain = DOMAIN
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

func listApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		respondWithError(w, http.StatusForbidden, "")
		return
	}
	if !u.Admin {
		respondWithError(w, http.StatusForbidden, "")
		return
	}
}
