package intcode

import (
	"github.com/rs/zerolog/log"
)

type mode int

const (
	Position mode = iota
	Immediate
)

type Opcode struct {
	name       string
	opcode     int
	parameters int
	execute    func(*Memory, []int, chan int, chan int)
}

var Opcodes map[int]Opcode = map[int]Opcode{
	1: {
		name:       "ADD",
		opcode:     1,
		parameters: 3,
		execute: func(memory *Memory, parameters []int, input chan int, output chan int) {
			leftHandSide := memory.Get(parameters[0])
			rightHandSide := memory.Get(parameters[1])
			result := leftHandSide + rightHandSide
			memory.Set(parameters[2], result)

			log.
				Debug().
				Int("leftHandSide", leftHandSide).
				Int("rightHandSide", rightHandSide).
				Int("result", result).
				Msg("[OPCODE] ADD")
		},
	},
	2: {
		name:       "MULTIPLY",
		opcode:     2,
		parameters: 3,
		execute: func(memory *Memory, parameters []int, input chan int, output chan int) {
			leftHandSide := memory.Get(parameters[0])
			rightHandSide := memory.Get(parameters[1])
			result := leftHandSide * rightHandSide
			memory.Set(parameters[2], result)

			log.
				Debug().
				Int("leftHandSide", leftHandSide).
				Int("rightHandSide", rightHandSide).
				Int("result", result).
				Msg("[OPCODE] MULTIPLY")
		},
	},
	3: {
		name:       "INPUT",
		opcode:     3,
		parameters: 1,
		execute: func(memory *Memory, parameters []int, input chan int, output chan int) {
			value := <-input

			memory.Set(parameters[0], value)

			log.
				Debug().
				Int("input", value).
				Msg("[OPCODE] INPUT")
		},
	},
	4: {
		name:       "OUTPUT",
		opcode:     4,
		parameters: 1,
		execute: func(memory *Memory, parameters []int, input chan int, output chan int) {
			value := memory.Get(parameters[0])

			output <- value

			log.
				Debug().
				Int("output", value).
				Msg("[OPCODE] OUTPUT")
		},
	},
	99: {
		name:       "HALT",
		opcode:     99,
		parameters: 0,
		execute: func(memory *Memory, parameters []int, input chan int, output chan int) {
			log.
				Debug().
				Msg("[OPCODE] HALT")
		},
	},
}
