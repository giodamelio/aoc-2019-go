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

func pipe(from chan intcode.AddressValue, to []chan intcode.AddressValue) {
	for i := range from {
		for _, t := range to {
			t <- i
		}
	}
}

func send(to chan intcode.AddressValue, value intcode.AddressValue) {
	to <- value
}

func amplifierChain(program []intcode.AddressValue, phaseSequence []int) int {
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
	send(computerA.Input, intcode.AddressValue(phaseSequence[0]))
	send(computerB.Input, intcode.AddressValue(phaseSequence[1]))
	send(computerC.Input, intcode.AddressValue(phaseSequence[2]))
	send(computerD.Input, intcode.AddressValue(phaseSequence[3]))
	send(computerE.Input, intcode.AddressValue(phaseSequence[4]))

	// Chain the computers inputs and outputs togather
	go pipe(computerA.Output, []chan intcode.AddressValue{computerB.Input})
	go pipe(computerB.Output, []chan intcode.AddressValue{computerC.Input})
	go pipe(computerC.Output, []chan intcode.AddressValue{computerD.Input})
	go pipe(computerD.Output, []chan intcode.AddressValue{computerE.Input})

	// Pass data to the start of the chain
	go send(computerA.Input, 0)

	return int(<-computerE.Output)
}

func amplifierChainFeedbackMode(program []intcode.AddressValue, phaseSequence []int) int {
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
	send(computerA.Input, intcode.AddressValue(phaseSequence[0]))
	send(computerB.Input, intcode.AddressValue(phaseSequence[1]))
	send(computerC.Input, intcode.AddressValue(phaseSequence[2]))
	send(computerD.Input, intcode.AddressValue(phaseSequence[3]))
	send(computerE.Input, intcode.AddressValue(phaseSequence[4]))

	// Chain the computers inputs and outputs togather
	go pipe(computerA.Output, []chan intcode.AddressValue{computerB.Input})
	go pipe(computerB.Output, []chan intcode.AddressValue{computerC.Input})
	go pipe(computerC.Output, []chan intcode.AddressValue{computerD.Input})
	go pipe(computerD.Output, []chan intcode.AddressValue{computerE.Input})

	// Hook E back up to A, but record all of it's outputs to a second channel
	computerEOutputs := make(chan intcode.AddressValue)

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
	lastComputerEOutput := make(chan intcode.AddressValue)

	returnLastItem := func() {
		var last intcode.AddressValue

		for out := range computerEOutputs {
			log.Info().Int64("value", int64(out)).Msg("Computer E output")
			last = out
		}

		lastComputerEOutput <- last
	}
	go returnLastItem()

	return int(<-lastComputerEOutput)
}

func part1(input []intcode.AddressValue) int {
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

func part2(input []intcode.AddressValue) int {
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
