package main

import (
	"testing"

	"github.com/giodamelio/aoc-2020-go/intcode"
	"github.com/stretchr/testify/assert"
)

func TestPart1(t *testing.T) {
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int64(12490719), part1(discardLogger, parsedInput))
}

func TestPart2(t *testing.T) {
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int64(2003), part2(discardLogger, parsedInput))
}
