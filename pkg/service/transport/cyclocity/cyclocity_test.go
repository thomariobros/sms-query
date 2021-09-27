package cyclocity

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"testing"

	"github.com/stretchr/testify/assert"

	"main/pkg/config"
)

func TestCyclocity(t *testing.T) {
	config.InitWithRootPath("../../../config")

	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := ioutil.ReadFile("cyclocity_test_resp.json")
		w.Write(resp)
	}))
	defer ts.Close()

	result, err := NewCyclocityServiceRootURL(ts.URL).GetStations("nantes", "gare maritime", 3)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(1, len(result))
	assert.Equal("00042-gare maritime", result[0].Name)
	assert.NotEmpty(result[0].Address)
	assert.True(result[0].AvailableBikes >= 0)
	assert.True(result[0].AvailableBikeStands >= 0)
}
