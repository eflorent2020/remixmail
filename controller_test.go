package main

import (
	"net/http/httptest"
	"testing"

	"google.golang.org/appengine/aetest"
)

// generic request build
func buildReqRes(t *testing.T, method string, url string) *httptest.ResponseRecorder {
	inst, _ := aetest.NewInstance(nil)
	defer inst.Close()

	req, _ := inst.NewRequest(method, url, nil)
	rec := httptest.NewRecorder()

	r := makeRouter()
	r.ServeHTTP(rec, req)
	return rec
}

// check a valid request / response for:
// "GET", "/api/alias/me@privacy.net/jonh%20doe
/*
func TestPutAliasValid(t *testing.T) {
	rec := buildReqRes(t, "PUT", "/api/alias/me@privacy.net/jonh%20doe")
	// test correct response code
	assert.Equal(t, 200, rec.Code, "should respond 200")

	var alias Alias
	err := json.Unmarshal([]byte(rec.Body.String()), &alias)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "me@privacy.net", alias.Email, "should store email field")
}

func TestPutAliasInvalid(t *testing.T) {
	rec := buildReqRes(t, "PUT", "/api/alias/meprivacy.net/jonh%20doe")
	assert.Equal(t, 400, rec.Code, "bad email should respond 400")
}

func TestPutFullnameInvalid(t *testing.T) {
	rec := buildReqRes(t, "PUT", "/api/alias/meprivacy.net/%20a")
	assert.Equal(t, rec.Code, 400, "bad fullname should respond 400")
}
*/
