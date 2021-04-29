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

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func TestChainingComputers(t *testing.T) {
	// Take an input, double it and output the result
	program := []int{
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

	go send(computer1.Input, 10)
	go pipe(computer1.Output, computer2.Input)

	go computer1.Run()
	go computer2.Run()

	output := <-computer2.Output
	assert.Equal(t, 40, output)
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
	input := make(chan int)
	output := make(chan int)

	go send(input, 10)
	go pipe(input, output)

	assert.Equal(t, 10, <-output)
}

func TestSend(t *testing.T) {
	output := make(chan int)

	go send(output, 10)

	assert.Equal(t, 10, <-output)
}

func TestAmplifyChain(t *testing.T) {
	exampleProgram1 := []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0}
	output1 := amplifierChain(exampleProgram1, []int{4, 3, 2, 1, 0})
	assert.Equal(t, 43210, output1)

	exampleProgram2 := []int{
		3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23,
		101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0,
	}
	output2 := amplifierChain(exampleProgram2, []int{0, 1, 2, 3, 4})
	assert.Equal(t, 54321, output2)

	exampleProgram3 := []int{
		3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33,
		1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0,
	}
	output3 := amplifierChain(exampleProgram3, []int{1, 0, 4, 3, 2})
	assert.Equal(t, 65210, output3)
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

	assert.Equal(t, 0, output)
}
