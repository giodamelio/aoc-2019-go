package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/giodamelio/aoc-2020-go/intcode"
	"github.com/stretchr/testify/assert"
)

var discardLog = log.New(ioutil.Discard, "", 0)

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
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int(3232358), part1(discardLog, parsedInput))
}

func TestPart2(t *testing.T) {
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	assert.Equal(t, int(4845669), part2(discardLog, parsedInput))
}
