package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

// generic request build
func buildReqRes(t *testing.T, method string, url string) *httptest.ResponseRecorder {
	inst, _ := aetest.NewInstance(nil)
	defer inst.Close()

	req, err := inst.NewRequest(method, url, nil)
	if err != nil {
		assert.Fail(t, err.Error())
		return nil
	}
	ctx := appengine.NewContext(req)
	createTestApiKey(ctx)
	req.Header.Set("X-HMAC", "test@example.com,f13431fe08548c842dad4df339266cd80f2074c90147ddc829a6cb5588288c11")

	rec := httptest.NewRecorder()

	r := makeRouter()
	r.ServeHTTP(rec, req)
	return rec
}

// check a valid request / response for:
// "GET", "0

func createTestApiKey(ctx context.Context) (*ApiKey, error) {
	apiKey := new(ApiKey)
	apiKey.Email = "test@example.com"
	apiKey.ApiKey = "77706003-66f5-47e5-b52d-7d5d3f83a5c8"
	key := datastore.NewKey(ctx, "ApiKey", "test@example.com", 0, nil)
	_, err := datastore.Put(ctx, key, apiKey)
	if err != nil {
		return apiKey, err
	}
	return apiKey, nil
}

func TestPutAliasValid(t *testing.T) {

	rec := buildReqRes(t, "PUT", "/api/alias/test@example.com/john")
	// test correct response code
	assert.Equal(t, 200, rec.Code, "should respond 200")
	var alias Alias
	err := json.Unmarshal([]byte(rec.Body.String()), &alias)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test@example.com", alias.Email, "but "+rec.Body.String())
}

func TestPutAliasInvalid(t *testing.T) {
	// rec, _ := buildReqRes(t, "PUT", "/api/alias/meprivacy.net/jonh%20doe")
	// assert.Equal(t, 400, rec.Code, "bad email should respond 400")
}

func TestPutFullnameInvalid(t *testing.T) {
	// rec, _ := buildReqRes(t, "PUT", "/api/alias/meprivacy.net/%20a")
	// assert.Equal(t, rec.Code, 400, "bad fullname should respond 400")
}
