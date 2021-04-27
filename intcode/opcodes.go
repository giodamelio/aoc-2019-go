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

const (
	ADD         = 1
	MULTIPLY    = 2
	INPUT       = 3
	OUTPUT      = 4
	JUMPIFTRUE  = 5
	JUMPIFFALSE = 6
	LESSTHAN    = 7
	EQUALS      = 8
	HALT        = 99
)

type Opcode struct {
	name       string
	opcode     int
	parameters []readWrite
	execute    func(*Computer, Opcode, []int)
}

// The total length of the opcode including parameters.
func (o Opcode) length() int {
	return 1 + len(o.parameters)
}

func (o Opcode) incrementInstructionPointer(computer *Computer) {
	computer.SetInstructionPointer(computer.instructionPointer + o.length())
}

var Opcodes = map[int]Opcode{
	ADD: {
		name:       "ADD",
		opcode:     ADD,
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
	MULTIPLY: {
		name:       "MULTIPLY",
		opcode:     MULTIPLY,
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
	INPUT: {
		name:       "INPUT",
		opcode:     INPUT,
		parameters: []readWrite{Write},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			address := parameters[0]
			value := <-computer.input

			computer.Memory.Set(address, value)

			log.
				Debug().
				Int("input", value).
				Int("address", address).
				Msg("[OPCODE] INPUT")

			operation.incrementInstructionPointer(computer)
		},
	},
	OUTPUT: {
		name:       "OUTPUT",
		opcode:     OUTPUT,
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
	JUMPIFTRUE: {
		name:       "JUMP-IF-TRUE",
		opcode:     JUMPIFTRUE,
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
			} else {
				operation.incrementInstructionPointer(computer)
			}
		},
	},
	JUMPIFFALSE: {
		name:       "JUMP-IF-FALSE",
		opcode:     JUMPIFFALSE,
		parameters: []readWrite{Read, Read},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			condition := parameters[0]
			address := parameters[1]

			log.
				Debug().
				Int("condition", condition).
				Int("address", address).
				Msg("[OPCODE] JUMP-IF-FALSE")

			if condition == 0 {
				computer.SetInstructionPointer(address)
			} else {
				operation.incrementInstructionPointer(computer)
			}
		},
	},
	LESSTHAN: {
		name:       "LESS-THAN",
		opcode:     LESSTHAN,
		parameters: []readWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			lhs := parameters[0]
			rhs := parameters[1]
			outputAddress := parameters[2]

			var output int
			if lhs < rhs {
				output = 1
			} else {
				output = 0
			}

			log.
				Debug().
				Int("lhs", lhs).
				Int("rhs", rhs).
				Int("outputAddress", outputAddress).
				Int("output", output).
				Msg("[OPCODE] LESS-THAN")

			computer.Memory.Set(outputAddress, output)

			operation.incrementInstructionPointer(computer)
		},
	},
	EQUALS: {
		name:       "EQUALS",
		opcode:     EQUALS,
		parameters: []readWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			lhs := parameters[0]
			rhs := parameters[1]
			outputAddress := parameters[2]

			var output int
			if lhs == rhs {
				output = 1
			} else {
				output = 0
			}

			log.
				Debug().
				Int("lhs", lhs).
				Int("rhs", rhs).
				Int("outputAddress", outputAddress).
				Int("output", output).
				Msg("[OPCODE] EQUALS")

			computer.Memory.Set(outputAddress, output)

			operation.incrementInstructionPointer(computer)
		},
	},
	HALT: {
		name:       "HALT",
		opcode:     HALT,
		parameters: []readWrite{},
		execute: func(computer *Computer, operation Opcode, parameters []int) {
			log.
				Debug().
				Msg("[OPCODE] HALT")
		},
	},
}
