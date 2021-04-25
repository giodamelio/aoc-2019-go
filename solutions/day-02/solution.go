package main

import (
	_ "embed"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var discardLogger *log.Logger = log.New(ioutil.Discard, "", 0)

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

// TODO: add a panic handler
type intcodeComputer struct {
	memory         *intcodeMemory
	programCounter int
	opcodes        map[int]intcodeOpcode
	logger         *log.Logger
}

func newIntcodeComputer(initialMemory []int) *intcodeComputer {
	logger := discardLogger
	comp := new(intcodeComputer)
	comp.memory = newIntcodeMemory(initialMemory, logger)
	comp.programCounter = 0
	comp.logger = logger

	// Create the opcodes
	comp.opcodes = make(map[int]intcodeOpcode)
	comp.opcodes[1] = intcodeOpcode{
		name:      "ADD",
		opcode:    1,
		arguments: 3,
		execute: func(memory *intcodeMemory, arguments []int) {
			leftHandSide := memory.Get(arguments[0])
			rightHandSide := memory.Get(arguments[1])
			memory.Set(arguments[2], leftHandSide+rightHandSide)
		},
	}
	comp.opcodes[2] = intcodeOpcode{
		name:      "MULTIPLY",
		opcode:    2,
		arguments: 3,
		execute: func(memory *intcodeMemory, arguments []int) {
			leftHandSide := memory.Get(arguments[0])
			rightHandSide := memory.Get(arguments[1])
			memory.Set(arguments[2], leftHandSide*rightHandSide)
		},
	}
	comp.opcodes[99] = intcodeOpcode{
		name:      "HALT",
		opcode:    99,
		arguments: 0,
		execute: func(memory *intcodeMemory, arguments []int) {
		},
	}

	return comp
}

func (ic *intcodeComputer) Step() int {
	ic.logger.Println("Fetching OPCode at", ic.programCounter)
	opcode := ic.memory.Get(ic.programCounter)

	// Check if it is a valid opcode
	if _, ok := ic.opcodes[opcode]; !ok {
		// TODO handle this better
		panic("Invalid opcode")
	}

	opcodeArguments := ic.memory.GetRange(ic.programCounter+1, ic.opcodes[opcode].arguments)
	ic.logger.Println("Opcode", opcode)
	ic.logger.Println("Opcode args", opcodeArguments)
	ic.logger.Println("Memory before", ic.memory.rawMemory)
	ic.opcodes[opcode].execute(ic.memory, opcodeArguments)
	ic.logger.Println("Memory after", ic.memory.rawMemory)

	// Increment program counter
	ic.logger.Println("PC before", ic.programCounter)
	ic.programCounter = ic.programCounter + ic.opcodes[opcode].arguments + 1
	ic.logger.Println("PC after", ic.programCounter)

	return opcode
}

func (ic intcodeComputer) Run() {
	ic.logger.Println("Running!")
	for {
		opcode := ic.Step()
		// Special case for HALT
		if opcode == 99 {
			break
		}
	}
}

type intcodeMemory struct {
	rawMemory []int
	logger    *log.Logger
}

type intcodeOpcode struct {
	name      string
	opcode    int
	arguments int
	execute   func(*intcodeMemory, []int)
}

func newIntcodeMemory(initialMemory []int, logger *log.Logger) *intcodeMemory {
	mem := new(intcodeMemory)
	mem.rawMemory = initialMemory
	mem.logger = logger
	return mem
}

func (im intcodeMemory) Get(location int) int {
	value := im.rawMemory[location]
	im.logger.Printf("MEMORY[%d] get %d", location, value)
	return value
}

func (im intcodeMemory) GetRange(location int, length int) []int {
	value := im.rawMemory[location : location+length]
	im.logger.Printf("MEMORY[%d:%d] get %v", location, location+length, value)
	return value
}

func (im intcodeMemory) Set(location int, value int) {
	oldValue := im.rawMemory[location]
	im.logger.Printf("MEMORY[%d]: set to %d, old value %d", location, value, oldValue)
	im.rawMemory[location] = value
}

func part1(log *log.Logger, input []int) int {
	log.Println("Day 2 Part 1")

	computer := newIntcodeComputer(input)
	computer.logger = log

	computer.memory.Set(1, 12)
	computer.memory.Set(2, 2)
	computer.Run()

	return computer.memory.Get(0)
}

func part2(log *log.Logger, input []int) int {
	log.Println("Day 2 Part 2")

	max := 99
	for a := 0; a <= max; a++ {
		for b := 0; b <= max; b++ {
			// Make a new copy of the input
			inputCopy := make([]int, len(input))
			copy(inputCopy, input)

			computer := newIntcodeComputer(inputCopy)
			computer.memory.Set(1, a)
			computer.memory.Set(2, b)
			computer.Run()

			if computer.memory.Get(0) == 19690720 {
				return 100*a + b
			}
		}
	}

	panic("Should never happen")
}
