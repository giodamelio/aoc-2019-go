package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// Read the raw input
//go:embed input.txt
var rawInput string

func main() {
	log := log.New(os.Stdout, "", 0)

	// Parse input flags
	partPtr := flag.String("part", "both", "The part you want to run \"1\", \"2\" or \"both\"")
	flag.Parse()

	// Ensure a valid part
	if *partPtr != "1" && *partPtr != "2" && *partPtr != "both" {
		log.Fatal("Part must be 1, 2, or both")
	}

	// Parse the input
	parsedInput, err := parseInput(rawInput)
	if err != nil {
		log.Fatal("Failed to parse input")
	}

	if *partPtr == "1" {
		log.Printf("Part 1 solution: %d", part1(log, parsedInput))
	}
	if *partPtr == "2" {
		log.Printf("Part 2 solution: %d", part2(log, parsedInput))
	}
	if *partPtr == "both" {
		log.Printf("Part 1 solution: %d", part1(log, parsedInput))
		log.Printf("Part 2 solution: %d", part2(log, parsedInput))
	}
}

func parseInput(input string) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	numbers := make([]int, len(lines))

	for index, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		numbers[index] = number
	}

	return numbers, nil
}

func calculateMass(mass int) int {
	// `/` does floor division with integers
	return mass/3 - 2
}

func calculateMassWithFuel(mass int) int {
	var totalMass int = 0
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
