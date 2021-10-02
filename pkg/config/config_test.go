package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := InitWithPath("../../config/config.yml")
	if err != nil {
		t.Fatal(err)
	}
	config := GetInstance()

	assert := assert.New(t)
	assert.Equal("jcdecaux", config.Cyclocity.ID)
	assert.Empty(config.Cyclocity.APIKey)
}
