package maps

import (
	"fmt"
)

type Geocoder interface {
	Geocode(country, city string) (*Location, error)
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (l *Location) GetLatStr() string {
	return fmt.Sprintf("%.6f", l.Lat)
}

func (l *Location) GetLngStr() string {
	return fmt.Sprintf("%.6f", l.Lat)
}
