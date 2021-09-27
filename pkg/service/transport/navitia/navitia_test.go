package navitia

import (
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"main/pkg/config"
	"main/pkg/service/transport"
)

func TestStifServiceGetNextDepartureTime(t *testing.T) {
	config.InitWithRootPath("../../../config")

	// mock server response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp []byte
		if strings.Contains(r.URL.String(), "type[]=stop_area") {
			resp, _ = ioutil.ReadFile("navitia_test_resp_stops.json")
		} else if strings.Contains(r.URL.String(), "type[]=line") {
			resp, _ = ioutil.ReadFile("navitia_test_resp_lines.json")
		} else {
			resp, _ = ioutil.ReadFile("navitia_test_resp_departures.json")
		}
		w.Write(resp)
	}))
	defer ts.Close()

	loc, _ := time.LoadLocation("Europe/Paris")
	nowBefore := time.Now().In(loc)

	result, err := NewNavitiaStifServiceRootURL(ts.URL).GetNextDepartureTime("m", "1", "bastille", "la defense", 3)
	if err != nil {
		t.Fatal(err)
	}

	nowAfter := time.Now().In(loc)

	assert := assert.New(t)
	assert.NotNil(result)
	assert.Equal(transport.StatusOk, result.Status)
	assert.NotNil(result.Now)
	assert.NotNil(result.Now)
	assert.True(result.Now.After(nowBefore))
	assert.True(result.Now.Before(nowAfter))
}
