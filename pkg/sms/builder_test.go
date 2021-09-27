package sms

import (
	"testing"
	"time"

	"github.com/goodsign/monday"
	forecast "github.com/mlbright/forecast/v2"
	"github.com/mmcdole/gofeed"

	"github.com/stretchr/testify/assert"

	"main/pkg/i18n"
	"main/pkg/service/search"
	"main/pkg/service/transport"
	"main/pkg/service/transport/cyclocity"
	"main/pkg/service/weather"
)

func setup() {
	i18n.LoadWithRootPath("../../config")
}

func TestBuildArrayStrings(t *testing.T) {
	setup()

	assert := assert.New(t)

	values := []string{"value1", "value2"}
	assert.Equal("value1\nvalue2", BuildArrayStrings("fr_FR", values))
}

func TestBuildCyclocityStations(t *testing.T) {
	setup()

	assert := assert.New(t)

	stations := []cyclocity.CyclocityStation{
		cyclocity.CyclocityStation{
			Name:                "name1",
			Address:             "address1",
			AvailableBikes:      1,
			AvailableBikeStands: 2,
		},
		cyclocity.CyclocityStation{
			Name:                "name2",
			Address:             "address2",
			AvailableBikes:      3,
			AvailableBikeStands: 4,
		},
	}
	assert.Equal("name1, address1\n"+
		"vélos disponibles : 1\n"+
		"points d'attache disponibles : 2\n"+
		"name2, address2\n"+
		"vélos disponibles : 3\n"+
		"points d'attache disponibles : 4",
		BuildCyclocityStations("fr_FR", stations))
}

func TestBuildSearchResult(t *testing.T) {
	setup()

	assert := assert.New(t)

	searchResult := search.SearchResult{
		Source: "www.nantes.fr",
		Text:   "Ville de Nantes",
	}
	assert.Equal("www.nantes.fr - Ville de Nantes",
		BuildSearchResult("fr_FR", &searchResult))
}

func TestBuildRss(t *testing.T) {
	setup()

	assert := assert.New(t)

	items := []gofeed.Item{
		gofeed.Item{
			Title:       "Title1",
			Description: "Description1",
		},
		gofeed.Item{
			Title: "Title2",
		},
	}
	assert.Equal("1 - Title1 - Description1\n"+
		"2 - Title2",
		BuildRss("fr_FR", items))
}

func TestBuildRainHour(t *testing.T) {
	setup()

	assert := assert.New(t)

	rainHour := weather.RainHour{
		IsAvailable: true,
		HasData:     true,
		Deadline:    "201707201815",
		Data: []weather.RainHourData{
			weather.RainHourData{Level: 0},
			weather.RainHourData{Level: 1},
			weather.RainHourData{Level: 2},
			weather.RainHourData{Level: 3},
			weather.RainHourData{Level: 4},
			weather.RainHourData{Level: 5},
			weather.RainHourData{Level: 2},
		},
	}
	assert.Equal("18:15 : données indisponibles\n"+
		"18:20 : pas de précipitations\n"+
		"18:25 : précipitations faibles\n"+
		"18:30 : précipitations modérées\n"+
		"18:35 : précipitations fortes\n"+
		"18:40 : données indisponibles\n"+
		"18:45 : précipitations faibles",
		BuildRainHour("fr_FR", &rainHour))
}

func TestBuildWeatherDataBlock(t *testing.T) {
	setup()

	assert := assert.New(t)

	loc, _ := time.LoadLocation("Europe/Paris")
	now := time.Now().In(loc)
	nowPlus1Hour := now.Add(time.Hour)

	data := weather.WeatherData{
		Hour: true,
		Data: &forecast.DataBlock{},
	}
	assert.Equal("données indisponibles",
		BuildWeather("fr_FR", loc, &data))

	data = weather.WeatherData{
		Hour: true,
		Data: &forecast.DataBlock{
			Data: []forecast.DataPoint{
				forecast.DataPoint{
					Time:        now.Unix(),
					Summary:     "Sunny",
					Temperature: 32.1,
				},
				forecast.DataPoint{
					Time:        nowPlus1Hour.Unix(),
					Summary:     "Clear",
					Temperature: 28.7,
				},
			},
		},
	}
	assert.Equal(monday.Format(now, "Mon 15h", monday.LocaleFrFR)+" : Sunny, 32°C\n"+
		monday.Format(nowPlus1Hour, "Mon 15h", monday.LocaleFrFR)+" : Clear, 29°C",
		BuildWeather("fr_FR", loc, &data))
}

func TestBuildNextDeparture(t *testing.T) {
	setup()

	assert := assert.New(t)

	nextDeparture := transport.NextDeparture{}
	assert.Equal("pas de service",
		BuildNextDeparture("fr_FR", &nextDeparture))

	nextDeparture = transport.NextDeparture{Status: transport.StatusUnknownLine}
	assert.Equal("ligne inconnue",
		BuildNextDeparture("fr_FR", &nextDeparture))

	t1 := time.Now()
	t2 := t1.Add(time.Second * time.Duration(30))
	t3 := t1.Add(time.Minute * time.Duration(5))
	nextDeparture = transport.NextDeparture{
		Status: transport.StatusOk,
		Now:    t1,
		StopTimes: []transport.DateTimeOrMoreThan1Hour{
			transport.DateTimeOrMoreThan1Hour{
				MoreThan1Hour: false,
				DateTime:      t2,
			},
			transport.DateTimeOrMoreThan1Hour{
				MoreThan1Hour: false,
				DateTime:      t3,
			},
			transport.DateTimeOrMoreThan1Hour{
				MoreThan1Hour: true,
			},
		},
	}
	assert.Equal("à "+t1.Format("15:04")+" : \n"+
		t2.Format("15:04")+" - < 1min\n"+
		t3.Format("15:04")+" - 5min\n"+
		"> 1h",
		BuildNextDeparture("fr_FR", &nextDeparture))
}
