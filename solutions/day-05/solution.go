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
	sendInput := func() {
		computer.Input <- 1
	}
	go sendInput()

	// Listen for outputs and when they are done send them on a channel
	outputChan := computer.Output
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
	sendInput := func() {
		computer.Input <- 5
	}
	go sendInput()

	// Listen for outputs and when they are done send them on a channel
	outputChan := computer.Output
	output := make(chan int)

	forwardOutputs := func() {
		output <- <-outputChan
	}
	go forwardOutputs()

	computer.Run()

	return <-output
}
