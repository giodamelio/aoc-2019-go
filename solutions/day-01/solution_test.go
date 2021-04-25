package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var discardLog *log.Logger = log.New(ioutil.Discard, "", 0)

func TestParseInput(t *testing.T) {
	parsedInput, err := parseInput("1\n2\n3\n")

	assert.Nil(t, err, "Parsing returned an error")
	assert.Equal(t, []int64{1, 2, 3}, parsedInput)
}

func TestParseInputNonInteger(t *testing.T) {
	parsedInput, err := parseInput("1\nhaha\n3\n")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "strconv.ParseInt: parsing \"haha\": invalid syntax")
	assert.Nil(t, parsedInput)
}

func TestCalculateMass(t *testing.T) {
	assert.Equal(t, int64(2), calculateMass(12))
	assert.Equal(t, int64(2), calculateMass(14))
	assert.Equal(t, int64(654), calculateMass(1969))
	assert.Equal(t, int64(33583), calculateMass(100756))
}

func TestCalculateMassWithFuel(t *testing.T) {
	assert.Equal(t, int64(2), calculateMassWithFuel(14))
	assert.Equal(t, int64(966), calculateMassWithFuel(1969))
	assert.Equal(t, int64(50346), calculateMassWithFuel(100756))
}

func TestPart1(t *testing.T) {
	parsedInput, err := parseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int64(3232358), part1(discardLog, parsedInput))
}

func TestPart2(t *testing.T) {
	parsedInput, err := parseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int64(4845669), part2(discardLog, parsedInput))
}
