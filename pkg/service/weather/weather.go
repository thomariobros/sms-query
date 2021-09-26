package weather

import (
	forecast "github.com/mlbright/forecast/v2"
	"main/pkg/config"
	"main/pkg/service"
)

type WeatherService struct {
}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

// GetWeather get weather through dark sky API (formerly forecast)
func (s *WeatherService) GetWeather(locale, country, city string, hour bool) (*WeatherData, error) {
	location, err := service.NewMapsService().Geocode(country, city)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, nil
	}

	weatherConfig := config.GetInstance().Weather
	data, err := forecast.Get(weatherConfig.APIKey, location.GetLatStr(), location.GetLngStr(), "now", forecast.SI, getLang(locale))
	if err != nil {
		return nil, err
	}

	if hour {
		return &WeatherData{
			Hour: true,
			Data: &data.Hourly,
		}, nil
	}

	return &WeatherData{
		Hour: false,
		Data: &data.Daily,
	}, nil
}

func getLang(locale string) forecast.Lang {
	switch locale {
	case "fr_FR":
		return forecast.French
	default:
		return forecast.English
	}
}

type WeatherData struct {
	Hour bool
	Data *forecast.DataBlock
}
