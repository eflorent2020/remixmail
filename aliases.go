package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/i18n"
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
	ID            int64
	Email         string    `datastore:"email"`
	Alias         string    `datastore:"alias"`
	Fullname      string    `datastore:"fullname"`
	CreatedDate   time.Time `datastore:"created_at"`
	Validated     bool      `datastore:"validated"`
	ValidationKey string    `datastore:"validation_key"`
	Domain        string
}

func aliasKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Alias", "default_alias", 0, nil)
}

// take an email, a fullname, check email format and put in datastore
func dsPutAliasSendValidationLink(ctx context.Context, T i18n.TranslateFunc, email string, fullname string) (*Alias, error) {
	err := ValidateEmailFormat(email)
	alias := new(Alias)
	if err != nil {
		return nil, err
	}
	err = ValidateFullnameFormat(fullname)
	if err != nil {
		return nil, err
	}
	checkExist, err := dsGetAlias(ctx, email, fullname)
	if checkExist != nil {
		return nil, errors.New("alias already exists")
	}
	alias.CreatedDate = time.Now()
	alias.Email = email
	alias.Fullname = fullname
	alias.Validated = false
	alias.Alias = uuid.NewV4().String()
	alias.ValidationKey = uuid.NewV4().String()
	key := datastore.NewIncompleteKey(ctx, "Alias", aliasKey(ctx))
	_, err = datastore.Put(ctx, key, alias)
	if err != nil {
		return nil, err
	}
	sendValidationLink(ctx, T, alias)
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
	aliasKey := strings.Split(email, "@")[0]
	q := datastore.NewQuery("Alias").Filter("alias = ", aliasKey)
	var aliases []Alias
	var alias Alias
	if _, err := q.GetAll(ctx, &aliases); err != nil {
		return alias, err
	}
	if len(aliases) != 1 {
		return alias, errors.New("inconsistent aliased " + aliasKey)
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
	aliases2 := filterByName(aliases, fullname)
	if fullname != "" {
		return filterByName(aliases2, fullname), nil
	} else {
		return aliases2, nil
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
	for _, a := range aliases {
		//func NewKey(c context.Context, kind, stringID string, intID int64, parent *Key) *Key
		key := datastore.NewKey(ctx, "Alias", "", a.ID, nil)
		err = dsDeleteAlias(r, key)
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

func createConfirmationURL(alias *Alias) string {
	return APP_ROOT_URL + "/#/alias/validate/" + alias.ValidationKey
}

func sendValidationLink(ctx context.Context, T i18n.TranslateFunc, alias *Alias) {
	addr := alias.Email
	url := createConfirmationURL(alias)
	msg := &mail.Message{
		Sender:  APP_NAME + " <" + SENDER + ">",
		To:      []string{addr},
		Subject: T("confirm_mail_alias"),
		Body: fmt.Sprintf(T("confirm_message",
			map[string]interface{}{"link": url})),
	}
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, T("couldnt_send_email"), err)
	}
	log.Infof(ctx, "validation email sent to "+addr+" "+url)
}
