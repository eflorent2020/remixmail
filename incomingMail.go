package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/mail"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	AppEngineMail "google.golang.org/appengine/mail"
)

// Handle all incoming mails
func incomingMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer r.Body.Close()
	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Errorf(ctx, "Error reading mail body: %v", err)
		return
	}
	// Todo parse address
	rcptTo := msg.Header.Get("To")
	if rcptTo == SERVICEMAIL {
		serviceMail(ctx, msg)
		return
	}

	aliasFrom, err := getAliasFrom(ctx, *msg)
	if err != nil {
		log.Errorf(ctx, "Error getting aliasFrom : %v", err)
		return
	}
	aliasTo, err := getAliasFrom(ctx, *msg)
	if err != nil {
		log.Errorf(ctx, "Error getting aliasTo: %v", err)
		return
	}
	forward(ctx, *aliasFrom, *aliasTo, *msg)
	log.Infof(ctx, "Done with : %v  -> %v", aliasFrom.ID, aliasTo.ID)
}

// taking an incoming message get and check validation of an alias for the recipent
func getAliasTo(ctx context.Context, msgReceived mail.Message) (Alias, error) {
	parsedTo, err := mail.ParseAddress(msgReceived.Header.Get("From"))
	if err != nil {
		log.Errorf(ctx, "unable to parse from adress", err)
		var alias Alias
		return alias, err
	}
	alias, err := dsFindAliased(ctx, parsedTo.Address)
	if err != nil {
		log.Errorf(ctx, "Error finding aliased", err)
		return alias, err
	}
	if alias.Validated == false {
		log.Errorf(ctx, "Error received not validated recipient", err)
		// TODO inform or blacklist sender
		return alias, errors.New("the recipient did not validate")
	}
	return alias, nil
}

// taking an incoming message get or create an alias for the sender
func getAliasFrom(ctx context.Context, msgReceived mail.Message) (*Alias, error) {
	parsedFrom, err := mail.ParseAddress(msgReceived.Header.Get("From"))
	if err != nil {
		log.Errorf(ctx, "unable to parse from adress", err)
		var alias Alias
		return &alias, err
	}
	// try find and check if active
	// TODO dsFindOrCreate
	// if not found
	translater := getTranslater(msgReceived.Header.Get("Accept-Language"))
	aliasFrom, err := dsPutAlias(ctx, translater, parsedFrom.Address, parsedFrom.Name)
	if err != nil {
		log.Errorf(ctx, "unable to put alias", err)
		var alias Alias
		return &alias, err
	}
	return aliasFrom, nil
}

// TODO should handle mail sent to a service bot
func serviceMail(ctx context.Context, msg *mail.Message) {
	log.Infof(ctx, "Received service mail : %v", msg.Body)
}

// Taking a mail.Message, from alias, to alias, try do a forward
// TODO more on spam protection & reply on error
func forward(ctx context.Context, aliasFrom Alias, aliasTo Alias, msgReceived mail.Message) {
	var to []string
	to[0] = aliasTo.Fullname + "<" + aliasTo.Email + ">"
	buf := new(bytes.Buffer)
	buf.ReadFrom(msgReceived.Body)
	body := buf.String()
	msg := &AppEngineMail.Message{
		Sender:  aliasFrom.Fullname + " <" + aliasFrom.Alias + ">",
		To:      to,
		Subject: msgReceived.Header.Get("subject"),
		Body:    body,
	}
	if err := AppEngineMail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "couldnt_send_email", err)
	}
}
