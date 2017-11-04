package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/jhillyerd/go.enmime"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	aeMail "google.golang.org/appengine/mail"
)

// take a mail.Message, decode it in Appengine Mail Message format
// Appengine mail format distingish text, html and attachments
func decodeMail(ctx context.Context, msg *mail.Message) (string, string, []aeMail.Attachment) {
	// Parse message body with enmime
	mime, err := enmime.ParseMIMEBody(msg)
	if err != nil {
		log.Errorf(ctx, "During enmime.ParseMIMEBody: %v", err)
	}
	var atchs []aeMail.Attachment

	for i, a := range mime.Attachments {
		ath := aeMail.Attachment{
			a.FileName(),
			a.Content(),
			"<content" + strconv.Itoa(i) + ">"}
		atchs = append(atchs, ath)
	}
	return mime.Text, mime.HTML, atchs
}

// take a PGP public key and a mail.Message, decode it in Appengine Mail Message format
// then encrypt all parts it using PGP
func encryptAndDecodeMail(ctx context.Context, key string, msg *mail.Message) (string, string, []aeMail.Attachment) {
	pubkey, _ := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(key))
	// Parse message body with enmime
	mime, err := enmime.ParseMIMEBody(msg)
	if err != nil {
		log.Errorf(ctx, "During enmime.ParseMIMEBody: %v", err)
	}
	var atchs []aeMail.Attachment

	for i, a := range mime.Attachments {
		//content, _ := gpg_pubkey_encrypt(ctx, nil, string(a.Content()), key)
		buf := new(bytes.Buffer)
		w, _ := openpgp.Encrypt(buf, pubkey, nil, nil, nil)
		_, err = w.Write(a.Content())
		if err != nil {
		}
		err = w.Close()
		bytes, _ := ioutil.ReadAll(buf)
		ath := aeMail.Attachment{
			a.FileName() + ".pgp",
			bytes,
			"<content" + strconv.Itoa(i) + ">"}
		atchs = append(atchs, ath)
	}
	text, _ := gpg_pubkey_encrypt(ctx, "PGP MESSAGE", mime.Text, key)
	html, _ := gpg_pubkey_encrypt(ctx, "PGP MESSAGE", mime.HTML, key)
	return text, html, atchs
}

// build a part Mime-Part encrypted message, that could be text or html
func gpg_pubkey_encrypt(ctx context.Context, header, s, key string) (string, error) {
	buf := new(bytes.Buffer)
	msg, err := armor.Encode(buf, header, nil)
	if err != nil {
		log.Errorf(ctx, "error armoring %+v", err)
		return s, err
	}
	pubkey, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(key))
	if err != nil {
		log.Errorf(ctx, "error reading pub key %+v", err)
		return s, err
	}
	gpg, err := openpgp.Encrypt(msg, pubkey, nil, nil, nil)
	if err != nil {
		log.Errorf(ctx, "error encrypting %+v", err)
		return s, err
	}
	fmt.Fprintf(gpg, s)
	gpg.Close()
	msg.Close()
	return buf.String(), nil
}

// Taking a mail.Message, from alias, to alias, try do a forward
func buildForward(ctx context.Context, aliasFrom *Alias, aliasTo *Alias, msgReceived *mail.Message) (msg aeMail.Message) {
	// recover/fallback using straight body mail if parse failed ...
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(ctx, "was panic, recovered value: %v", r)
			buf := new(bytes.Buffer)
			buf.ReadFrom(msgReceived.Body)
			// (re) build a mail
			msg = aeMail.Message{
				Sender:  aliasFrom.Fullname + " <" + aliasFrom.Alias + ">",
				To:      []string{aliasTo.Fullname + " <" + aliasTo.Email + ">"},
				Subject: msgReceived.Header.Get("subject"),
				Body:    buf.String()}
		}
	}()

	// build a mail
	if len(aliasTo.PGPPubKey) > 0 {
		log.Infof(ctx, "enforcing PGP mail")
		body, html, atchs := encryptAndDecodeMail(ctx, string(aliasTo.PGPPubKey), msgReceived)
		msg = aeMail.Message{
			Sender:      aliasFrom.Fullname + " <" + aliasFrom.Alias + ">",
			To:          []string{aliasTo.Fullname + " <" + aliasTo.Email + ">"},
			Subject:     msgReceived.Header.Get("subject"),
			Body:        body,
			HTMLBody:    html,
			Attachments: atchs}
		return msg
	} else {
		log.Infof(ctx, "sending unencrypted mail")
		body, html, atchs := decodeMail(ctx, msgReceived)
		msg = aeMail.Message{
			Sender:      aliasFrom.Fullname + " <" + aliasFrom.Alias + ">",
			To:          []string{aliasTo.Fullname + " <" + aliasTo.Email + ">"},
			Subject:     msgReceived.Header.Get("subject"),
			Body:        body,
			HTMLBody:    html,
			Attachments: atchs}
		return msg
	}
}

// Handle all incoming mails
func incomingMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer r.Body.Close()
	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Errorf(ctx, "Error reading mail body: %v", err)
		return
	}
	parsedTo, err := mail.ParseAddress(msg.Header.Get("To"))
	if err != nil {
		log.Errorf(ctx, "Error parsing to address : %v", err)
		return
	}
	rcptTo := parsedTo.Address
	if rcptTo == SERVICE_MAIL {
		alias, err := serviceMail(ctx, msg)
		if err == nil {
			respondWithJSON(w, http.StatusOK, alias)
		} else {
			log.Errorf(ctx, "service received junk mail : %v", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	// must check to exist before from because from will create alias
	aliasTo, err := getAliasTo(ctx, msg)
	if err != nil {
		log.Errorf(ctx, "Error getting aliasTo: %v", msg.Header.Get("To"))
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	aliasFrom, err := getAliasFrom(ctx, msg)
	if err != nil {
		log.Errorf(ctx, "Error getting aliasFrom : %v", err)
		return
	}

	// TODO more on error handling
	msgToForward := buildForward(ctx, aliasFrom, aliasTo, msg)
	forwardSend(ctx, &msgToForward)
	return
}

// taking an incoming message get and check validation of an alias for the recipent
func getAliasTo(ctx context.Context, msgReceived *mail.Message) (*Alias, error) {
	parsedTo, err := mail.ParseAddress(msgReceived.Header.Get("To"))
	if err != nil {
		log.Errorf(ctx, "unable to parse from adress", err)
		var alias Alias
		return &alias, err
	}
	alias, err := dsFindAliased(ctx, parsedTo.Address)
	if err != nil {
		log.Errorf(ctx, "Error finding aliased", err)
		return &alias, err
	}
	if alias.Validated == false {
		log.Errorf(ctx, "The recipient did not yet validate it's address", err)
		// TODO inform or blacklist sender
		return &alias, errors.New("the recipient did not validate")
	}
	return &alias, nil
}

// taking an incoming message get or create an alias for the sender
func getAliasFrom(ctx context.Context, msgReceived *mail.Message) (*Alias, error) {
	parsedFrom, err := mail.ParseAddress(msgReceived.Header.Get("From"))
	if err != nil {
		log.Errorf(ctx, "unable to parse from adress", err)
		var alias Alias
		return &alias, err
	}

	// try find and check if active can ignore error not suposed to get the line
	aliasesFrom, _ := dsGetAlias(ctx, parsedFrom.Address, "")

	//sender does to exist in the database
	if len(aliasesFrom) < 1 {
		log.Infof(ctx, "creating user from getAliasFrom")
		aliasFrom, err := dsPutAliasSendValidationLink(ctx, msgReceived.Header.Get("Accept-Language"), parsedFrom.Address, parsedFrom.Name)
		if err != nil {
			log.Errorf(ctx, "unable to put alias", err)
			return aliasFrom, err
		}
		return aliasFrom, nil
	}

	// sender exist in the database
	aliasFrom := aliasesFrom[0]
	return &aliasFrom, nil
}

// handle mail sent to a service bot, if a message is received,
// with the subject register, create an account
func serviceMail(ctx context.Context, msg *mail.Message) (*Alias, error) {
	var alias *Alias
	if !strings.Contains("register", strings.ToLower(msg.Header.Get("Subject"))) {
		return alias, errors.New("KO: subject must contain register was " + strings.ToLower(msg.Header.Get("Subject")))
	}
	address, err := mail.ParseAddress(msg.Header.Get("From"))
	if err != nil {
		return alias, err
	}
	log.Infof(ctx, "creating user from serviceMail")
	alias, err = dsPutAliasSendValidationLink(ctx, msg.Header.Get("Accept-Language"), address.Address, address.Name)
	if err != nil {
		return alias, err
	}
	log.Infof(ctx, "Received service mail")
	return alias, nil
}

func forwardSend(ctx context.Context, msg *aeMail.Message) error {
	if err := aeMail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "couldnt_send_email", err)
		return err
	}
	return nil
}
