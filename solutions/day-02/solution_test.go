package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	parsedInput, err := parseInput("1,2,3")

	assert.Nil(t, err, "Parsing returned an error")
	assert.Equal(t, []int{1, 2, 3}, parsedInput)
}

func TestParseInputNonInteger(t *testing.T) {
	parsedInput, err := parseInput("1,haha,3")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "strconv.Atoi: parsing \"haha\": invalid syntax")
	assert.Nil(t, parsedInput)
}

func TestPart1(t *testing.T) {
	parsedInput, err := parseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, 12490719, part1(discardLogger, parsedInput))
}

func TestPart2(t *testing.T) {
	parsedInput, err := parseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, 2003, part2(discardLogger, parsedInput))
}
