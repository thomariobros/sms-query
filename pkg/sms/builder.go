package sms

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/golang/glog"
	"github.com/goodsign/monday"
	"github.com/mmcdole/gofeed"

	"main/pkg/config"
	"main/pkg/errors"
	"main/pkg/i18n"
	"main/pkg/service/search"
	"main/pkg/service/transport"
	"main/pkg/service/transport/cyclocity"
	"main/pkg/service/weather"
	"main/pkg/util"
)

func Build(from string, query string, data interface{}, err error) string {
	result := ""
	preferences := config.GetInstance().GetPreferences(from)
	if preferences != nil && preferences.DisplayQuery {
		result += query + "\n"
	}
	if err == nil && data != nil {
		result += buildContent(preferences, data)
	} else {
		result += buildError(from, err)
	}
	return strings.ToLower(result)
}

func buildError(from string, err error) string {
	if _, ok := err.(errors.CustomError); ok {
		return i18n.GetTranslationForPhoneNumber(from, "error.wrong_input")
	}
	glog.Errorf("Sending error message to %s because of %s", from, err)
	return i18n.GetTranslationForPhoneNumber(from, "error")
}

func buildContent(preferences *config.Preferences, data interface{}) string {
	locale := preferences.Defaults.Locale
	location := preferences.Defaults.GetLocation()
	switch data.(type) {
	case string:
		return BuildString(locale, data.(string))
	case []string:
		return BuildArrayStrings(locale, data.([]string))
	case []cyclocity.CyclocityStation:
		return BuildCyclocityStations(locale, data.([]cyclocity.CyclocityStation))
	case *search.SearchResult:
		return BuildSearchResult(locale, data.(*search.SearchResult))
	case []gofeed.Item:
		return BuildRss(locale, data.([]gofeed.Item))
	case *weather.RainHour:
		return BuildRainHour(locale, data.(*weather.RainHour))
	case *weather.WeatherData:
		return BuildWeather(locale, location, data.(*weather.WeatherData))
	case *transport.NextDeparture:
		return BuildNextDeparture(locale, data.(*transport.NextDeparture))
	default:
		return ""
	}
}

func BuildString(locale string, str string) string {
	return str
}

func BuildArrayStrings(locale string, values []string) string {
	params := map[string]interface{}{
		"values": values,
	}
	template := "{{ range $index, $value := $.values }}" +
		"{{ if gt $index 0 }}" +
		"\n" +
		"{{ end }}" +
		"{{ $value }}" +
		"{{ end }}"
	return buildFromTemplate(template, params)
}

func BuildCyclocityStations(locale string, stations []cyclocity.CyclocityStation) string {
	params := map[string]interface{}{
		"i18n": map[string]string{
			"common.separator.colon":        i18n.GetTranslation(locale, "common.separator.colon"),
			"cyclocity.availableBikes":      i18n.GetTranslation(locale, "cyclocity.available_bikes"),
			"cyclocity.availableBikeStands": i18n.GetTranslation(locale, "cyclocity.available_bike_stands"),
		},
		"stations": stations,
	}
	template := "{{ range $index, $station := $.stations }}" +
		"{{ if gt $index 0 }}" +
		"\n" +
		"{{ end }}" +
		"{{ $station.Name }}, {{ $station.Address }}\n" +
		"{{ index $.i18n \"cyclocity.availableBikes\" }}{{ index $.i18n \"common.separator.colon\" }}{{ $station.AvailableBikes }}\n" +
		"{{ index $.i18n \"cyclocity.availableBikeStands\" }}{{ index $.i18n \"common.separator.colon\" }}{{ $station.AvailableBikeStands }}" +
		"{{ end }}"
	return buildFromTemplate(template, params)
}

func BuildSearchResult(locale string, item *search.SearchResult) string {
	if item == nil || len((*item).Text) == 0 {
		return i18n.GetTranslation(locale, "search.no_data")
	}
	params := map[string]interface{}{
		"i18n": map[string]string{
			"common.separator.hyphen": i18n.GetTranslation(locale, "common.separator.hyphen"),
		},
		"item": item,
	}
	template := "{{ $.item.Source }}" +
		"{{ index $.i18n \"common.separator.hyphen\" }}" +
		"{{ $.item.Text }}"
	return buildFromTemplate(template, params)
}

func BuildRss(locale string, items []gofeed.Item) string {
	params := map[string]interface{}{
		"i18n": map[string]string{
			"common.separator.hyphen": i18n.GetTranslation(locale, "common.separator.hyphen"),
		},
		"items": items,
	}
	template := "{{ range $index, $item := $.items }}" +
		"{{ if gt $index 0 }}" +
		"\n" +
		"{{ end }}" +
		"{{ inc $index }}{{ index $.i18n \"common.separator.hyphen\" }}{{ $item.Title }}" +
		"{{ if gt (len $item.Description) 0 }}" +
		"{{ index $.i18n \"common.separator.hyphen\" }}{{ $item.Description }}" +
		"{{ end }}" +
		"{{ end }}"
	return buildFromTemplate(template, params)
}

func BuildRainHour(locale string, rainHour *weather.RainHour) string {
	if rainHour == nil || !rainHour.IsAvailable {
		return i18n.GetTranslation(locale, "weather.rain.not_available")
	}
	if !rainHour.HasData || len(rainHour.Data) == 0 {
		return i18n.GetTranslation(locale, "weather.rain.no_data")
	}
	// build items
	var items []interface{}
	t, _ := time.Parse("200601021504", rainHour.Deadline)
	for _, data := range rainHour.Data {
		items = append(items, map[string]interface{}{
			"time":  t.Format("15:04"),
			"level": data.Level,
		})
		t = t.Add(time.Minute * time.Duration(5))
	}
	funcs := template.FuncMap{
		"level2I18n": func(level int) string {
			return fmt.Sprintf("level.%v", level)
		},
	}
	params := map[string]interface{}{
		"i18n": map[string]string{
			"common.separator.colon": i18n.GetTranslation(locale, "common.separator.colon"),
			"level.0":                i18n.GetTranslation(locale, "weather.rain.no_data_available"),
			"level.1":                i18n.GetTranslation(locale, "weather.rain.no_rainfall"),
			"level.2":                i18n.GetTranslation(locale, "weather.rain.low_rainfall"),
			"level.3":                i18n.GetTranslation(locale, "weather.rain.moderate_rainfall"),
			"level.4":                i18n.GetTranslation(locale, "weather.rain.heavy_rainfall"),
		},
		"items": items,
	}
	template := "{{ range $index, $item := $.items }}" +
		"{{ if gt $index 0 }}" +
		"\n" +
		"{{ end }}" +
		"{{ $item.time }}{{ index $.i18n \"common.separator.colon\" }}" +
		"{{ if and (ge $item.level 1) (le $item.level 4)}}" +
		"{{ index $.i18n (level2I18n $item.level) }}" +
		"{{ else }}" +
		"{{ index $.i18n \"level.0\" }}" +
		"{{ end }}" +
		"{{ end }}"
	return buildFromTemplateFuncs(template, &funcs, params)
}

func BuildWeather(locale string, location *time.Location, data *weather.WeatherData) string {
	if data == nil || data.Data == nil || len(data.Data.Data) == 0 {
		return i18n.GetTranslation(locale, "weather.no_data")
	}
	// build items
	var items []interface{}
	size := 5
	format := "Mon"
	if data.Hour {
		size = 12
		format = "Mon 15h"
	}
	for _, data := range data.Data.Data[:util.Min(size, len(data.Data.Data))] {
		items = append(items, map[string]interface{}{
			"time":           monday.Format(time.Unix(0, data.Time*int64(time.Second)).In(location), format, i18n.GetMondayLocale(locale)),
			"summary":        strings.Replace(data.Summary, ".", "", -1),
			"temperature":    util.Round(data.Temperature),
			"temperatureMin": util.Round(data.TemperatureMin),
			"temperatureMax": util.Round(data.TemperatureMax),
		})
	}

	params := map[string]interface{}{
		"i18n": map[string]string{
			"common.separator.colon":  i18n.GetTranslation(locale, "common.separator.colon"),
			"common.separator.comma":  i18n.GetTranslation(locale, "common.separator.comma"),
			"common.separator.hyphen": i18n.GetTranslation(locale, "common.separator.hyphen"),
		},
		"items": items,
	}
	template := "{{ range $index, $item := $.items }}" +
		"{{ if gt $index 0 }}" +
		"\n" +
		"{{ end }}" +
		"{{ $item.time }}{{ index $.i18n \"common.separator.colon\" }}{{ $item.summary }}" +
		"{{ if and (gt $item.temperatureMin 0) (gt $item.temperatureMax 0) }}" +
		"{{ index $.i18n \"common.separator.comma\" }}{{ $item.temperatureMin }}°C" +
		"{{ index $.i18n \"common.separator.hyphen\" }}{{ $item.temperatureMax }}°C" +
		"{{ else }}" +
		"{{ if gt $item.temperature 0 }}" +
		"{{ index $.i18n \"common.separator.comma\" }}{{ $item.temperature }}°C" +
		"{{ end }}" +
		"{{ end }}" +
		"{{ end }}"
	return buildFromTemplate(template, params)
}

func BuildNextDeparture(locale string, nextDeparture *transport.NextDeparture) string {
	if nextDeparture == nil {
		return i18n.GetTranslation(locale, "transport.no_service")
	}
	if nextDeparture.Status != transport.StatusOk {
		switch nextDeparture.Status {
		case transport.StatusUnknownLine:
			return i18n.GetTranslation(locale, "transport.unknown_line")
		case transport.StatusUnknownStop:
			return i18n.GetTranslation(locale, "transport.unknown_stop")
		case transport.StatusUnknownDirection:
			return i18n.GetTranslation(locale, "transport.unknown_direction")
		case transport.StatusNoService:
			return i18n.GetTranslation(locale, "transport.no_service")
		case transport.StatusError:
			return i18n.GetTranslation(locale, "transport.no_service")
		}
	}
	if len(nextDeparture.StopTimes) == 0 {
		return i18n.GetTranslation(locale, "transport.no_service")
	}
	funcs := template.FuncMap{
		"formatTime": func(dateTime time.Time) string {
			return dateTime.Format("15:04")
		},
		"formatDuration": func(dateTime time.Time) string {
			duration := dateTime.Sub(nextDeparture.Now)
			if duration.Minutes() <= 60 {
				if duration.Minutes() < 1 {
					return i18n.GetTranslation(locale, "transport.less_than_1_minute")
				}
				return fmt.Sprintf("%.f%s", duration.Minutes(), i18n.GetTranslation(locale, "common.minutes.short"))
			}
			return i18n.GetTranslation(locale, "transport.moreThan1Hour")
		},
	}
	params := map[string]interface{}{
		"i18n": map[string]string{
			"common.separator.colon":     i18n.GetTranslation(locale, "common.separator.colon"),
			"common.separator.hyphen":    i18n.GetTranslation(locale, "common.separator.hyphen"),
			"transport.at":               i18n.GetTranslation(locale, "transport.at"),
			"transport.more_than_1_hour": i18n.GetTranslation(locale, "transport.more_than_1_hour"),
		},
		"nextDeparture": nextDeparture,
	}
	template := "{{ index $.i18n \"transport.at\" }} {{ formatTime $.nextDeparture.Now }}{{ index $.i18n \"common.separator.colon\" }}\n" +
		"{{ range $index, $stopTime := $.nextDeparture.StopTimes }}" +
		"{{ if gt $index 0 }}" +
		"\n" +
		"{{ end }}" +
		"{{ if $stopTime.MoreThan1Hour }}" +
		"{{ index $.i18n \"transport.more_than_1_hour\" }}" +
		"{{ else }}" +
		"{{ formatTime $stopTime.DateTime }}{{ index $.i18n \"common.separator.hyphen\" }}{{ formatDuration $stopTime.DateTime }}" +
		"{{ end }}" +
		"{{ end }}"
	return buildFromTemplateFuncs(template, &funcs, params)
}

func buildFromTemplate(content string, params map[string]interface{}) string {
	return buildFromTemplateFuncs(content, nil, params)
}

func buildFromTemplateFuncs(content string, funcs *template.FuncMap, params map[string]interface{}) string {
	inc := func(i int) int {
		return i + 1
	}
	if funcs == nil {
		funcs = &template.FuncMap{}
	}
	(*funcs)["inc"] = inc
	tmpl, err := template.New("template").Funcs(*funcs).Parse(content)
	if err != nil {
		glog.Errorf("error while parsing template: %s", err)
		return ""
	}
	var bytes bytes.Buffer
	tmpl.Execute(&bytes, params)
	return bytes.String()
}
