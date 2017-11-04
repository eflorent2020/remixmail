package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

// simple controller to get sitename, tagline and service email
func getEntrepriseInfo(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	addCorsIfNeeded(w)
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

// take the client language and return our i18n or default
func getLang(lang string) string {
	if lang == "" {
		lang = "en"
	}
	for _, a := range LANGS {
		if a == lang {
			return lang
		}
	}
	// try to seek fr-BE for fr-FR ...
	for _, a := range LANGS {
		if a[0:2] == lang[0:2] {
			return lang
		}
	}
	return DEFAULT_LANG
}

// take the client request and return the best
// github.com/nicksnyder/go-i18n/i18n translater function
func getTranslaterFromReq(r *http.Request) i18n.TranslateFunc {
	acceptLang := r.Header.Get("Accept-Language")
	return getTranslater(acceptLang)
}

// take the client language and return the best
// github.com/nicksnyder/go-i18n/i18n translater function
func getTranslater(acceptLang string) i18n.TranslateFunc {
	T, err := i18n.Tfunc(getLang(acceptLang), DEFAULT_LANG, DEFAULT_LANG)
	if err != nil {
		println(err.Error())
	}
	return T
}

// update alias by http, secured by validation key,
// at confirmation mail user should update it's Fullname
// as shown in mail
func updateAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	q := datastore.NewQuery("Alias").Filter("validation_key = ", strings.TrimSpace(vars["validationKey"]))
	var aliases []Alias
	keys, err := q.GetAll(ctx, &aliases)
	addCorsIfNeeded(w)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(aliases) >= 1 {
		decoder := json.NewDecoder(r.Body)
		var t Alias
		err := decoder.Decode(&t)
		if err != nil {
			println(err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		// update fullname here maybe we'll also handle PGP here one day
		aliases[0].Fullname = t.Fullname
		aliases[0].PGPPubKey = t.PGPPubKey
		datastore.Put(ctx, keys[0], &aliases[0])
		respondWithJSON(w, http.StatusOK, aliases[0])
		return
	}
	respondWithError(w, http.StatusInternalServerError, "something went wrong")
}

// on mail feedback user may delete it's account,
// request secured by validation key
func deleteAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	q := datastore.NewQuery("Alias").Filter("validation_key = ", strings.TrimSpace(vars["validationKey"]))
	var aliases []Alias
	keys, err := q.GetAll(ctx, &aliases)
	if err != nil {
		println("ERrroooooorr")
		println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(aliases) >= 1 {
		datastore.Delete(ctx, keys[0])
		respondWithJSON(w, http.StatusOK, aliases[0])
		return
	}
	respondWithError(w, http.StatusInternalServerError, "something went wrong")
}

// called when user have received email, secured y validation key
// call the datastore func who pass the validation key to true
func validateAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	addCorsIfNeeded(w)
	validationKey := strings.TrimSpace(vars["validationKey"])
	alias, err := dsValidateAlias(ctx, validationKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, alias)
}

// list all system api keys
func listApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := handleAdminLogin(w, r)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}
	apiKeys, err := listAPiKey(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, apiKeys)
}

func putApiKey(w http.ResponseWriter, r *http.Request) {
	err := handleAdminLogin(w, r)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	email := strings.TrimSpace(vars["email"])
	T := getTranslaterFromReq(r)
	_, err = dsPutAPiKey(ctx, T, email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "mail sent")
}

// delete an api key
func deleteApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	println("ici")
	fmt.Println(vars)
	email := strings.TrimSpace(vars["email"])
	println("deleting " + email)
	err := handleAdminLogin(w, r)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	key := datastore.NewKey(ctx, "ApiKey", email, 0, nil)

	err = datastore.Delete(ctx, key)

	if err == nil {
		respondWithJSON(w, http.StatusOK, "OK")
		return
	}
	respondWithError(w, http.StatusInternalServerError, err.Error())
}

func extDeleteAlias(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	println(vars)
	email := strings.TrimSpace(vars["email"])
	err := checkAuthorisation(r)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
	}

	dsDeleteAliases(r, email)
}

func extPutAlias(w http.ResponseWriter, r *http.Request) {
	err := checkAuthorisation(r)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}
	vars := mux.Vars(r)
	email := strings.TrimSpace(vars["email"])
	fullname := strings.TrimSpace(vars["fullname"])
	ctx := appengine.NewContext(r)
	name, err := GetFreeName(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	alias := new(Alias)
	alias.CreatedDate = time.Now()
	alias.Email = email
	alias.Fullname = fullname
	alias.Validated = true
	alias.Alias = name + "@" + MAIL_DOMAIN
	alias.ValidationKey = uuid.NewV4().String()
	key := datastore.NewIncompleteKey(ctx, "Alias", aliasKey(ctx))
	_, err = datastore.Put(ctx, key, alias)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, alias)
}

func checkAuthorisation(r *http.Request) error {
	ctx := appengine.NewContext(r)
	queryString := r.URL.Path
	header := r.Header.Get("X-HMAC")
	if len(strings.Split(header, ",")) != 2 {
		return errors.New("invalid X-HMAC header")
	}

	clientID := strings.Split(header, ",")[0]
	signature := strings.Split(header, ",")[1]
	key := datastore.NewKey(ctx, "ApiKey", clientID, 0, nil)
	var apiKey ApiKey
	err := datastore.Get(ctx, key, &apiKey)
	if err != nil {
		return errors.New("client_id not found " + clientID)
	}
	expected := hmac.New(sha256.New, []byte(apiKey.ApiKey))
	expected.Write([]byte(strings.TrimSpace(queryString)))
	if signature != hex.EncodeToString(expected.Sum(nil)) {
		return errors.New("invalid signature " + signature + " for " + queryString)
	}
	// No validation errors.  Signature is good.
	return nil

}

// list all system api keys
func handleOption(w http.ResponseWriter, r *http.Request) {
	if appengine.IsDevAppServer() {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		respondWithJSON(w, 200, "{}")
	}
}

// utility function to add CORS header if in dev
func addCorsIfNeeded(w http.ResponseWriter) {
	if appengine.IsDevAppServer() {
		//w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}

// utity function to check if user is an admin
// additionnaly set CORS header for dev
func handleAdminLogin(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	addCorsIfNeeded(w)
	if appengine.IsDevAppServer() {
		return nil
	}
	if u == nil {
		return errors.New("user not logged in")
	}
	if !u.Admin {
		return errors.New("user is not an app admin")
	}
	return nil
}

// utility function to format an API error in JSON format
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// utility function to format an API success response in JSON format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
