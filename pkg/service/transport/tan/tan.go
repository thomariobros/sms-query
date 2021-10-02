package tan

import (
	"encoding/json"
	"strings"
	"time"

	"sms-query/pkg/http"
	"sms-query/pkg/service/transport"
	"sms-query/pkg/util"
)

const (
	rootURL = "http://open.tan.fr/ewp"
)

type TanService struct {
	rootURL string
}

func NewTanService() *TanService {
	return &TanService{
		rootURL: rootURL,
	}
}

func NewTanServiceRootURL(rootURL string) *TanService {
	return &TanService{
		rootURL: rootURL,
	}
}

// GetNextDepartureTime get next departure time for Nantes city through tan API
func (s *TanService) GetNextDepartureTime(lineType, lineID, stopName, direction string, max int) (*transport.NextDeparture, error) {
	// check line
	lineTypeCode := getInternalLineType(lineType)
	if lineTypeCode == 0 {
		return &transport.NextDeparture{Status: transport.StatusUnknownLine}, nil
	}
	// find stop
	stopAndStatus, err := s.findStop(lineID, stopName)
	if err != nil {
		return nil, err
	}
	if stopAndStatus.status != transport.StatusOk {
		return &transport.NextDeparture{Status: stopAndStatus.status}, nil
	}
	// retrieve stop times
	stopTimesAndStatus, err := s.retrieveStopTimes(lineTypeCode, lineID, stopAndStatus.stop, direction)
	if err != nil {
		return nil, err
	}
	if stopTimesAndStatus.status == transport.StatusOk {
		return stopTimesToNextDeparture(stopTimesAndStatus.stopTimes, max), nil
	}
	return &transport.NextDeparture{Status: stopTimesAndStatus.status}, nil
}

func getInternalLineType(lineType string) int {
	switch lineType {
	case "t":
		return 1
	case "bw":
		return 2
	case "b":
		return 3
	case "n":
		return 4
	}
	return 0
}

func (s *TanService) findStop(lineID, stopName string) (*stopAndStatus, error) {
	// get stops list
	resp, err := http.Get(s.rootURL + "/arrets.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var allStops []stop
	err = json.NewDecoder(resp.Body).Decode(&allStops)
	if err != nil {
		return nil, err
	}
	// find stop
	var stop *stop
	globalFoundStop := false
	globalFoundLine := false
	for _, stopLoop := range allStops {
		foundStop := strings.Contains(util.StandardizeStringLower(stopLoop.Name), stopName)
		foundLine := stopLoop.hasLineID(lineID)
		if !globalFoundStop && foundStop {
			globalFoundStop = true
		}
		if !globalFoundLine && foundLine {
			globalFoundLine = true
		}
		if foundStop && foundLine {
			stop = &stopLoop
			break
		}
	}
	result := &stopAndStatus{}
	if !globalFoundStop {
		result.status = transport.StatusUnknownStop
	} else if !globalFoundLine {
		result.status = transport.StatusUnknownLine
	} else if stop != nil {
		result.status = transport.StatusOk
		result.stop = *stop
	} else {
		result.status = transport.StatusNoService
	}
	return result, nil
}

func (s *TanService) retrieveStopTimes(lineTypeCode int, lineID string, stop stop, direction string) (*stopTimesAndStatus, error) {
	// get stop times
	var stopTimes []stopTime
	resp, err := http.Get(s.rootURL + "/tempsattente.json/" + stop.ID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var allStopTimes []stopTime
	err = json.NewDecoder(resp.Body).Decode(&allStopTimes)
	if err != nil {
		return nil, err
	}
	globalFoundDirection := false
	globalFoundLine := false
	for _, stopTime := range allStopTimes {
		foundDirection := strings.Contains(util.StandardizeStringLower(stopTime.Direction), direction)
		foundLine := stopTime.Line.Type == lineTypeCode && stopTime.Line.ID == lineID
		if !globalFoundDirection && foundDirection {
			globalFoundDirection = true
		}
		if !globalFoundLine && foundLine {
			globalFoundLine = true
		}
		if foundDirection && foundLine {
			stopTimes = append(stopTimes, stopTime)
		}
	}
	result := &stopTimesAndStatus{}
	if !globalFoundDirection {
		result.status = transport.StatusUnknownDirection
	} else if !globalFoundLine {
		result.status = transport.StatusUnknownLine
	} else {
		result.status = transport.StatusOk
		result.stopTimes = stopTimes
	}
	return result, nil
}

func stopTimesToNextDeparture(stopTimes []stopTime, max int) *transport.NextDeparture {
	if len(stopTimes) == 0 {
		return &transport.NextDeparture{Status: transport.StatusNoService}
	}
	loc, _ := time.LoadLocation("Europe/Paris")
	now := time.Now().In(loc)
	stopTimesMax := stopTimes[:util.Min(max, len(stopTimes))]
	var stopTimesResult []transport.DateTimeOrMoreThan1Hour
	for _, stopTime := range stopTimesMax {
		var dateTimeOrMoreThan1Hour *transport.DateTimeOrMoreThan1Hour
		if stopTime.Time == "Proche" {
			dateTimeOrMoreThan1Hour = &transport.DateTimeOrMoreThan1Hour{DateTime: now}
		} else if stopTime.Time == ">1h" {
			// > 1h
			dateTimeOrMoreThan1Hour = &transport.DateTimeOrMoreThan1Hour{MoreThan1Hour: true}
		} else {
			d, _ := time.Parse("04 mn", stopTime.Time)
			t := now.Add(time.Minute * time.Duration(d.Minute()))
			dateTimeOrMoreThan1Hour = &transport.DateTimeOrMoreThan1Hour{DateTime: t}
		}
		if dateTimeOrMoreThan1Hour != nil {
			stopTimesResult = append(stopTimesResult, *dateTimeOrMoreThan1Hour)
		}
	}
	return &transport.NextDeparture{Now: now, StopTimes: stopTimesResult}
}

type stopAndStatus struct {
	status string
	stop   stop
}

type stopTimesAndStatus struct {
	status    string
	stopTimes []stopTime
}

type stop struct {
	ID    string `json:"codeLieu"`
	Name  string `json:"libelle"`
	Lines []line `json:"ligne"`
}

func (s *stop) hasLineID(lineID string) bool {
	for _, line := range s.Lines {
		if line.ID == lineID {
			return true
		}
	}
	return false
}

type stopTime struct {
	Direction string `json:"terminus"`
	Time      string `json:"temps"`
	Line      line   `json:"ligne"`
}

type line struct {
	ID   string `json:"numLigne"`
	Type int    `json:"typeLigne"`
}
