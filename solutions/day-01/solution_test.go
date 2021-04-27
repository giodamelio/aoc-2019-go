package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var discardLog = log.New(ioutil.Discard, "", 0)

func TestParseInput(t *testing.T) {
	parsedInput, err := parseInput("1\n2\n3\n")

	assert.Nil(t, err, "Parsing returned an error")
	assert.Equal(t, []int{1, 2, 3}, parsedInput)
}

func TestParseInputNonInteger(t *testing.T) {
	parsedInput, err := parseInput("1\nhaha\n3\n")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "strconv.Atoi: parsing \"haha\": invalid syntax")
	assert.Nil(t, parsedInput)
}

func TestCalculateMass(t *testing.T) {
	assert.Equal(t, int(2), calculateMass(12))
	assert.Equal(t, int(2), calculateMass(14))
	assert.Equal(t, int(654), calculateMass(1969))
	assert.Equal(t, int(33583), calculateMass(100756))
}

func TestCalculateMassWithFuel(t *testing.T) {
	assert.Equal(t, int(2), calculateMassWithFuel(14))
	assert.Equal(t, int(966), calculateMassWithFuel(1969))
	assert.Equal(t, int(50346), calculateMassWithFuel(100756))
}

func TestPart1(t *testing.T) {
	parsedInput, err := parseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int(3232358), part1(discardLog, parsedInput))
}

func TestPart2(t *testing.T) {
	parsedInput, err := parseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int(4845669), part2(discardLog, parsedInput))
}
