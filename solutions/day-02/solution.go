package main

import (
	_ "embed"
	"io/ioutil"
	"log"

	"github.com/giodamelio/aoc-2020-go/intcode"
)

var discardLogger = log.New(ioutil.Discard, "", 0)

// Read the raw input
//go:embed input.txt
var rawInput string

func part1(log *log.Logger, input []intcode.AddressValue) int64 {
	log.Println("Day 2 Part 1")

	computer := intcode.NewComputer(input)

	computer.Memory.Set(1, 12)
	computer.Memory.Set(2, 2)
	computer.Run()

	return int64(computer.Memory.Get(0))
}

func part2(log *log.Logger, input []intcode.AddressValue) int64 {
	log.Println("Day 2 Part 2")

	max := 99
	for a := 0; a <= max; a++ {
		for b := 0; b <= max; b++ {
			// Make a new copy of the input
			inputCopy := make([]intcode.AddressValue, len(input))
			copy(inputCopy, input)

			computer := intcode.NewComputer(inputCopy)
			computer.Memory.Set(1, intcode.AddressValue(a))
			computer.Memory.Set(2, intcode.AddressValue(b))
			computer.Run()

			if computer.Memory.Get(0) == 19690720 {
				return int64(100*a + b)
			}
		}
	}

	panic("Should never happen")
}
