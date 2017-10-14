package main

import (
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
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	tr := getTranslater("en-EN")
	alias, error := dsPutAliasSendValidationLink(ctx, tr, testEmail, testFullname)
	if error != nil {
		t.Fatal(error)
	}
	assert.Equal(t, alias.Email, testEmail, "email should be stored")
	assert.Equal(t, alias.Fullname, testFullname, "fullname")
	assert.Equal(t, alias.Validated, false, "validated should be false")
	assert.Equal(t, len(alias.ValidationKey), 36, "validation key should exist")
	assert.Equal(t, len(alias.Alias), 36, "alias should exist")
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

func TestDatastoreGetAlias(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	var testEmail = "me@privacy.net"
	var testFullname = "John Doe"

	key := datastore.NewKey(ctx, "Alias", "", 1, nil)
	alias := &Alias{1,
		"me@privacy.net",
		"sfddsqfsdf@privacy.net",
		"John Doe",
		time.Now(),
		false,
		""}
	if _, err := datastore.Put(ctx, key, alias); err != nil {
		t.Fatal(err)
	}

	aliases, error := dsGetAlias(ctx, testEmail, testFullname)
	if error != nil {
		t.Fatal(error)
	}
	assert.Equal(t, len(aliases), 1, "there should not be aliases yet")
}

/**
func TestFilterByName(t *testing.T) {
	t.Fatal("TestFilterByName not implemented")
}

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
