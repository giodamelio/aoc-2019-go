package intcode

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Computer struct {
	Memory             *Memory
	instructionPointer AddressLocation
	opcodes            map[AddressValue]Opcode
	errorHandler       func(error)
	Input              chan AddressValue
	Output             chan AddressValue
	State              string
	Name               string
}

func NewComputer(initialMemory []AddressValue) *Computer {
	log.Debug().Msg("[COMPUTER] Computer created")

	copyOfInitialMemory := copyMemory(initialMemory)

	comp := new(Computer)
	comp.Memory = newMemory(copyOfInitialMemory)
	comp.instructionPointer = 0
	comp.opcodes = Opcodes
	comp.Input = make(chan AddressValue)
	comp.Output = make(chan AddressValue)
	comp.State = "pre-run"
	comp.Name = "computer"

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

func (ic Computer) parseOpcode(rawOpcode AddressValue) (AddressValue, []mode, error) {
	// Get opcode from 1s and 10s columns
	rawOpcodeString := fmt.Sprintf("%02d", rawOpcode)

	opcode, err := strconv.Atoi(rawOpcodeString[len(rawOpcodeString)-2:])
	if err != nil {
		return 0, nil, err
	}

	// Check if it is a valid opcode
	opcodeDefinition, ok := ic.opcodes[AddressValue(opcode)]
	if !ok {
		err := fmt.Errorf("invalid opcode: %d", opcode)

		return 0, nil, err
	}

	// Get the number of parameters of that opcode
	parameterCount := len(opcodeDefinition.Parameters)

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

	return AddressValue(opcode), reversedOutputModes, nil
}

func (ic Computer) resolveParameters(
	memory *Memory, opcode AddressValue,
	opcodeParameters []AddressValue,
	parameterModes []mode,
) ([]AddressValue, error) {
	resolvedParameters := make([]AddressValue, len(opcodeParameters))

	for index, opcodeParameter := range opcodeParameters {
		parameterMode := ic.opcodes[opcode].Parameters[index]

		switch parameterModes[index] {
		case Position:
			switch parameterMode {
			case Write:
				resolvedParameters[index] = opcodeParameter
			case Read:
				resolvedParameters[index] = memory.Get(AddressLocation(opcodeParameter))
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

func (ic *Computer) SetInstructionPointer(address AddressLocation) {
	ic.instructionPointer = address
}

func (ic *Computer) Step() (AddressValue, error) {
	// Get the opcode at the address of the instruction pointer
	opcode := ic.Memory.Get(ic.instructionPointer)
	log.Trace().Int64("opcode", int64(opcode)).Msg("[COMPUTER] Retrieved opcode")

	// Parse the opcode
	opcode, parameterModes, err := ic.parseOpcode(opcode)
	if err != nil {
		return -1, err
	}

	log.
		Trace().
		Int64("opcode", int64(opcode)).
		Msg("[COMPUTER] Parsed opcode")

	// Get the parameters for the opcode
	parametersLength := len(ic.opcodes[opcode].Parameters)
	opcodeParameters := ic.Memory.GetRange(ic.instructionPointer+1, int64(parametersLength))
	// TODO: fix this log
	// log.Trace().Ints("parameters", opcodeParameters).Msg("[COMPUTER] Retrieved opcode parameters")

	// Resolve the parameters based on the modes
	opcodeParameters, err = ic.resolveParameters(ic.Memory, opcode, opcodeParameters, parameterModes)
	if err != nil {
		return -1, err
	}

	// TODO: fix this log
	// log.Trace().Ints("parameters", opcodeParameters).Msg("[COMPUTER] Resolved opcode parameters")

	// Execute the opcode
	operation := ic.opcodes[opcode]
	log.Trace().Str("opcodeName", operation.Name).Msg("[COMPUTER] Executing operation")
	operation.execute(ic, operation, opcodeParameters)

	return opcode, nil
}

func (ic *Computer) Run() {
	ic.State = "running"

	for {
		opcode, err := ic.Step()
		if err != nil {
			ic.errorHandler(err)
		}

		// Special case for HALT
		if opcode == HALT {
			ic.State = "halted"

			log.Info().Str("name", ic.Name).Msg("[COMPUTER] Halt")

			close(ic.Input)
			close(ic.Output)

			break
		}
	}
}
