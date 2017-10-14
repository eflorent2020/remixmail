package main

const APP_NAME = "RemixMail"
const APP_ID = "snapmail-182207"

// All email addresses on the Email API Authorized Senders list
// need to be valid Gmail or Google-hosted domain accounts
// https://cloud.google.com/appengine/docs/standard/go/mail/#who_can_send_mail
// Note: Even if your app is deployed on a custom domain,
// you can't receive email sent to addresses in that domain.
const DOMAIN = APP_ID + ".appspotmail.com"

const SENDER = "service@" + DOMAIN
const APP_ROOT_URL = "https://" + APP_ID + ".appengine.com"

const SERVICEMAIL = "system@" + DOMAIN
