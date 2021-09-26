package news

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestSearch(t *testing.T) {
	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := ioutil.ReadFile("news_test_resp.xml")
		w.Write(resp)
	}))
	defer ts.Close()

	result, err := Rss(ts.URL, 3)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(3, len(result))
	for i := 0; i < len(result); i++ {
		assert.NotEmpty(result[0].Title)
		assert.NotEmpty(result[0].Description)
	}
}
