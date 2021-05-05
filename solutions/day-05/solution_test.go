package main

import (
	"os"
	"testing"

	"github.com/giodamelio/aoc-2020-go/intcode"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	out := zerolog.NewConsoleWriter()
	out.Out = os.Stderr
	out.NoColor = true
	log.Logger = log.Output(out)

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func TestPart1(t *testing.T) {
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	output := part1(parsedInput)
	allButLast := output[:len(output)-1]
	last := output[len(output)-1]

	// Ensure the last output is the solution
	assert.Equal(t, intcode.AddressValue(4511442), last)

	// Ensure that all the preceding outputs are zero
	for _, outputElement := range allButLast {
		assert.Equal(t, intcode.AddressValue(0), outputElement)
	}
}

func TestPart2(t *testing.T) {
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	output := part2(parsedInput)

	assert.Equal(t, intcode.AddressValue(12648139), output)
}
