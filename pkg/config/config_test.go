package config

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestInit(t *testing.T) {
	err := InitWithRootPath("../../config")
	if err != nil {
		t.Fatal(err)
	}
	config := GetInstance()

	assert := assert.New(t)
	assert.Equal("jcdecaux", config.Cyclocity.ID)
	assert.NotEmpty(config.Cyclocity.APIKey)
}
