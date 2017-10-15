package main

import (
	"net/http"

	"github.com/nicksnyder/go-i18n/i18n"
)

func init() {
	i18n.MustLoadTranslationFile("lang/en-us.all.json")
	i18n.MustLoadTranslationFile("lang/fr-fr.all.json")
	r := makeRouter()
	http.Handle("/", r)
	http.HandleFunc("/_ah/mail/", incomingMail)
}
