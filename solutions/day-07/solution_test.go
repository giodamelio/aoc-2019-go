package main

import (
	"os"
	"testing"

	"github.com/giodamelio/aoc-2020-go/intcode"
	"github.com/gitchander/permutation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	out := zerolog.NewConsoleWriter()
	out.Out = os.Stderr
	out.NoColor = true
	log.Logger = log.Output(out)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func TestChainingComputers(t *testing.T) {
	// Take an input, double it and output the result
	program := []intcode.AddressValue{
		// Program
		3, 9, //         INPUT			read input into address 9
		102, 2, 9, 9, // MULTIPLY 	address 9 by 2
		4, 9, //         OUTPUT			send the contents of address 9 to output
		99, //           HALT
		// Data
		0,
	}

	computer1 := intcode.NewComputer(program)
	computer2 := intcode.NewComputer(program)

	go send(computer1.Input, intcode.AddressValue(10))
	go pipe(computer1.Output, []chan intcode.AddressValue{computer2.Input})

	go computer1.Run()
	go computer2.Run()

	output := <-computer2.Output
	assert.Equal(t, intcode.AddressValue(40), output)
}

func TestPermutations(t *testing.T) {
	a := []int{1, 2, 3}
	p := permutation.New(permutation.IntSlice(a))

	i := 0
	for p.Next() {
		i++
	}
	assert.Equal(t, 6, i)
}

func TestPipe(t *testing.T) {
	input := make(chan intcode.AddressValue)
	output := make(chan intcode.AddressValue)
	output2 := make(chan intcode.AddressValue)

	go send(input, 10)
	go pipe(input, []chan intcode.AddressValue{output, output2})

	assert.Equal(t, intcode.AddressValue(10), <-output)
	assert.Equal(t, intcode.AddressValue(10), <-output2)
}

func TestSend(t *testing.T) {
	output := make(chan intcode.AddressValue)

	go send(output, 10)

	assert.Equal(t, intcode.AddressValue(10), <-output)
}

func TestAmplifyChain(t *testing.T) {
	exampleProgram1 := []intcode.AddressValue{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0}
	output1 := amplifierChain(exampleProgram1, []int{4, 3, 2, 1, 0})
	assert.Equal(t, 43210, output1)

	exampleProgram2 := []intcode.AddressValue{
		3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23,
		101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0,
	}
	output2 := amplifierChain(exampleProgram2, []int{0, 1, 2, 3, 4})
	assert.Equal(t, 54321, output2)

	exampleProgram3 := []intcode.AddressValue{
		3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33,
		1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0,
	}
	output3 := amplifierChain(exampleProgram3, []int{1, 0, 4, 3, 2})
	assert.Equal(t, 65210, output3)
}

func TestAmplifyChainFeedbackMode(t *testing.T) {
	exampleProgram1 := []intcode.AddressValue{
		3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26,
		27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5,
	}
	output1 := amplifierChainFeedbackMode(exampleProgram1, []int{9, 8, 7, 6, 5})
	assert.Equal(t, 139629729, output1)

	exampleProgram2 := []intcode.AddressValue{
		3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54,
		-5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4,
		53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10,
	}
	output2 := amplifierChainFeedbackMode(exampleProgram2, []int{9, 7, 8, 5, 6})
	assert.Equal(t, 18216, output2)
}

func TestPart1(t *testing.T) {
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	output := part1(parsedInput)

	assert.Equal(t, 359142, output)
}

func TestPart2(t *testing.T) {
	parsedInput, err := intcode.ParseInput(rawInput)
	assert.Nil(t, err)

	output := part2(parsedInput)

	assert.Equal(t, 4374895, output)
}
