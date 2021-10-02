package command

import (
	"regexp"
	"strings"

	"github.com/golang/glog"

	"sms-query/pkg/config"
	"sms-query/pkg/service"
	"sms-query/pkg/service/news"
	"sms-query/pkg/service/transport"
	"sms-query/pkg/service/transport/cyclocity"
	"sms-query/pkg/service/transport/navitia"
	"sms-query/pkg/service/transport/tan"
	"sms-query/pkg/service/weather"
)

var commands []*Command = []*Command{
	searchCmd,
	translateCmd,
	// news
	newsCmd,
	leMondeCmd,
	leMondeInternationalCmd,
	leMondePoliticsCmd,
	leMondeSocietyCmd,
	leMondeEconomyCmd,
	leMondeCultureCmd,
	leMondeIdeasCmd,
	leMondePlanetCmd,
	leMondeSportCmd,
	leMondeSciencesCmd,
	leMondePixelsCmd,
	// weather
	weatherCmd,
	rainHourCmd,
	// transport
	cyclocityCmd,
	biclooCmd,
	tanCmd,
	stifCmd,
}

type Command struct {
	Keys               map[string]string
	NeedSpaceSeparator bool
	OptionsPattern     string
	DefaultValues      map[string]string
	InvalidMessageKey  string
	Run                func(from string, args map[string]string) (interface{}, error)
}

func (cmd *Command) getRegexp(locale, query string) *regexp.Regexp {
	regexp, err := regexp.Compile(cmd.GetPattern(locale))
	if err != nil {
		glog.Errorf("Error while compiling regexp: %s", err)
	}
	return regexp
}

func (cmd *Command) GetPattern(locale string) string {
	key := cmd.GetKey(locale)
	if cmd.OptionsPattern == "" {
		return key
	}
	if cmd.NeedSpaceSeparator {
		return key + " " + cmd.OptionsPattern
	}
	return key + cmd.OptionsPattern
}

func (cmd *Command) GetKey(locale string) string {
	value, ok := cmd.Keys[locale]
	if !ok {
		return cmd.Keys["en_US"]
	}
	return value
}

func (cmd *Command) IsMatch(locale, query string) bool {
	regexp := cmd.getRegexp(locale, query)
	if regexp == nil {
		return false
	}
	return regexp.MatchString(query)
}

func (cmd *Command) IsMatchKeyOnly(locale, query string) bool {
	key := cmd.GetKey(locale)
	if key == "" {
		return false
	}
	return strings.HasPrefix(query, key)
}

func (cmd *Command) GetArgs(locale, query string) map[string]string {
	regexp := cmd.getRegexp(locale, query)
	if regexp == nil {
		return nil
	}
	match := regexp.FindStringSubmatch(query)
	args := make(map[string]string)
	for i, name := range regexp.SubexpNames() {
		if i != 0 && name != "" {
			if match[i] == "" && cmd.DefaultValues != nil {
				// set default value when optional param is not provided
				args[name] = cmd.DefaultValues[name]
			} else {
				// trim spaces due to optional params declaration
				args[name] = strings.Trim(match[i], " ")
			}
		}
	}
	return args
}

// helpCmd is added into NewDispatcher
var helpCmd = &Command{
	Keys: map[string]string{"en_US": "help", "fr_FR": "aide"},
	Run: func(from string, args map[string]string) (interface{}, error) {
		defaults := config.GetInstance().GetDefaults(from)
		var result []string
		for _, cmd := range commands {
			result = append(result, cmd.GetKey(defaults.Locale))
		}
		return result, nil
	},
}

var searchCmd = &Command{
	Keys:               map[string]string{"en_US": "search", "fr_FR": "recherche"},
	NeedSpaceSeparator: true,
	OptionsPattern:     "(?P<text>.+)",
	Run: func(from string, args map[string]string) (interface{}, error) {
		defaults := config.GetInstance().GetDefaults(from)
		return service.NewSearchService().Search(defaults.Locale, args["text"])
	},
	InvalidMessageKey: "search.invalid",
}

var translateCmd = &Command{
	Keys:               map[string]string{"en_US": "translate", "fr_FR": "traduction"},
	NeedSpaceSeparator: true,
	OptionsPattern:     "((?P<source>[a-z]{2}) (?P<target>[a-z]{2}) (?P<text>.+)|(?P<textOnly>.+))",
	DefaultValues:      map[string]string{"source": "en", "target": "fr"},
	Run: func(from string, args map[string]string) (interface{}, error) {
		text := args["text"]
		if args["textOnly"] != "" {
			text = args["textOnly"]
		}
		return service.NewTranslateService().Translate(args["source"], args["target"], text)
	},
	InvalidMessageKey: "translate.invalid",
}

// news
var newsCmd = &Command{
	Keys: map[string]string{"en_US": "news", "fr_FR": "actu"},
	Run: func(from string, args map[string]string) (interface{}, error) {
		return news.Rss(config.GetInstance().GetDefaults(from).Rss, 3)
	},
}
var leMondeCmd = buildLeMondeCategoryCommand("")
var leMondeInternationalCmd = buildLeMondeCategoryCommand("international")
var leMondePoliticsCmd = buildLeMondeCategoryCommand("politique")
var leMondeSocietyCmd = buildLeMondeCategoryCommand("societe")
var leMondeEconomyCmd = buildLeMondeCategoryCommand("economie")
var leMondeCultureCmd = buildLeMondeCategoryCommand("culture")
var leMondeIdeasCmd = buildLeMondeCategoryCommand("idees")
var leMondePlanetCmd = buildLeMondeCategoryCommand("planete")
var leMondeSportCmd = buildLeMondeCategoryCommand("sport")
var leMondeSciencesCmd = buildLeMondeCategoryCommand("sciences")
var leMondePixelsCmd = buildLeMondeCategoryCommand("pixels")

func buildLeMondeCategoryCommand(category string) *Command {
	key := "lemonde"
	if category != "" {
		key += " " + category
	}
	return &Command{
		Keys: map[string]string{"en_US": key},
		Run: func(from string, args map[string]string) (interface{}, error) {
			rss := "https://www.lemonde.fr/rss/une.xml"
			if category != "" {
				rss = "https://www.lemonde.fr/" + category + "/rss_full.xml"
			}
			return news.Rss(rss, 3)
		},
	}
}

// weather
var weatherCmd = &Command{
	Keys:               map[string]string{"en_US": "weather", "fr_FR": "meteo"},
	NeedSpaceSeparator: false,
	OptionsPattern:     "( (?P<hour>h))?( (?P<city>.+))?",
	DefaultValues:      map[string]string{"hour": ""},
	Run: func(from string, args map[string]string) (interface{}, error) {
		defaults := config.GetInstance().GetDefaults(from)
		city := args["city"]
		if city == "" {
			city = defaults.City
		}
		data, err := weather.NewWeatherService().GetWeather(defaults.Locale, defaults.Country, city, args["hour"] != "")
		if err != nil {
			return nil, err
		}
		return data, nil
	},
	InvalidMessageKey: "weather.invalid",
}
var rainHourCmd = &Command{
	Keys:               map[string]string{"en_US": "rain", "fr_FR": "pluie"},
	NeedSpaceSeparator: false,
	OptionsPattern:     "( (?P<city>.+))?",
	Run: func(from string, args map[string]string) (interface{}, error) {
		defaults := config.GetInstance().GetDefaults(from)
		city := args["city"]
		if city == "" {
			city = defaults.City
		}
		return weather.NewRainHourService().GetRainHour(city)
	},
	InvalidMessageKey: "weather.rain.invalid",
}

// transport
var cyclocityCmd = buildCyclocityCmdCommand("cyclocity", "")
var biclooCmd = buildCyclocityCmdCommand("bicloo", "nantes")

func buildCyclocityCmdCommand(key, city string) *Command {
	optionsPattern := "(?P<station>.+)"
	if city == "" {
		optionsPattern = "(?P<city>(nantes|bordeaux)) " + optionsPattern
	}
	return &Command{
		Keys:               map[string]string{"en_US": key},
		NeedSpaceSeparator: true,
		OptionsPattern:     optionsPattern,
		Run: func(from string, args map[string]string) (interface{}, error) {
			contract := city
			if contract == "" {
				contract = args["city"]
			}
			return cyclocity.NewCyclocityService().GetStations(contract, args["station"], 3)
		},
		InvalidMessageKey: key + ".invalid",
	}
}

var tanCmd = buildTransportNextDepartureTimeCommand("tan")
var stifCmd = buildTransportNextDepartureTimeCommand("stif")

func buildTransportNextDepartureTimeCommand(key string) *Command {
	return &Command{
		Keys:               map[string]string{"en_US": key},
		NeedSpaceSeparator: true,
		OptionsPattern:     "((?P<lineTypeLineId>[a-z]+[1-9]+?)|(?P<lineType>.+?) (?P<lineId>.+?)) '(?P<stop>.+?)' '(?P<direction>.+?)'",
		Run: func(from string, args map[string]string) (interface{}, error) {
			lineType := args["lineType"]
			lineID := args["lineId"]
			if lineTypeLineID, ok := args["lineTypeLineId"]; ok {
				lineType = lineTypeLineID[:1]
				lineID = lineTypeLineID[1:]
			}
			return getNextDepartureTimeService(key).GetNextDepartureTime(lineType, lineID, args["stop"], args["direction"], 3)
		},
		InvalidMessageKey: "transport." + key + ".invalid",
	}
}

func getNextDepartureTimeService(key string) transport.TransportService {
	switch key {
	case "tan":
		return tan.NewTanService()
	case "stif":
		return navitia.NewNavitiaStifService()
	}
	return nil
}
