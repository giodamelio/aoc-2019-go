package intcode

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Computer struct {
	Memory             *Memory
	instructionPointer int
	opcodes            map[int]Opcode
	errorHandler       func(error)
	Input              chan int
	Output             chan int
}

func NewComputer(initialMemory []int) *Computer {
	log.Debug().Msg("[COMPUTER] Computer created")

	copyOfInitialMemory := copyMemory(initialMemory)

	comp := new(Computer)
	comp.Memory = newMemory(copyOfInitialMemory)
	comp.instructionPointer = 0
	comp.opcodes = Opcodes
	comp.Input = make(chan int)
	comp.Output = make(chan int)

	// Default to panicing when things go wrong
	comp.errorHandler = func(err error) {
		log.Err(err).Msg("[COMPUTER] Uncaught error")
		panic(err)
	}

	return comp
}

func reverseOutputModes(input []mode) []mode {
	inputLength := len(input)
	output := make([]mode, inputLength)

	for i, item := range input {
		output[inputLength-i-1] = item
	}

	return output
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
	parameterCount := len(opcodeDefinition.parameters)

	// Format the rawOpcode to be as long as the opcode + parameters
	baseOpcodeLength := 2
	rawOpcodeStringWithParameterModes := fmt.Sprintf("%0*d", baseOpcodeLength+parameterCount, rawOpcode)

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

	// Reverse the output modes since they are read right to left
	reversedOutputModes := reverseOutputModes(outputModes)

	return opcode, reversedOutputModes, nil
}

func (ic Computer) resolveParameters(
	memory *Memory, opcode int,
	opcodeParameters []int,
	parameterModes []mode,
) ([]int, error) {
	resolvedParameters := make([]int, len(opcodeParameters))

	for index, opcodeParameter := range opcodeParameters {
		parameterMode := ic.opcodes[opcode].parameters[index]

		switch parameterModes[index] {
		case Position:
			switch parameterMode {
			case Write:
				resolvedParameters[index] = opcodeParameter
			case Read:
				resolvedParameters[index] = memory.Get(opcodeParameter)
			default:
				err := fmt.Errorf("invalid parameter mode: %d", parameterModes[index])

				return nil, err
			}
		case Immediate:
			if parameterMode == Write {
				err := fmt.Errorf("write parameter cannot be in immediate mode: %d", opcodeParameter)

				return nil, err
			}

			resolvedParameters[index] = opcodeParameter
		default:
			err := fmt.Errorf("invalid mode: %d", parameterModes[index])

			return nil, err
		}
	}

	return resolvedParameters, nil
}

func (ic *Computer) SetInstructionPointer(address int) {
	ic.instructionPointer = address
}

func (ic *Computer) Step() (int, error) {
	// Get the opcode at the address of the instruction pointer
	opcode := ic.Memory.Get(ic.instructionPointer)
	log.Trace().Int("opcode", opcode).Msg("[COMPUTER] Retrieved opcode")

	// Parse the opcode
	opcode, parameterModes, err := ic.parseOpcode(opcode)
	if err != nil {
		return -1, err
	}

	log.
		Trace().
		Int("opcode", opcode).
		Msg("[COMPUTER] Parsed opcode")

	// Get the parameters for the opcode
	parametersLength := len(ic.opcodes[opcode].parameters)
	opcodeParameters := ic.Memory.GetRange(ic.instructionPointer+1, parametersLength)
	log.Trace().Ints("parameters", opcodeParameters).Msg("[COMPUTER] Retrieved opcode parameters")

	// Resolve the parameters based on the modes
	opcodeParameters, err = ic.resolveParameters(ic.Memory, opcode, opcodeParameters, parameterModes)
	if err != nil {
		return -1, err
	}

	log.Trace().Ints("parameters", opcodeParameters).Msg("[COMPUTER] Resolved opcode parameters")

	// Execute the opcode
	operation := ic.opcodes[opcode]
	log.Trace().Str("opcodeName", operation.name).Msg("[COMPUTER] Executing operation")
	operation.execute(ic, operation, opcodeParameters)

	return opcode, nil
}

func (ic Computer) Run() {
	for {
		opcode, err := ic.Step()
		if err != nil {
			ic.errorHandler(err)
		}

		// Special case for HALT
		if opcode == HALT {
			close(ic.Input)
			close(ic.Output)

			break
		}
	}
}
