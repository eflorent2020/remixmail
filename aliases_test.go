package main

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestValidateFullnameFormat(t *testing.T) {
	err := ValidateFullnameFormat("<script>")
	assert.NotNil(t, err, "tagged fullname raise error")
	err = ValidateFullnameFormat("")
	assert.NotNil(t, err, "empty fullname raise error")
	err = ValidateFullnameFormat("中英字典 YellowBridge")
	assert.Nil(t, err, "foreign chars should be accepted")
}

func TestValidateEmailFormat(t *testing.T) {
	err := ValidateEmailFormat("<script>alert</script>@test.net")
	assert.NotNil(t, err, "scripted email raise error")
	err = ValidateEmailFormat("")
	assert.NotNil(t, err, "empty fullname raise error")
	err = ValidateEmailFormat("meprivacy.net")
	assert.NotNil(t, err, "invalid mails should not be accepted")
	err = ValidateEmailFormat("me@pricacy.net")
	assert.Nil(t, err, "valid mail should be accepted")
}

func TestDatastorePutAlias(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	var testEmail = "me@privacy.net"
	var testFullname = "John Doe"
	assert.Nil(t, err, "test machinery should ge a context")
	defer done()
	alias, err := dsPutAliasSendValidationLink(ctx, "en-EN", testEmail, testFullname)
	assert.Nil(t, err, "dsPutAliasSendValidationLink should not return error")
	assert.Equal(t, alias.Email, testEmail, "email should be stored")
	assert.Equal(t, alias.Fullname, testFullname, "fullname")
	assert.Equal(t, alias.Validated, false, "validated should be false")
	assert.Equal(t, len(alias.ValidationKey), 36, "validation key should exist")
	assert.Contains(t, alias.Alias, "_", "alias should contain _")
	assert.Contains(t, alias.Alias, "@"+MAIL_DOMAIN, "alias should contain @MAIL_DOMAIN")

	year, month, day := alias.CreatedDate.Date()
	tyear, tmonth, tday := time.Now().Date()
	assert.Equal(t, year, tyear, "should store date - this year")
	assert.Equal(t, month, tmonth, "should store date - this month")
	assert.Equal(t, day, tday, "should store date - this day")
}

// get a context for test with datastore consistency activated
func getTestContext(t *testing.T) (context.Context, aetest.Instance) {
	inst, err := aetest.NewInstance(
		&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := appengine.NewContext(req)
	return ctx, inst
}

func makeTestAlias(ctx context.Context) (*Alias, error) {
	// key := datastore.NewKey(ctx, "Alias", "", 0, nil)
	key := aliasKey(ctx)
	alias := &Alias{1,
		"me@privacy.net",
		"test_test1@" + MAIL_DOMAIN,
		"John Doe",
		time.Now(),
		false,
		"aze123"}
	if _, err := datastore.Put(ctx, key, alias); err != nil {
		panic("cannot put alias")
		return alias, err
	}
	return alias, nil
}

func TestDatastoreGetAlias(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	var testEmail = "me@privacy.net"
	var testFullname = "John Doe"
	alias, err := makeTestAlias(ctx)
	assert.Nil(t, err, "cannot make test alias")
	aliases, err := dsGetAlias(ctx, testEmail, testFullname)
	assert.Nil(t, err, "cannot get alias")
	assert.Equal(t, alias.Email, aliases[0].Email, "should keep data after store")
	assert.Equal(t, len(aliases), 1, "there should not be aliases yet")
}

func TestSendValidationLink(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	alias, err := makeTestAlias(ctx)
	assert.Nil(t, err, "cannot make test alias")
	msg, err := sendValidationLink(ctx, "en-EN", alias)
	assert.Nil(t, err, "SendValidationLink should not return error")
	testLink := createConfirmationURL(alias)
	assert.Contains(t, msg.Body, alias.Fullname, "mail body should be personnalized")
	assert.Contains(t, msg.Body, testLink, "mail body should contains link")
	assert.Contains(t, msg.Sender, SENDER, "should send mail with configured sender")
	assert.Contains(t, msg.To[0], alias.Email, "should send mail to client address")
	// assert.Contains(t, msg.Subject, "hello", "subject well formed")
	assert.Contains(t, msg.Body, APP_NAME, "message sould be signed")
	assert.Contains(t, msg.Body, TAGLINE, "message sould be signed (2)")
	// println(msg.Body)
}

func TestCreateConfirmationURL(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	alias, err := makeTestAlias(ctx)
	assert.Nil(t, err, "cannot make test alias")
	testUrl := createConfirmationURL(alias)
	parsed, err := url.Parse(testUrl)
	assert.Nil(t, err, "cannot parse CreateConfirmationURL")
	assert.Equal(t, parsed.Scheme, "https")
	assert.Equal(t, parsed.Host, DOMAIN)
	assert.Equal(t, parsed.Path, "/")
	assert.Equal(t, parsed.Fragment, "/alias/validate/"+alias.ValidationKey)
}

/**
func TestDsValidateAlias(t *testing.T) {
	t.Fatal("TestDsValidateAlias not implemented")
}

func TestDsDeleteAliases(t *testing.T) {
	t.Fatal("TestDsDeleteAliases not implemented")
}

func TestDsDeleteAlias(t *testing.T) {
	t.Fatal("TestDsDeleteAlias not implemented")
}

func TestSendValidationLink(t *testing.T) {
	t.Fatal("TestSendValidationLink not implemented")
}

func TestCreateConfirmationURL(t *testing.T) {
	t.Fatal("TestCreateConfirmationURL not implemented")
}

func TestDsFindAliased(t *testing.T) {
	t.Fatal("TestCreateConfirmationURL not implemented")
}

*/
