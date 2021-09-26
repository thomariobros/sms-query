package nominatim

import (
	"encoding/json"
	"net/url"
	"strconv"

	"main/pkg/http"
	"main/pkg/service/maps"
)

const (
	nominatimRootURL = "https://nominatim.openstreetmap.org"
)

type NominatimService struct {
	rootURL string
}

func NewNominatimService() *NominatimService {
	return &NominatimService{
		rootURL: nominatimRootURL,
	}
}

func NewNominatimServiceRootURL(rootURL string) *NominatimService {
	return &NominatimService{
		rootURL: rootURL,
	}
}

// Geocode geocode city using openstreetmap nominatim (https://wiki.openstreetmap.org/wiki/Nominatim)
func (s *NominatimService) Geocode(country, city string) (*maps.Location, error) {
	if city == "" {
		return nil, nil
	}
	params := "city=" + url.QueryEscape(city) +
		"&countrycodes=" + country +
		"&format=json"
	resp, err := http.Get(s.rootURL + "/search" + "?" + params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response []location
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	if len(response) > 0 {
		return &maps.Location{Lat: parseFloat(response[0].Lat), Lng: parseFloat(response[0].Lon)}, nil
	}
	return nil, nil
}

func parseFloat(input string) float64 {
	v, _ := strconv.ParseFloat(input, 64)
	return v
}

type location struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}
