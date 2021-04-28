package main

import (
	_ "embed"
	"log"

	"github.com/giodamelio/aoc-2020-go/intcode"
)

// Read the raw input
//go:embed input.txt
var rawInput string

func part1(input []int) []int {
	log.Println("Day 5 Part 1")

	computer := intcode.NewComputer(input)

	// Select Air Conditioning Unit
	computer.SendInput(1)

	// Listen for outputs and when they are done send them on a channel
	outputChan := computer.GetOutputChannel()
	allOutputs := make(chan []int)

	listenForOutputs := func() {
		var outputs []int
		for i := range outputChan {
			outputs = append(outputs, i)
		}

		allOutputs <- outputs
	}
	go listenForOutputs()

	computer.Run()

	return <-allOutputs
}

func part2(input []int) int {
	log.Println("Day 5 Part 2")

	computer := intcode.NewComputer(input)

	// Select Air Conditioning Unit
	computer.SendInput(5)

	// Listen for outputs and when they are done send them on a channel
	outputChan := computer.GetOutputChannel()
	output := make(chan int)

	forwardOutputs := func() {
		output <- <-outputChan
	}
	go forwardOutputs()

	computer.Run()

	return <-output
}
