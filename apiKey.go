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
	ApiKey string `datastore:"api_key"`
	Email  string `datastore:"email"`
}

// take an email,create an api key and put in datastore
// the api key will be sent by email
func dsPutAPiKey(ctx context.Context, T i18n.TranslateFunc, email string) (*ApiKey, error) {
	err := ValidateEmailFormat(email)
	var apiKey *ApiKey
	if err != nil {
		return apiKey, err
	}
	apiKey = new(ApiKey)
	apiKey.Email = email
	apiKey.ApiKey = uuid.NewV4().String()
	key := datastore.NewKey(ctx, "ApiKey", email, 0, nil)
	_, err = datastore.Put(ctx, key, apiKey)
	if err != nil {
		return apiKey, err
	}
	sendApiKey(ctx, T, apiKey)
	return apiKey, nil
}

// send a friendly message containing the api key
func sendApiKey(ctx context.Context, T i18n.TranslateFunc, apiKey *ApiKey) error {
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

// get the api key for a customer email address
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



// dsGetAliases take an email as argument and return an array of
// all Alias struc
func listAPiKey(ctx context.Context) ([]ApiKey, error) {
	q := datastore.NewQuery("ApiKey")
	var apiKeys []ApiKey
	if _, err := q.GetAll(ctx, &apiKeys); err != nil {
		return apiKeys, err
	}
	return apiKeys, nil
}
