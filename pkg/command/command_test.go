package command

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestGetCommand(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("help", helpCmd.GetKey("en_us"))

	assert.Equal("bicloo", biclooCmd.GetKey("en_us"))
	assert.Equal("bicloo", biclooCmd.GetKey("fr_FR"))

	assert.Equal(searchCmd.GetKey("en_us")+" "+searchCmd.OptionsPattern, searchCmd.GetPattern("en_us"))
	assert.Equal(searchCmd.GetKey("fr_FR")+" "+searchCmd.OptionsPattern, searchCmd.GetPattern("fr_FR"))
}

func TestGetArgs(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(map[string]string{"text": "toto tata"},
		searchCmd.GetArgs("en_us", "search toto tata"))

	assert.Equal(map[string]string{
		"source":   "en",
		"target":   "fr",
		"text":     "car",
		"textOnly": "",
	},
		translateCmd.GetArgs("en_us", "translate en fr car"))
	assert.Equal(map[string]string{
		"source":   "en",
		"target":   "fr",
		"text":     "",
		"textOnly": "car",
	},
		translateCmd.GetArgs("en_us", "translate car"))
}

func TestIsMatch(t *testing.T) {
	assert := assert.New(t)

	assert.True(helpCmd.IsMatch("en_us", "help"))

	assert.False(searchCmd.IsMatch("en_us", "search"))
	assert.True(searchCmd.IsMatch("en_us", "search toto"))
	assert.True(searchCmd.IsMatch("en_us", "search toto tata"))

	assert.True(weatherCmd.IsMatch("en_us", "weather"))
	assert.True(weatherCmd.IsMatch("en_us", "weather h"))
	assert.True(weatherCmd.IsMatch("en_us", "weather nantes"))
	assert.True(weatherCmd.IsMatch("en_us", "weather h nantes"))
	assert.True(rainHourCmd.IsMatch("en_us", "rain"))
	assert.True(rainHourCmd.IsMatch("en_us", "rain nantes"))

	assert.False(cyclocityCmd.IsMatch("en_us", "cyclocity"))
	assert.True(cyclocityCmd.IsMatch("en_us", "cyclocity nantes gare"))
	assert.True(cyclocityCmd.IsMatch("en_us", "cyclocity bordeaux gare"))
	assert.False(cyclocityCmd.IsMatch("en_us", "cyclocity paris gare"))
	assert.False(biclooCmd.IsMatch("en_us", "bicloo"))
	assert.True(biclooCmd.IsMatch("en_us", "bicloo gare maritime"))

	assert.False(tanCmd.IsMatch("en_us", "tan"))
	assert.True(tanCmd.IsMatch("en_us", "tan t 1 'stop' 'direction'"))
	assert.True(tanCmd.IsMatch("en_us", "tan tr 123 'stop' 'direction'"))
	assert.True(tanCmd.IsMatch("en_us", "tan t1 'stop' 'direction'"))
	assert.True(tanCmd.IsMatch("en_us", "tan tr123 'stop' 'direction'"))
}
