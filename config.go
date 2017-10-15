package main

const (
	// used in home page and mail signature

	APP_NAME = "RemixMail"

	// appengine hosting
	//

	APP_ID = "snapmail-182207"
	// DOMAIN = APP_ID + ".appspot.com"
	// base url for building links such as validation link
	APP_ROOT_URL = "https://" + APP_ID + ".appspot.com"

	MAIL_DOMAIN = APP_ID + ".appspotmail.com"

	// the mail from all our messages
	// https://cloud.google.com/appengine/docs/standard/go/mail/#who_can_send_mail
	// Many things are possible ex: All email addresses on the Email API Authorized
	// Senders list need to be valid Gmail or **Google-hosted domain accounts**
	SENDER = "service@" + MAIL_DOMAIN

	// the mail where to send registration requests
	SERVICE_MAIL = "system@" + MAIL_DOMAIN

	// tagline shown in homepage and confirmation mail
	TAGLINE = "An email-address proxy service "

	DEFAULT_LANG = "en-US"
)

// see main, init() :
// 	i18n.MustLoadTranslationFile("lang/en-us.all.json")
var LANGS = [2]string{"en-US", "fr-FR"}
