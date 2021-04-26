package intcode

import (
	"fmt"
	"strconv"

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

func (ic Computer) parseOpcode(rawOpcode int) (int, []mode, error) {
	// Get opcode from 1s and 10s columns
	rawOpcodeString := fmt.Sprintf("%02d", rawOpcode)
	opcode, err := strconv.Atoi(rawOpcodeString[len(rawOpcodeString)-2:])
	if err != nil {
		return 0, nil, err
	}

	// Check if it is a valid opcode
	opcodeDefinition, ok := ic.opcodes[opcode]
	if !ok {
		err := fmt.Errorf("invalid opcode: %d", opcode)
		return 0, nil, err
	}

	// Get the number of parameters of that opcode
	parameterCount := opcodeDefinition.parameters

	// Format the rawOpcode to be as long as the opcode + parameters
	rawOpcodeStringWithParameterModes := fmt.Sprintf("%0*d", 2+parameterCount, rawOpcode)

	// Get all the mode settings without the opcode itself
	parameterModeSettings := rawOpcodeStringWithParameterModes[:len(rawOpcodeStringWithParameterModes)-2]

	// Convert the mode numbers to their mode enum
	// 0 = Position
	// 1 = Immediate
	outputModes := make([]mode, len(parameterModeSettings))
	for i, m := range parameterModeSettings {
		intM, err := strconv.Atoi(string(m))
		if err != nil {
			return 0, nil, err
		}

		outputModes[i] = mode(intM)
	}

	return opcode, outputModes, nil
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

	// Get the parameters for the opcode
	opcodeParameters := ic.Memory.GetRange(ic.programCounter+1, ic.opcodes[opcode].parameters)
	log.Trace().Ints("parameters", opcodeParameters).Msg("[COMPUTER] Retrieved opcode parameters")

	// Execute the opcode
	operation := ic.opcodes[opcode]
	log.Trace().Str("opcodeName", operation.name).Msg("[COMPUTER] Executing operation")
	operation.execute(ic.Memory, opcodeParameters, ic.input, ic.output)

	// Increment program counter
	ic.programCounter = ic.programCounter + ic.opcodes[opcode].parameters + 1

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

func (ic *Computer) GetOutputChannel() chan int {
	return ic.output
}
