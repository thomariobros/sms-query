package config

import (
	"time"

	"main/pkg/util"
)

type Preferences struct {
	PhoneNumber  string   `yaml:"phoneNumber"`
	Aliases      []Alias  `yaml:"aliases"`
	DisplayQuery bool     `yaml:"displayQuery"`
	Defaults     Defaults `yaml:"defaults"`
}

func (p *Preferences) FindAlias(query string) string {
	standardizedQuery := util.StandardizeStringLower(query)
	for _, alias := range p.Aliases {
		if alias.Key == standardizedQuery {
			return alias.Value
		}
	}
	return ""
}

type Alias struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Defaults struct {
	TimeZoneID   string `yaml:"timeZoneId"`
	Locale       string `yaml:"locale"`
	Country      string `yaml:"country"`
	City         string `yaml:"city"`
	Rss          string `yaml:"rss"`
	dateTimeZone *time.Location
}

func (d *Defaults) GetLocation() *time.Location {
	if d.dateTimeZone == nil {
		d.dateTimeZone, _ = time.LoadLocation(d.TimeZoneID)
	}
	return d.dateTimeZone
}

func (d *Defaults) GetLanguage() string {
	return d.Locale[:2]
}
