package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFreeName(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	name, err := GetFreeName(ctx)
	assert.NotEqual(t, "!", name)
	assert.Nil(t, err, "should generate name without error")
}

func TestCheckNameExists(t *testing.T) {
	ctx, inst := getTestContext(t)
	defer inst.Close()
	alias, _ := makeTestAlias(ctx)
	result2 := checkNameExists(ctx, "new_one_iota")
	assert.Equal(t, result2, false)
	result := checkNameExists(ctx, strings.Split(alias.Alias, "@")[0])
	assert.Equal(t, result, true)
}
