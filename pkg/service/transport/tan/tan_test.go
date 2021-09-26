package tan

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"strings"
	"testing"
	"time"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestTanServiceGetNextDepartureTime(t *testing.T) {
	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp []byte
		if strings.Contains(r.URL.String(), "arrets") {
			resp, _ = ioutil.ReadFile("tan_test_resp_arrets.json")
		} else {
			resp, _ = ioutil.ReadFile("tan_test_resp_tempsattente.json")
		}
		w.Write(resp)
	}))
	defer ts.Close()

	loc, _ := time.LoadLocation("Europe/Paris")
	nowBefore := time.Now().In(loc)

	result, err := NewTanServiceRootURL(ts.URL).GetNextDepartureTime("t", "1", "lauriers", "beaujoire", 3)
	if err != nil {
		t.Fatal(err)
	}

	nowAfter := time.Now().In(loc)

	assert := assert.New(t)
	assert.NotNil(result)
	assert.NotNil(result.Now)
	assert.True(result.Now.After(nowBefore))
	assert.True(result.Now.Before(nowAfter))
}
