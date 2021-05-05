package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
)

// Read the raw input
//go:embed input.txt
var rawInput string

func parseInput(input string) ([]int, error) {
	bytes := strings.Split(strings.TrimSpace(input), "\n")
	numbers := make([]int, len(bytes))

	for index, programByte := range bytes {
		number, err := strconv.Atoi(programByte)
		if err != nil {
			return nil, err
		}

		numbers[index] = (number)
	}

	return numbers, nil
}

func calculateMass(mass int) int {
	// `/` does floor division with integers
	return mass/3 - 2
}

func calculateMassWithFuel(mass int) int {
	totalMass := 0

	nextMass := mass

	// Repeatedly calculate mass of fuel until it reaches zero or less
	for {
		newMass := calculateMass(nextMass)
		if newMass <= 0 {
			break
		}

		nextMass = newMass
		totalMass += newMass
	}

	return totalMass
}

func part1(log *log.Logger, input []int) int {
	log.Println("Day 1 Part 1")

	sum := int(0)
	for _, moduleMass := range input {
		sum += calculateMass(moduleMass)
	}

	return sum
}

func part2(log *log.Logger, input []int) int {
	log.Println("Day 1 Part 2")

	sum := int(0)
	for _, moduleMass := range input {
		sum += calculateMassWithFuel(moduleMass)
	}

	return sum
}
