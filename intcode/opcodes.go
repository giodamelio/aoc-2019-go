package intcode

import (
	"github.com/rs/zerolog/log"
)

type Mode int

const (
	Position Mode = iota
	Immediate
)

type ReadWrite int

const (
	Read ReadWrite = iota
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
	Name       string
	Opcode     AddressValue
	Parameters []ReadWrite
	execute    func(*Computer, Opcode, []AddressValue)
}

// The total length of the opcode including parameters.
func (o Opcode) length() int {
	return 1 + len(o.Parameters)
}

func (o Opcode) incrementInstructionPointer(computer *Computer) {
	computer.SetInstructionPointer(computer.instructionPointer + AddressLocation(o.length()))
}

var Opcodes = map[AddressValue]Opcode{
	ADD: {
		Name:       "ADD",
		Opcode:     ADD,
		Parameters: []ReadWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			leftHandSide := parameters[0]
			rightHandSide := parameters[1]
			result := leftHandSide + rightHandSide
			computer.Memory.Set(parameters[2], result)

			log.
				Debug().
				Int64("leftHandSide", int64(leftHandSide)).
				Int64("rightHandSide", int64(rightHandSide)).
				Int64("result", int64(result)).
				Msg("[OPCODE] ADD")

			operation.incrementInstructionPointer(computer)
		},
	},
	MULTIPLY: {
		Name:       "MULTIPLY",
		Opcode:     MULTIPLY,
		Parameters: []ReadWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			leftHandSide := parameters[0]
			rightHandSide := parameters[1]
			result := leftHandSide * rightHandSide
			computer.Memory.Set(parameters[2], result)

			log.
				Debug().
				Int64("leftHandSide", int64(leftHandSide)).
				Int64("rightHandSide", int64(rightHandSide)).
				Int64("result", int64(result)).
				Msg("[OPCODE] MULTIPLY")

			operation.incrementInstructionPointer(computer)
		},
	},
	INPUT: {
		Name:       "INPUT",
		Opcode:     INPUT,
		Parameters: []ReadWrite{Write},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			address := parameters[0]
			value := <-computer.Input

			computer.Memory.Set(address, value)

			log.
				Debug().
				Int64("input", int64(value)).
				Int64("address", int64(address)).
				Msg("[OPCODE] INPUT")

			operation.incrementInstructionPointer(computer)
		},
	},
	OUTPUT: {
		Name:       "OUTPUT",
		Opcode:     OUTPUT,
		Parameters: []ReadWrite{Read},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			value := parameters[0]

			computer.Output <- value

			log.
				Debug().
				Int64("output", int64(value)).
				Msg("[OPCODE] OUTPUT")

			operation.incrementInstructionPointer(computer)
		},
	},
	JUMPIFTRUE: {
		Name:       "JUMP-IF-TRUE",
		Opcode:     JUMPIFTRUE,
		Parameters: []ReadWrite{Read, Read},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			condition := parameters[0]
			address := parameters[1]

			log.
				Debug().
				Int64("condition", int64(condition)).
				Int64("address", int64(address)).
				Msg("[OPCODE] JUMP-IF-TRUE")

			if condition != 0 {
				computer.SetInstructionPointer(AddressLocation(address))
			} else {
				operation.incrementInstructionPointer(computer)
			}
		},
	},
	JUMPIFFALSE: {
		Name:       "JUMP-IF-FALSE",
		Opcode:     JUMPIFFALSE,
		Parameters: []ReadWrite{Read, Read},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			condition := parameters[0]
			address := parameters[1]

			log.
				Debug().
				Int64("condition", int64(condition)).
				Int64("address", int64(address)).
				Msg("[OPCODE] JUMP-IF-FALSE")

			if condition == 0 {
				computer.SetInstructionPointer(AddressLocation(address))
			} else {
				operation.incrementInstructionPointer(computer)
			}
		},
	},
	LESSTHAN: {
		Name:       "LESS-THAN",
		Opcode:     LESSTHAN,
		Parameters: []ReadWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			lhs := parameters[0]
			rhs := parameters[1]
			outputAddress := parameters[2]

			var output AddressValue
			if lhs < rhs {
				output = 1
			} else {
				output = 0
			}

			log.
				Debug().
				Int64("lhs", int64(lhs)).
				Int64("rhs", int64(rhs)).
				Int64("outputAddress", int64(outputAddress)).
				Int64("output", int64(output)).
				Msg("[OPCODE] LESS-THAN")

			computer.Memory.Set(outputAddress, output)

			operation.incrementInstructionPointer(computer)
		},
	},
	EQUALS: {
		Name:       "EQUALS",
		Opcode:     EQUALS,
		Parameters: []ReadWrite{Read, Read, Write},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			lhs := parameters[0]
			rhs := parameters[1]
			outputAddress := parameters[2]

			var output AddressValue
			if lhs == rhs {
				output = 1
			} else {
				output = 0
			}

			log.
				Debug().
				Int64("lhs", int64(lhs)).
				Int64("rhs", int64(rhs)).
				Int64("outputAddress", int64(outputAddress)).
				Int64("output", int64(output)).
				Msg("[OPCODE] EQUALS")

			computer.Memory.Set(outputAddress, output)

			operation.incrementInstructionPointer(computer)
		},
	},
	HALT: {
		Name:       "HALT",
		Opcode:     HALT,
		Parameters: []ReadWrite{},
		execute: func(computer *Computer, operation Opcode, parameters []AddressValue) {
			log.
				Debug().
				Msg("[OPCODE] HALT")
		},
	},
}
