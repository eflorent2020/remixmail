package main

import (
	"errors"
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
)

// ApiKey is requested to create new Alias or to make operations
// on aliases without a validation key
// in request header api_key: your-api-key
type ApiKey struct {
	ID     int64
	ApiKey string `datastore:"email"`
	Email  string `datastore:"created_at"`
}

func apiKeyKey(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "ApiKey", "default_api_key", 0, nil)
}

// take an email,create an api key and put in datastore
// the api key will be sent by email
func dsPutAPiKey(ctx context.Context, T i18n.TranslateFunc, email string) (ApiKey, error) {
	err := ValidateEmailFormat(email)
	var apiKey ApiKey
	if err != nil {
		return apiKey, err
	}
	apiKey, err = dsGetApiKeyFor(ctx, email)
	if err == nil {
		sendApiKey(ctx, T, apiKey)
		return apiKey, nil
	}
	apiKey.Email = email
	apiKey.ApiKey = uuid.NewV4().String()
	key := datastore.NewIncompleteKey(ctx, "ApiKey", apiKeyKey(ctx))
	_, err = datastore.Put(ctx, key, apiKey)
	if err != nil {
		return apiKey, err
	}
	sendApiKey(ctx, T, apiKey)
	return apiKey, nil
}

// send a friendly message containing the api key
func sendApiKey(ctx context.Context, T i18n.TranslateFunc, apiKey ApiKey) error {
	addr := apiKey.Email
	key := apiKey.ApiKey
	msg := &mail.Message{
		Sender:  APP_NAME + " <" + SENDER + ">",
		To:      []string{addr},
		Subject: T("your_api_key"),
		Body: fmt.Sprintf(T("api_key_message",
			map[string]interface{}{"Key": key})),
	}
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "couldn't send email", err.Error())
		return err
	}
	return nil
}

func dsGetApiKeyFor(ctx context.Context, email string) (ApiKey, error) {
	q := datastore.NewQuery("ApiKey").Filter("Email = ", email)
	var apiKeys []ApiKey
	var apiKey ApiKey
	if _, err := q.GetAll(ctx, &apiKeys); err != nil {
		return apiKey, err
	}
	if len(apiKeys) == 0 {
		return apiKey, errors.New("not api key found")
	}
	return apiKeys[0], nil
}

func checkAPiKey(ctx context.Context, key string) (ApiKey, error) {
	q := datastore.NewQuery("ApiKey").Filter("ApiKey = ", key)
	var apiKeys []ApiKey
	if _, err := q.GetAll(ctx, &apiKeys); err != nil {
		var apiKey ApiKey
		return apiKey, err
	}
	return apiKeys[0], nil
}
