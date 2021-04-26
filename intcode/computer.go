package intcode

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Computer struct {
	Memory         *Memory
	programCounter int
	opcodes        map[int]Opcode
	errorHandler   func(error)
	input          chan int
	output         chan int
}

func NewComputer(initialMemory []int) *Computer {
	log.Debug().Msg("[COMPUTER] Computer created")

	comp := new(Computer)
	comp.Memory = newMemory(initialMemory)
	comp.programCounter = 0
	comp.opcodes = Opcodes
	comp.input = make(chan int)
	comp.output = make(chan int)

	// Default to panicing when things go wrong
	comp.errorHandler = func(err error) {
		log.Err(err).Msg("[COMPUTER] Uncaught error")
		panic(err)
	}

	return comp
}

func (ic *Computer) Step() (int, error) {
	// Get the opcode at the address of the program counter
	opcode := ic.Memory.Get(ic.programCounter)
	log.Trace().Int("opcode", opcode).Msg("[COMPUTER] Retrieved opcode")

	// Check if it is a valid opcode
	if _, ok := ic.opcodes[opcode]; !ok {
		err := fmt.Errorf("invalid opcode: %d", opcode)
		return -1, err
	}

	// Get the arguments for the opcode
	opcodeArguments := ic.Memory.GetRange(ic.programCounter+1, ic.opcodes[opcode].arguments)
	log.Trace().Ints("arguments", opcodeArguments).Msg("[COMPUTER] Retrieved opcode arguments")

	// Execute the opcode
	operation := ic.opcodes[opcode]
	log.Trace().Str("opcodeName", operation.name).Msg("[COMPUTER] Executing operation")
	operation.execute(ic.Memory, opcodeArguments, ic.input, ic.output)

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
			close(ic.input)
			close(ic.output)
			break
		}
	}
}

func (ic *Computer) SendInput(input int) {
	go func() {
		ic.input <- input
	}()
}
