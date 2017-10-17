package main

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
)

// MailAlias is the main struc of the app,
// an email may have several Alias with different fullname
// the application must validate any insert by email
type Alias struct {
	Email         string    `datastore:"email"`
	Alias         string    `datastore:"alias"`
	Fullname      string    `datastore:"fullname"`
	CreatedDate   time.Time `datastore:"created_at"`
	Validated     bool      `datastore:"validated"`
	ValidationKey string    `datastore:"validation_key"`
}

// utitlity function to avoid repeat
// and handle namespaces in appengine datastore
func aliasKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Alias", "default_alias", 0, nil)
}

// take an email, a fullname, check email format and put in datastore
func dsPutAliasSendValidationLink(ctx context.Context, lang string, email string, fullname string) (*Alias, error) {
	err := ValidateEmailFormat(email)
	if err != nil {
		return nil, err
	}
	err = ValidateFullnameFormat(fullname)
	if err != nil {
		return nil, err
	}
	aliasExists, err := dsGetAlias(ctx, email, "")
	if aliasExists != nil {
		sendValidationLink(ctx, lang, &aliasExists[0])
		return &aliasExists[0], nil
	}
	name, err := GetFreeName(ctx)
	if err != nil {
		return nil, err
	}
	alias := new(Alias)
	alias.CreatedDate = time.Now()
	alias.Email = email
	alias.Fullname = fullname
	alias.Validated = false
	alias.Alias = name + "@" + MAIL_DOMAIN
	alias.ValidationKey = uuid.NewV4().String()
	key := datastore.NewIncompleteKey(ctx, "Alias", aliasKey(ctx))
	_, err = datastore.Put(ctx, key, alias)
	if err != nil {
		return nil, err
	}
	sendValidationLink(ctx, lang, alias)
	return alias, nil
}

// storeGetAliases take an email as argument and return an array of
// all Alias struc
func dsGetAliases(ctx context.Context, email string) ([]Alias, error) {
	q := datastore.NewQuery("Alias").Filter("email = ", email)
	var aliases []Alias
	if _, err := q.GetAll(ctx, &aliases); err != nil {
		return aliases, err
	}
	return aliases, nil
}

// storeGetAliases take an email as argument and return an array of
// all Alias struc
func dsFindAliased(ctx context.Context, email string) (Alias, error) {
	// aliasKey := strings.Split(email, "@")[0]
	q := datastore.NewQuery("Alias").Filter("alias = ", email)
	var aliases []Alias
	var alias Alias
	if _, err := q.GetAll(ctx, &aliases); err != nil {
		log.Criticalf(ctx, err.Error())
		return alias, err
	}
	if len(aliases) != 1 {
		return alias, errors.New("inconsistent aliased:" + email)
	}
	return aliases[0], nil
}

// storeGetAlias take an email and a fullname as argument
// and return an array  []Alias,nill or nill,error,
// appengine will restrict to maximum 1000 results
func dsGetAlias(ctx context.Context, email string, fullname string) ([]Alias, error) {
	aliases, err := dsGetAliases(ctx, email)
	if err != nil {
		return nil, err
	}
	if fullname != "" {
		return filterByName(aliases, fullname), nil
	} else {
		return aliases, nil
	}
}

// validate an email using a regexp only
func ValidateEmailFormat(email string) error {
	if strings.ContainsAny(email, "<>&") {
		return errors.New("invalid char in string")
	}
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return errors.New("invalid mail format:" + email + "")
	}
	return nil
}

// validate an email using a regexp only
// using Unicode Character Property that matches any kind of letter from any language
func ValidateFullnameFormat(fullname string) error {
	if strings.ContainsAny(fullname, "<>&") {
		return errors.New("invalid char in string")
	}
	// \\p{L} is a Unicode Character Property
	fullnameRegexp := regexp.MustCompile("^[\\PN .'-]+$")
	if !fullnameRegexp.MatchString(fullname) {
		return errors.New("invalid fullname format: \"" + fullname + "\"")
	}
	return nil
}

// This util take an Alias array and filter it using a test function
// exemple: return aliases[] with fullname="Doe"
func filterByName(aliases []Alias, fullname string) (ret []Alias) {
	for _, a := range aliases {
		if a.Fullname == fullname {
			ret = append(ret, a)
		}
	}
	return ret
}

// For a given validationKey (previously sent by mail)
// set the Alias as Valitaed by the owner of the real email
func dsValidateAlias(ctx context.Context, validationKey string) (Alias, error) {
	q := datastore.NewQuery("Alias").Filter("validation_key = ", validationKey)
	var aliases []Alias
	var alias Alias
	keys, err := q.GetAll(ctx, &aliases)
	if err != nil {
		return alias, err
	}
	if len(aliases) == 1 {
		aliases[0].Validated = true
		datastore.Put(ctx, keys[0], &aliases[0])
		return aliases[0], nil
	}
	return alias, errors.New("inconsitency with validation key " + validationKey)
}

// delete all aliases for a customer email
func dsDeleteAliases(r *http.Request, email string) error {
	ctx := appengine.NewContext(r)
	aliases, err := dsGetAliases(ctx, email)
	if err != nil {
		return err
	}
	q := datastore.NewQuery("Alias").Filter("email = ", email)
	keys, err := q.GetAll(ctx, &aliases)
	for _, key := range keys {
		err := datastore.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// delete one alias by appengine key
func dsDeleteAlias(r *http.Request, key *datastore.Key) error {
	ctx := appengine.NewContext(r)
	if err := datastore.Delete(ctx, key); err != nil {
		return err
	}
	return nil
}

// create a confirmation link (url) to be put in email sent
func createConfirmationURL(alias *Alias) string {
	return APP_ROOT_URL + "/#/alias/validate/" + alias.ValidationKey
}

// make some checks, build the templated, translated ,
// validation mail and call send mail function
func sendValidationLink(ctx context.Context, lang string, alias *Alias) (*mail.Message, error) {
	// T, lang := getTranslater(lang)
	addr := alias.Email
	url := createConfirmationURL(alias)
	templateData := struct {
		Name     string
		Url      string
		Alias    string
		APP_NAME string
		TAGLINE  string
	}{
		Name:     alias.Fullname,
		Url:      url,
		Alias:    alias.Alias,
		APP_NAME: APP_NAME,
		TAGLINE:  TAGLINE}
	templateFileName := "templates/" + lang + "/mail-confirm.html"

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Criticalf(ctx, err.Error())
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return nil, err
	}
	body := buf.String()
	T := getTranslater(lang)
	subject := T("confirm_registration")
	msg := &mail.Message{Sender: APP_NAME + " <" + SENDER + ">",
		To:       []string{addr},
		Subject:  subject,
		HTMLBody: body}
	return sendMail(ctx, msg)
}

// simply send a mail
func sendMail(ctx context.Context, msg *mail.Message) (*mail.Message, error) {
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "couldnt_send_email %v", err)
		return msg, err
	}
	print(msg.Body)
	return msg, nil
}

/*
// create an alias from an http request and send validation mail
func putAlias(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	email := strings.TrimSpace(vars["email"])
	fullname := strings.TrimSpace(vars["fullname"])
	alias := new(Alias)
	lang := r.Header.Get("Accept-Language")
	alias, err := dsPutAliasSendValidationLink(ctx, lang, email, fullname)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, alias)
}
*/
