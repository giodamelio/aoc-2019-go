package intcode

import "github.com/rs/zerolog/log"

type Opcode struct {
	name      string
	opcode    int
	arguments int
	execute   func(*Memory, []int)
}

var Opcodes map[int]Opcode = map[int]Opcode{
	1: {
		name:      "ADD",
		opcode:    1,
		arguments: 3,
		execute: func(memory *Memory, arguments []int) {
			leftHandSide := memory.Get(arguments[0])
			rightHandSide := memory.Get(arguments[1])
			output := leftHandSide + rightHandSide
			memory.Set(arguments[2], output)

			log.
				Debug().
				Int("leftHandSide", leftHandSide).
				Int("rightHandSide", rightHandSide).
				Int("output", output).
				Msg("[OPCODE] ADD")
		},
	},
	2: {
		name:      "MULTIPLY",
		opcode:    2,
		arguments: 3,
		execute: func(memory *Memory, arguments []int) {
			leftHandSide := memory.Get(arguments[0])
			rightHandSide := memory.Get(arguments[1])
			output := leftHandSide * rightHandSide
			memory.Set(arguments[2], output)

			log.
				Debug().
				Int("leftHandSide", leftHandSide).
				Int("rightHandSide", rightHandSide).
				Int("output", output).
				Msg("[OPCODE] MULTIPLY")
		},
	},
	99: {
		name:      "HALT",
		opcode:    99,
		arguments: 0,
		execute: func(memory *Memory, arguments []int) {
			log.
				Debug().
				Msg("[OPCODE] HALT")
		},
	},
}
