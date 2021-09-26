package util

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestStandardizeStringLower(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("translate en fr 'car 'ecau'", StandardizeStringLower("traNslAte  én    fr  \"car ’éçàù\""))
}

func TestStandardizeString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("traNslAte en fr 'car ecau'", StandardizeString("traNslAte  én    fr  \"car éçàù\"", false))
}
