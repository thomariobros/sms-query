package navitia

import (
	"encoding/json"
	net_http "net/http"
	"net/url"
	"strings"
	"time"

	"main/pkg/config"
	"main/pkg/http"
	"main/pkg/service/transport"
	"main/pkg/util"
)

// http://doc.navitia.io/

const (
	rootURL = "https://api.navitia.io/v1/coverage"
)

type NavitiaService struct {
	rootURL  string
	coverage string
	location time.Location
}

func NewNavitiaStifService() *NavitiaService {
	loc, _ := time.LoadLocation("Europe/Paris")
	return &NavitiaService{
		rootURL:  rootURL,
		coverage: "fr-idf",
		location: *loc,
	}
}

func NewNavitiaStifServiceRootURL(rootURL string) *NavitiaService {
	loc, _ := time.LoadLocation("Europe/Paris")
	return &NavitiaService{
		rootURL:  rootURL,
		coverage: "fr-idf",
		location: *loc,
	}
}

func getInternalLineType(lineType string) string {
	switch lineType {
	case "m":
		return "Metro"
	case "t":
		return "Tram"
	case "tr":
		return "Train"
	case "b":
		return "Bus"
	}
	return ""
}

// GetNextDepartureTime get next departure time through navitia API
func (s *NavitiaService) GetNextDepartureTime(lineType, lineID, stopName, direction string, max int) (*transport.NextDeparture, error) {
	// find line
	internalLineType := getInternalLineType(lineType)
	if internalLineType == "" {
		return &transport.NextDeparture{
			Status: transport.StatusUnknownLine,
		}, nil
	}
	status, lw, err := s.findLines(internalLineType, lineID, direction)
	if err != nil {
		return nil, err
	}
	if status != transport.StatusOk {
		return &transport.NextDeparture{
			Status: status,
		}, nil
	}
	// find stops
	stops, err := s.findStops(stopName)
	if err != nil {
		return nil, err
	}
	if len(stops) == 0 {
		return &transport.NextDeparture{
			Status: transport.StatusUnknownStop,
		}, nil
	}
	// loop stops and return first match
	for _, stop := range stops {
		departures, err := s.retrieveDepartures(stop, *lw, direction)
		if err != nil {
			return nil, err
		}
		nextDeparture := s.calculateNextDeparture(departures, max)
		return &nextDeparture, nil
	}
	return &transport.NextDeparture{
		Status: transport.StatusNoService,
	}, nil
}

func (s *NavitiaService) getRootURL() string {
	return s.rootURL + "/" + s.coverage
}

func (s *NavitiaService) get(url string) (*net_http.Response, error) {
	headers := map[string]string{"Authorization": config.GetInstance().Navitia.APIKey}
	return http.Do("GET", url, headers, nil)
}

// http://doc.navitia.io/#autocomplete-on-public-transport-objects
func (s *NavitiaService) findStops(stopName string) ([]stop, error) {
	params := "q=" + url.QueryEscape(stopName) + "&type[]=stop_area"
	resp, err := s.get(s.getRootURL() + "/pt_objects" + "?" + params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var wrapper stopsWrapper
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper.StopWrappers, nil
}

// http://doc.navitia.io/#autocomplete-on-public-transport-objects
func (s *NavitiaService) findLines(lineType, lineID, direction string) (string, *lineWrapper, error) {
	params := "q=" + url.QueryEscape(lineType+lineID) + "&type[]=line"
	resp, err := s.get(s.getRootURL() + "/pt_objects" + "?" + params)
	if err != nil {
		return transport.StatusError, nil, err
	}
	defer resp.Body.Close()
	var wrapper linesWrapper
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return transport.StatusError, nil, err
	}
	var lw *lineWrapper
	globalFoundDirection := false
	globalFoundLine := false
	if len(wrapper.LineWrappers) > 0 {
		for _, lwLoop := range wrapper.LineWrappers {
			foundDirection := lwLoop.Line.hasDirection(direction)
			foundLine := lwLoop.Line.isEquals(lineType, lineID)
			if !globalFoundDirection && foundDirection {
				globalFoundDirection = true
			}
			if !globalFoundLine && foundLine {
				globalFoundLine = true
			}
			if foundDirection && foundLine {
				lw = &lwLoop
				break
			}
		}
	} else {
		globalFoundLine = false
	}
	if !globalFoundLine {
		return transport.StatusUnknownLine, nil, nil
	} else if !globalFoundDirection {
		return transport.StatusUnknownDirection, nil, nil
	} else {
		return transport.StatusOk, lw, nil
	}
}

// http://doc.navitia.io/#next-departures-and-arrivals
func (s *NavitiaService) retrieveDepartures(stop stop, lineWrapper lineWrapper, destination string) ([]departure, error) {
	now := time.Now().In(&s.location)
	params := "from_datetime=" + now.Format("20060102'T'150405") +
		"&count=100" // http://doc.navitia.io/#paging
	resp, err := s.get(s.getRootURL() + "/stop_areas/" + stop.ID + "/departures" + "?" + params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var dw departuresWrapper
	err = json.NewDecoder(resp.Body).Decode(&dw)
	if err != nil {
		return nil, err
	}

	var departures []departure
	for _, departure := range dw.Departures {
		// ' - ' is origin / destination separator
		if strings.Contains(util.StandardizeStringLower(departure.Route.Name), " - "+destination) &&
			strings.EqualFold(departure.Route.Line.ID, lineWrapper.ID) {
			departures = append(departures, departure)
		}
	}
	return departures, nil
}

func (s *NavitiaService) calculateNextDeparture(departures []departure, max int) transport.NextDeparture {
	var stopTimesResult []transport.DateTimeOrMoreThan1Hour
	if len(departures) > 0 {
		departuresMax := departures[:util.Min(max, len(departures))]
		for _, departure := range departuresMax {
			stopTimesResult = append(stopTimesResult, transport.DateTimeOrMoreThan1Hour{
				DateTime: departure.StopTime.getDepartureDateTime().In(&s.location),
			})
		}
	}
	return transport.NextDeparture{
		Status:    transport.StatusOk,
		Now:       time.Now().In(&s.location),
		StopTimes: stopTimesResult,
	}
}

type object struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type line struct {
	object
	Code          string         `json:"code"`
	PhysicalModes []physicalMode `json:"physical_modes"`
	Routes        []object       `json:"routes"`
}

func (l *line) isEquals(t, code string) bool {
	return strings.EqualFold(l.Code, code) && l.hasPhysicalModes(t)
}

func (l *line) hasPhysicalModes(code string) bool {
	for _, physicalMode := range l.PhysicalModes {
		if strings.EqualFold(strings.Split(physicalMode.ID, ":")[1], code) {
			return true
		}
	}
	return false
}

func (l *line) hasDirection(direction string) bool {
	for _, route := range l.Routes {
		if strings.Contains(util.StandardizeStringLower(route.Name), direction) {
			return true
		}
	}
	return false
}

type lineWrapper struct {
	object
	Line line `json:"line"`
}

type linesWrapper struct {
	LineWrappers []lineWrapper `json:"pt_objects"`
}

type route struct {
	object
	Line line `json:"line"`
}

type physicalMode struct {
	ID string `json:"id"`
}

type stop struct {
	object
}

type stopsWrapper struct {
	StopWrappers []stop `json:"pt_objects"`
}

type stopTime struct {
	DepartureDateTimeString string `json:"departure_date_time"`
}

func (s *stopTime) getDepartureDateTime() time.Time {
	t, _ := time.Parse("20060102'T'150405", s.DepartureDateTimeString)
	return t
}

type departure struct {
	Route    route    `json:"route"`
	StopTime stopTime `json:"stop_date_time"`
}

type departuresWrapper struct {
	Departures []departure `json:"departures"`
}
