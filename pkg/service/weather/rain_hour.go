package weather

import (
	"encoding/json"

	"net/url"

	"main/pkg/http"
)

const (
	rainHourRootURL = "https://www.meteofrance.com/mf3-rpc-portlet/rest"
)

type RainHourService struct {
	rootURL string
}

func NewRainHourService() *RainHourService {
	return &RainHourService{
		rootURL: rainHourRootURL,
	}
}

func NewRainHourServiceRootURL(rootURL string) *RainHourService {
	return &RainHourService{
		rootURL: rootURL,
	}
}

// GetRainHour get rain information from meteofrance int the next hour
func (s *RainHourService) GetRainHour(city string) (*RainHour, error) {
	if city == "" {
		return nil, nil
	}
	// search city
	resp, err := http.Get(s.rootURL + "/lieu/facet/previsions/search/" + url.QueryEscape(city))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var searchCityResult []searchCityResult
	err = json.NewDecoder(resp.Body).Decode(&searchCityResult)
	if err != nil {
		return nil, err
	}
	if len(searchCityResult) > 0 {
		// load city data
		resp, err := http.Get(s.rootURL + "/pluie/" + searchCityResult[0].ID)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var result RainHour
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		if len(result.Data) > 0 {
			return &result, nil
		}
	}

	return nil, nil
}

type searchCityResult struct {
	ID string `json:"id"`
}

type RainHour struct {
	IsAvailable bool           `json:"isAvailable"`
	HasData     bool           `json:"hasData"`
	Deadline    string         `json:"echeance"`
	Data        []RainHourData `json:"dataCadran"`
}

type RainHourData struct {
	Level int `json:"niveauPluie"`
}
