package main

import (
	"encoding/json"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// set some data when test need it
func setSampleAlias(t *testing.T, ctx context.Context) {
	key := datastore.NewKey(ctx, "Alias", "", 1, nil)
	aliasTest := &Alias{1,
		"me@privacy.net",
		"test_test1@" + MAIL_DOMAIN,
		"John Doe",
		time.Now(),
		true,
		""}
	if _, err := datastore.Put(ctx, key, aliasTest); err != nil {
		t.Fatal(err)
	}
}

// build a sample mail for test purpose
func getTestMail(t *testing.T) *mail.Message {
	sample := `Return-path: <sender@senderdomain.tld>
Delivery-date: Wed, 13 Apr 2011 00:31:13 +0200
Message-ID: <xxxxxxxx.xxxxxxxx@senderdomain.tld>
Date: Tue, 12 Apr 2011 20:36:01 -0100
X-Mailer: Mail Client
From: Sender Name <sender@senderdomain.tld>
To: Recipient Name <test_test1@` + MAIL_DOMAIN + `>
Subject: Message Subject

This is the body...
`
	r := strings.NewReader(sample)
	m, err := mail.ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestGetAliasFrom(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	m := getTestMail(t)
	alias, err := getAliasFrom(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "sender@senderdomain.tld", alias.Email, "should find alias for any incoing mail")
	assert.Contains(t, alias.Alias, "_", "alias should exist for incoing mail")
}

func TestGetAliasTo(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	// check without alias
	_, err := getAliasTo(ctx, getTestMail(t))
	assert.NotNil(t, err, "alias should not exist")
	setSampleAlias(t, ctx)
	alias, err := getAliasTo(ctx, getTestMail(t))
	assert.Nil(t, err, "valid recipent does not return errors")
	assert.Equal(t, "me@privacy.net", alias.Email, "alias to should be decoded")
}

func TestBuildForward(t *testing.T) {

	ctx, inst := getTestContext(t)
	defer inst.Close()

	msg := getTestMail(t)

	aliasFrom := &Alias{1,
		"bob@privacy.net",
		"bob_test@" + MAIL_DOMAIN,
		"Bob Doe",
		time.Now(),
		true,
		""}
	aliasTo := &Alias{1,
		"alice@privacy.net",
		"alice_test@" + MAIL_DOMAIN,
		"Alice Doe",
		time.Now(),
		true,
		""}

	aeMsg := buildForward(ctx, aliasFrom, aliasTo, msg)
	assert.Equal(t, "Bob Doe <bob_test@"+MAIL_DOMAIN+">", aeMsg.Sender, "should send mail from aliased")
	assert.Equal(t, "Alice Doe <alice@privacy.net>", aeMsg.To[0], "should send mail to real")
}

func TestIncomingMail4ServiceRegister(t *testing.T) {
	_, inst := getTestContext(t)
	defer inst.Close()

	sample := `Return-path: <sender@senderdomain.tld>
Delivery-date: Wed, 13 Apr 2011 00:31:13 +0200
Message-ID: <xxxxxxxx.xxxxxxxx@senderdomain.tld>
Date: Tue, 12 Apr 2011 20:36:01 -0100
X-Mailer: Mail Client
From: Sender Name <sender@senderdomain.tld>
To: Recipient Name <system@snapmail-182207.appspotmail.com>
Subject: ReGister

This is the body...
`
	r := strings.NewReader(sample)

	req1, err := inst.NewRequest("GET", "/_ah/mail/", r)
	if err != nil {
		t.Fatalf("Failed to create req1: %v", err)
	}
	_ = appengine.NewContext(req1)
	w := httptest.NewRecorder()

	incomingMail(w, req1)
	str := w.Body.String()
	res := Alias{}
	json.Unmarshal([]byte(str), &res)
	assert.Equal(t, "sender@senderdomain.tld", res.Email, "email should be acquired acquired")
}
