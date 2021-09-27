package weather

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRainHour(t *testing.T) {
	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp []byte
		if strings.Contains(r.URL.String(), "lieu") {
			resp, _ = ioutil.ReadFile("rain_hour_resp_lieu.json")
		} else {
			resp, _ = ioutil.ReadFile("rain_hour_resp_pluie.json")
		}
		w.Write(resp)
	}))
	defer ts.Close()

	result, err := NewRainHourServiceRootURL(ts.URL).GetRainHour("nantes")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.NotNil(result)
	assert.True(result.IsAvailable)
	assert.True(result.HasData)
	assert.NotNil(result.Deadline)
	assert.NotEmpty(result.Data)
	for _, data := range result.Data {
		assert.NotNil(data.Level)
	}
}
