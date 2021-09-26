package nominatim

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	"main/pkg/config"
)

func TestGeocode(t *testing.T) {
	config.InitWithRootPath("../../../config")

	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := ioutil.ReadFile("nominatim_test_resp.json")
		w.Write(resp)
	}))
	defer ts.Close()

	result, err := NewNominatimServiceRootURL(ts.URL).Geocode("FR", "nantes")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.NotNil(result)
	assert.NotNil(result.Lat)
	assert.NotNil(result.Lng)
}
