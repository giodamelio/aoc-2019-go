package main

import (
	_ "embed"

	"github.com/giodamelio/aoc-2020-go/intcode"
	"github.com/gitchander/permutation"
	"github.com/rs/zerolog/log"
)

// Read the raw input
//go:embed input.txt
var rawInput string

func pipe(from chan int, to []chan int) {
	for i := range from {
		for _, t := range to {
			t <- i
		}
	}
}

func send(to chan int, value int) {
	to <- value
}

func amplifierChain(program []int, phaseSequence []int) int {
	computerA := intcode.NewComputer(program)
	computerB := intcode.NewComputer(program)
	computerC := intcode.NewComputer(program)
	computerD := intcode.NewComputer(program)
	computerE := intcode.NewComputer(program)

	// Start the computers
	go computerA.Run()
	go computerB.Run()
	go computerC.Run()
	go computerD.Run()
	go computerE.Run()

	// Pass the phase settings
	send(computerA.Input, phaseSequence[0])
	send(computerB.Input, phaseSequence[1])
	send(computerC.Input, phaseSequence[2])
	send(computerD.Input, phaseSequence[3])
	send(computerE.Input, phaseSequence[4])

	// Chain the computers inputs and outputs togather
	go pipe(computerA.Output, []chan int{computerB.Input})
	go pipe(computerB.Output, []chan int{computerC.Input})
	go pipe(computerC.Output, []chan int{computerD.Input})
	go pipe(computerD.Output, []chan int{computerE.Input})

	// Pass data to the start of the chain
	go send(computerA.Input, 0)

	return <-computerE.Output
}

func amplifierChainFeedbackMode(program []int, phaseSequence []int) int {
	computerA := intcode.NewComputer(program)
	computerB := intcode.NewComputer(program)
	computerC := intcode.NewComputer(program)
	computerD := intcode.NewComputer(program)
	computerE := intcode.NewComputer(program)

	// Start the computers
	go computerA.Run()
	go computerB.Run()
	go computerC.Run()
	go computerD.Run()
	go computerE.Run()

	// Pass the phase settings
	send(computerA.Input, phaseSequence[0])
	send(computerB.Input, phaseSequence[1])
	send(computerC.Input, phaseSequence[2])
	send(computerD.Input, phaseSequence[3])
	send(computerE.Input, phaseSequence[4])

	// Chain the computers inputs and outputs togather
	go pipe(computerA.Output, []chan int{computerB.Input})
	go pipe(computerB.Output, []chan int{computerC.Input})
	go pipe(computerC.Output, []chan int{computerD.Input})
	go pipe(computerD.Output, []chan int{computerE.Input})

	// Hook E back up to A, but record all of it's outputs to a second channel
	computerEOutputs := make(chan int)

	pipeUnlessHalted := func() {
		for i := range computerE.Output {
			computerEOutputs <- i

			if computerE.State != "halted" {
				computerA.Input <- i
			}
		}

		close(computerEOutputs)
	}
	go pipeUnlessHalted()

	// Pass data to the start of the chain
	go send(computerA.Input, 0)
	// Listen to all the outputs from computerE, output just the last one
	lastComputerEOutput := make(chan int)

	returnLastItem := func() {
		last := 0

		for out := range computerEOutputs {
			log.Info().Int("value", out).Msg("Computer E output")
			last = out
		}

		lastComputerEOutput <- last
	}
	go returnLastItem()

	return <-lastComputerEOutput
}

func part1(input []int) int {
	log.Info().Msg("Day 7 Part 1")

	phaseSettingPermutation := []int{0, 1, 2, 3, 4}
	permutations := permutation.New(permutation.IntSlice(phaseSettingPermutation))

	maxOutput := 0

	for permutations.Next() {
		output := amplifierChain(input, phaseSettingPermutation)
		if output > maxOutput {
			maxOutput = output
		}
	}

	return maxOutput
}

func part2(input []int) int {
	log.Info().Msg("Day 7 Part 2")

	phaseSettingPermutation := []int{5, 6, 7, 8, 9}
	permutations := permutation.New(permutation.IntSlice(phaseSettingPermutation))

	maxOutput := 0

	for permutations.Next() {
		output := amplifierChainFeedbackMode(input, phaseSettingPermutation)
		if output > maxOutput {
			maxOutput = output
		}
	}

	return maxOutput
}
