package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoString(t *testing.T) {
	r := Repo{
		Org:  "foo",
		Name: "bar",
	}
	assert.Equal(t, "foo/bar", r.String())
}
