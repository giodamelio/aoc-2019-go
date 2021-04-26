package intcode

import (
	"fmt"
)

type Computer struct {
	Memory         *Memory
	programCounter int
	opcodes        map[int]Opcode
	errorHandler   func(error)
}

func NewComputer(initialMemory []int) *Computer {
	comp := new(Computer)
	comp.Memory = newMemory(initialMemory)
	comp.programCounter = 0
	comp.opcodes = Opcodes

	// Default to panicing when things go wrong
	comp.errorHandler = func(err error) {
		panic(err)
	}

	return comp
}

func (ic *Computer) Step() (int, error) {
	// Get the opcode at the address of the program counter
	opcode := ic.Memory.Get(ic.programCounter)

	// Check if it is a valid opcode
	if _, ok := ic.opcodes[opcode]; !ok {
		err := fmt.Errorf("invalid opcode: %d", opcode)
		return -1, err
	}

	// Get the arguments for the opcode
	opcodeArguments := ic.Memory.GetRange(ic.programCounter+1, ic.opcodes[opcode].arguments)

	// Execute the opcode
	ic.opcodes[opcode].execute(ic.Memory, opcodeArguments)

	// Increment program counter
	ic.programCounter = ic.programCounter + ic.opcodes[opcode].arguments + 1

	return opcode, nil
}

func (ic Computer) Run() {
	for {
		opcode, err := ic.Step()
		if err != nil {
			ic.errorHandler(err)
		}

		// Special case for HALT
		if opcode == 99 {
			break
		}
	}
}
