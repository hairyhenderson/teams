package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHumanDuration(t *testing.T) {
	assert.Equal(t, "Less than a second", HumanDuration(200*time.Millisecond))

	assert.Equal(t, "15 seconds", HumanDuration(15*time.Second))

	assert.Equal(t, "About a minute", HumanDuration(90*time.Second))

	assert.Equal(t, "5 minutes", HumanDuration(5*time.Minute))

	assert.Equal(t, "About an hour", HumanDuration(75*time.Minute))

	assert.Equal(t, "27 hours", HumanDuration(27*time.Hour))

	assert.Equal(t, "3 days", HumanDuration(76*time.Hour))
}
