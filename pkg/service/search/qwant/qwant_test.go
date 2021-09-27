package qwant

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := ioutil.ReadFile("qwant_test_resp.json")
		w.Write(resp)
	}))
	defer ts.Close()

	result, err := NewQwantSearchServiceRootURL(ts.URL).Search("fr_FR", "nantes")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal("https://fr.wikipedia.org/wiki/Le_Monde", result.Source)
	assert.Equal("https://fr.wikipedia.org/wiki/Le_Monde", result.URL)
	assert.NotEmpty(result.Text)
}
