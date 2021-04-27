package main

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

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

	output := part1(parsedInput)
	allButLast := output[:len(output)-1]
	last := output[len(output)-1]

	// Ensure the last output is the solution
	assert.Equal(t, 4511442, last)

	// Ensure that all the preceding outputs are zero
	for _, outputElement := range allButLast {
		assert.Equal(t, 0, outputElement)
	}
}

func TestPart2(t *testing.T) {
	parsedInput, err := parseInput(rawInput)
	assert.Nil(t, err)

	output := part2(parsedInput)

	assert.Equal(t, 12648139, output)
}
