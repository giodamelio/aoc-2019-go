package intcode

import (
	"github.com/rs/zerolog/log"
)

type mode int

const (
	Position mode = iota
	Immediate
)

type readWrite int

const (
	Read readWrite = iota
	Write
)

type Opcode struct {
	name       string
	opcode     int
	parameters []readWrite
	execute    func(*Computer, Opcode, []int)
}

// The total length of the opcode including parameters
func (o Opcode) length() int {
	return 1 + len(o.parameters)
}

func (o Opcode) incrementInstructionPointer(computer *Computer) {
	computer.SetInstructionPointer(computer.instructionPointer + o.length())
}

var Opcodes map[int]Opcode = map[int]Opcode{
	1: {
		name:       "ADD",
		opcode:     1,
		parameters: []readWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			leftHandSide := parameters[0]
			rightHandSide := parameters[1]
			result := leftHandSide + rightHandSide
			computer.Memory.Set(parameters[2], result)

			log.
				Debug().
				Int("leftHandSide", leftHandSide).
				Int("rightHandSide", rightHandSide).
				Int("result", result).
				Msg("[OPCODE] ADD")

			operation.incrementInstructionPointer(computer)
		},
	},
	2: {
		name:       "MULTIPLY",
		opcode:     2,
		parameters: []readWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			leftHandSide := parameters[0]
			rightHandSide := parameters[1]
			result := leftHandSide * rightHandSide
			computer.Memory.Set(parameters[2], result)

			log.
				Debug().
				Int("leftHandSide", leftHandSide).
				Int("rightHandSide", rightHandSide).
				Int("result", result).
				Msg("[OPCODE] MULTIPLY")

			operation.incrementInstructionPointer(computer)
		},
	},
	3: {
		name:       "INPUT",
		opcode:     3,
		parameters: []readWrite{Write},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			value := <-computer.input

			computer.Memory.Set(parameters[0], value)

			log.
				Debug().
				Int("input", value).
				Msg("[OPCODE] INPUT")

			operation.incrementInstructionPointer(computer)
		},
	},
	4: {
		name:       "OUTPUT",
		opcode:     4,
		parameters: []readWrite{Read},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			value := parameters[0]

			computer.output <- value

			log.
				Debug().
				Int("output", value).
				Msg("[OPCODE] OUTPUT")

			operation.incrementInstructionPointer(computer)
		},
	},
	5: {
		name:       "JUMP-IF-TRUE",
		opcode:     5,
		parameters: []readWrite{Read, Read},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			condition := parameters[0]
			address := parameters[1]

			log.
				Debug().
				Int("condition", condition).
				Int("address", address).
				Msg("[OPCODE] JUMP-IF-TRUE")

			if condition != 0 {
				computer.SetInstructionPointer(address)
			}
		},
	},
	99: {
		name:       "HALT",
		opcode:     99,
		parameters: []readWrite{},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			log.
				Debug().
				Msg("[OPCODE] HALT")
		},
	},
}
