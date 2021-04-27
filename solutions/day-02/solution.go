package main

import (
	_ "embed"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/giodamelio/aoc-2020-go/intcode"
)

var discardLogger = log.New(ioutil.Discard, "", 0)

// Read the raw input
//go:embed input.txt
var rawInput string

func parseInput(input string) ([]int, error) {
	bytes := strings.Split(strings.TrimSpace(input), ",")
	numbers := make([]int, len(bytes))

	for index, programByte := range bytes {
		number, err := strconv.Atoi(programByte)
		if err != nil {
			return nil, err
		}

		numbers[index] = number
	}

	return numbers, nil
}

func part1(log *log.Logger, input []int) int {
	log.Println("Day 2 Part 1")

	computer := intcode.NewComputer(input)

	computer.Memory.Set(1, 12)
	computer.Memory.Set(2, 2)
	computer.Run()

	return computer.Memory.Get(0)
}

func part2(log *log.Logger, input []int) int {
	log.Println("Day 2 Part 2")

	max := 99
	for a := 0; a <= max; a++ {
		for b := 0; b <= max; b++ {
			// Make a new copy of the input
			inputCopy := make([]int, len(input))
			copy(inputCopy, input)

			computer := intcode.NewComputer(inputCopy)
			computer.Memory.Set(1, a)
			computer.Memory.Set(2, b)
			computer.Run()

			if computer.Memory.Get(0) == 19690720 {
				return 100*a + b
			}
		}
	}

	panic("Should never happen")
}
