package main

import (
	"net/mail"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// set some data when test need it
func setSampleAlias(t *testing.T, ctx context.Context) {
	key := datastore.NewKey(ctx, "Alias", "", 1, nil)
	aliasTest := &Alias{1,
		"me@privacy.net",
		"5a700b3b-11d8-4874-bea6-8b653d3a0592",
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
To: Recipient Name <5a700b3b-11d8-4874-bea6-8b653d3a0592@privacy.net>
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
	assert.Equal(t, 36, len(alias.Alias), "alias should exist for incoing mail")
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
		"bob00b3b-11d8-4874-bea6-8b653d3a0592",
		"Bob Doe",
		time.Now(),
		true,
		""}
	aliasTo := &Alias{1,
		"alice@privacy.net",
		"aliceb3b-11d8-4874-bea6-8b653d3a0592",
		"Alice Doe",
		time.Now(),
		true,
		""}

	aeMsg := buildForward(ctx, aliasFrom, aliasTo, msg)
	assert.Equal(t, "Bob Doe <bob00b3b-11d8-4874-bea6-8b653d3a0592@"+DOMAIN+">", aeMsg.Sender, "from header should be translated")
	assert.Equal(t, "Alice Doe <alice@privacy.net>", aeMsg.To[0], "to header should be translated")
}
