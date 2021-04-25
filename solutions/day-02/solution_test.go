package main

import (
	"log"
	"os"
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

func TestNewIncodeComputer(t *testing.T) {
	computer := newIntcodeComputer([]int{1, 2, 3})

	assert.Equal(t, 0, computer.programCounter)
	assert.Equal(t, []int{1, 2, 3}, computer.memory.rawMemory)
}

func TestIntcodeComputerStep(t *testing.T) {
	computer := newIntcodeComputer([]int{1, 0, 0, 0})
	computer.logger = log.New(os.Stdout, "", 0)

	computer.Step()
	assert.Equal(t, []int{2, 0, 0, 0}, computer.memory.rawMemory)
}

func TestIntcodeComputerRun(t *testing.T) {
	computer := newIntcodeComputer([]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50})
	computer.logger = log.New(os.Stdout, "", 0)

	computer.Run()
	assert.Equal(t, 3500, computer.memory.rawMemory[0])
}

func TestIntcodeComputerMemoryGet(t *testing.T) {
	computer := newIntcodeComputer([]int{1, 2, 3})

	assert.Equal(t, 2, computer.memory.Get(1))
}

func TestIntcodeComputerMemoryGetRange(t *testing.T) {
	computer := newIntcodeComputer([]int{1, 2, 3, 4, 5, 6})

	assert.Equal(t, []int{3, 4, 5}, computer.memory.GetRange(2, 3))
}

func TestIntcodeComputerMemorySet(t *testing.T) {
	computer := newIntcodeComputer([]int{1, 2, 3})

	computer.memory.Set(1, 10)
	assert.Equal(t, 10, computer.memory.Get(1))
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
