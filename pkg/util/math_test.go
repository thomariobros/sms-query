package util

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestMin(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, Min(1, 2))
	assert.Equal(2, Min(3, 2))
}

func TestRound(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, Round(1.1))
	assert.Equal(2, Round(1.5))
	assert.Equal(2, Round(1.8))
}
