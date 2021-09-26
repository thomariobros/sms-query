package cyclocity

import (
	"encoding/json"
	"strings"

	"main/pkg/config"
	"main/pkg/http"
	"main/pkg/util"
)

const (
	rootURL = "https://api.jcdecaux.com/vls/v1/stations"
)

type CyclocityService struct {
	rootURL string
}

func NewCyclocityService() *CyclocityService {
	return &CyclocityService{
		rootURL: rootURL,
	}
}

func NewCyclocityServiceRootURL(rootURL string) *CyclocityService {
	return &CyclocityService{
		rootURL: rootURL,
	}
}

// GetStations get cyclocity stations through jcdecaux api
func (s *CyclocityService) GetStations(contract string, name string, max int) ([]CyclocityStation, error) {
	params := "contract=" + contract +
		"&apiKey=" + config.GetInstance().Cyclocity.APIKey
	resp, err := http.Get(s.rootURL + "?" + params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var stations []CyclocityStation
	err = json.NewDecoder(resp.Body).Decode(&stations)
	if err != nil {
		return nil, err
	}
	// search 'OPEN' stations that match
	var result []CyclocityStation
	for _, station := range stations {
		station.standardizeStrings()
		if station.Status == "open" && station.match(name) {
			result = append(result, station)
			if len(result) == max {
				break
			}
		}
	}

	return result, nil
}

type CyclocityStation struct {
	Name                string `json:"name"`
	Address             string `json:"address"`
	Status              string `json:"status"`
	AvailableBikeStands int    `json:"available_bike_stands"`
	AvailableBikes      int    `json:"available_bikes"`
}

func (station *CyclocityStation) standardizeStrings() {
	station.Status = util.StandardizeStringLower(station.Status)
	station.Name = strings.Replace(util.StandardizeStringLower(station.Name), "- ", "", -1)
	station.Address = strings.Replace(util.StandardizeStringLower(station.Address), "- ", "", -1)
}

func (station *CyclocityStation) match(name string) bool {
	return strings.Contains(station.Name, name) ||
		strings.Contains(station.Address, name)
}
