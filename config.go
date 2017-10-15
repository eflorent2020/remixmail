package main

const (
	// used in home page and mail signature

	APP_NAME = "RemixMail"

	// appengine hosting
	//
	// Many things are possible with domain names see
	// https://cloud.google.com/appengine/docs/standard/go/mail/#who_can_send_mail
	APP_ID = "snapmail-182207"
	DOMAIN = APP_ID + ".appengine.com"

	MAIL_DOMAIN = APP_ID + ".appspotmail.com"

	// the mail from all our messages
	SENDER = "service@" + DOMAIN

	// base url for building links such as validation link
	APP_ROOT_URL = "https://" + APP_ID + ".appengine.com"

	// the mail where to send registration requests
	SERVICE_MAIL = "system@" + MAIL_DOMAIN

	// tagline shown in homepage and confirmation mail
	TAGLINE = " email-address proxy service "

	DEFAULT_LANG = "en-EN"
)

var LANGS = [2]string{"en-EN", "fr-FR"}
