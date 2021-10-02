package nominatim

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"testing"

	"github.com/stretchr/testify/assert"

	"main/pkg/config"
)

func TestGeocode(t *testing.T) {
	err := config.InitWithPath("../../../../config/config.yml")
	if err != nil {
		t.Fatal(err)
	}

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
