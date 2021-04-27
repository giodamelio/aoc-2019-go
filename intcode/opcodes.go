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
	execute    func(*Computer, []int) bool
}

var Opcodes map[int]Opcode = map[int]Opcode{
	1: {
		name:       "ADD",
		opcode:     1,
		parameters: []readWrite{Read, Read, Write},
		execute: func(computer *Computer, parameters []int) bool {
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

			return true
		},
	},
	2: {
		name:       "MULTIPLY",
		opcode:     2,
		parameters: []readWrite{Read, Read, Write},
		execute: func(computer *Computer, parameters []int) bool {
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

			return true
		},
	},
	3: {
		name:       "INPUT",
		opcode:     3,
		parameters: []readWrite{Write},
		execute: func(computer *Computer, parameters []int) bool {
			value := <-computer.input

			computer.Memory.Set(parameters[0], value)

			log.
				Debug().
				Int("input", value).
				Msg("[OPCODE] INPUT")

			return true
		},
	},
	4: {
		name:       "OUTPUT",
		opcode:     4,
		parameters: []readWrite{Read},
		execute: func(computer *Computer, parameters []int) bool {
			value := parameters[0]

			computer.output <- value

			log.
				Debug().
				Int("output", value).
				Msg("[OPCODE] OUTPUT")

			return true
		},
	},
	99: {
		name:       "HALT",
		opcode:     99,
		parameters: []readWrite{},
		execute: func(computer *Computer, parameters []int) bool {
			log.
				Debug().
				Msg("[OPCODE] HALT")

			return false
		},
	},
}
