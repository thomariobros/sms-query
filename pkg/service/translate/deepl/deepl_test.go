package deepl

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslate(t *testing.T) {
	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := ioutil.ReadFile("deepl_test_resp.json")
		w.Write(resp)
	}))
	defer ts.Close()

	result, err := NewDeepLTranslateServiceRootURL(ts.URL).Translate("fr", "en", "maison")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal("house", result)
}
