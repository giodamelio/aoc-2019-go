package intcode

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	out := zerolog.NewConsoleWriter()
	out.Out = os.Stderr
	out.NoColor = true
	log.Logger = log.Output(out)

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func TestNewComputer(t *testing.T) {
	computer := NewComputer([]AddressValue{1, 2, 3})

	assert.Equal(t, AddressLocation(0), computer.instructionPointer)
	assert.Equal(t, []AddressValue{1, 2, 3}, computer.Memory.rawMemory)
	assert.Equal(t, "computer", computer.Name)
	assert.Equal(t, "pre-run", computer.State)
}

func TestNewComputerNotModifyInitialMemory(t *testing.T) {
	program := []AddressValue{1101, 1, 2, 0, 99}
	computer := NewComputer(program)

	computer.Run()

	assert.Equal(t, []AddressValue{3, 1, 2, 0, 99}, computer.Memory.rawMemory)
	assert.Equal(t, []AddressValue{1101, 1, 2, 0, 99}, program)
	assert.Equal(t, "halted", computer.State)
}

func TestReverOutputModes(t *testing.T) {
	out := reverseOutputModes([]mode{Immediate, Immediate, Position})
	assert.Equal(t, []mode{Position, Immediate, Immediate}, out)
}

func TestParseOpcode(t *testing.T) {
	computer := NewComputer([]AddressValue{})

	opcodes := []AddressValue{2, 1002, 4, 99}
	for _, opcode := range opcodes {
		opcode, parameterModes, err := computer.parseOpcode(opcode)

		assert.Nil(t, err)
		assert.Equal(t, len(computer.opcodes[opcode].parameters), len(parameterModes))
	}
}

func TestParseOpcodeProperlyReversed(t *testing.T) {
	computer := NewComputer([]AddressValue{})

	_, parameterModes, err := computer.parseOpcode(1101)
	assert.Nil(t, err)
	assert.Equal(t, []mode{Immediate, Immediate, Position}, parameterModes)
}

func TestParseOpcodeInvalid(t *testing.T) {
	computer := NewComputer([]AddressValue{})

	_, _, err := computer.parseOpcode(1050)
	assert.Equal(t, "invalid opcode: 50", err.Error())
}

func TestResolveParameters(t *testing.T) {
	computer := NewComputer([]AddressValue{11002, 11, 11, 0})
	opcodeParameters, err := computer.resolveParameters(
		computer.Memory,
		2,
		[]AddressValue{11, 11, 0},
		[]mode{Immediate, Immediate, Position},
	)

	assert.Nil(t, err)
	assert.Equal(t, []AddressValue{11, 11, 0}, opcodeParameters)
}

func TestResolveParametersErrors(t *testing.T) {
	computer := NewComputer([]AddressValue{11102, 11, 11, 0})

	opcodeParameters, err := computer.resolveParameters(
		computer.Memory,
		2,
		[]AddressValue{11, 11, 0},
		[]mode{Immediate, Immediate, Immediate},
	)

	assert.Equal(t, "write parameter cannot be in immediate mode: 0", err.Error())
	assert.Nil(t, opcodeParameters)

	opcodeParameters, err = computer.resolveParameters(
		computer.Memory,
		2,
		[]AddressValue{11, 11, 0},
		[]mode{Immediate, Immediate, 2},
	)

	assert.Equal(t, "invalid mode: 2", err.Error())
	assert.Nil(t, opcodeParameters)

	// Create new opcode with bad parameter mode
	computer.opcodes[98] = Opcode{
		name:       "FAKE",
		opcode:     98,
		parameters: []readWrite{Read, Read, 2},
		execute:    func(computer *Computer, operation Opcode, parameters []AddressValue) {},
	}
	opcodeParameters, err = computer.resolveParameters(
		computer.Memory,
		98,
		[]AddressValue{11, 11, 0},
		[]mode{Immediate, Immediate, Position},
	)

	assert.Equal(t, "invalid parameter mode: 0", err.Error())
	assert.Nil(t, opcodeParameters)
}

func TestSetInstructionPointer(t *testing.T) {
	computer := NewComputer([]AddressValue{1, 0, 0, 0})

	assert.Equal(t, AddressLocation(0), computer.instructionPointer)

	computer.SetInstructionPointer(3)

	assert.Equal(t, AddressLocation(3), computer.instructionPointer)
}

func TestStep(t *testing.T) {
	computer := NewComputer([]AddressValue{1, 0, 0, 0})

	opcode, err := computer.Step()

	assert.Nil(t, err)
	assert.Equal(t, AddressValue(1), opcode)
	assert.Equal(t, []AddressValue{2, 0, 0, 0}, computer.Memory.rawMemory)
}

func TestStepInvalidOpcode(t *testing.T) {
	computer := NewComputer([]AddressValue{-1, 0, 0, 0})

	opcode, err := computer.Step()

	assert.Equal(t, AddressValue(-1), opcode)
	assert.Equal(t, "invalid opcode: -1", err.Error())
	assert.Equal(t, []AddressValue{-1, 0, 0, 0}, computer.Memory.rawMemory)
}

func TestStepNoIncrementInstructionPointer(t *testing.T) {
	computer := NewComputer([]AddressValue{99})

	opcode, err := computer.Step()

	assert.Nil(t, err)
	assert.Equal(t, AddressValue(99), opcode)
	assert.Equal(t, []AddressValue{99}, computer.Memory.rawMemory)
	assert.Equal(t, AddressLocation(0), computer.instructionPointer)
}

func TestRun(t *testing.T) {
	computer := NewComputer([]AddressValue{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50})

	computer.Run()

	assert.Equal(t, AddressValue(3500), computer.Memory.rawMemory[0])
}

func TestRunInvalidOpcode(t *testing.T) {
	computer := NewComputer([]AddressValue{-1})

	assert.PanicsWithError(t, "invalid opcode: -1", func() {
		computer.Run()
	})
}

func TestInputChannel(t *testing.T) {
	computer := NewComputer([]AddressValue{3, 3, 99, 0})

	sendInput := func() {
		computer.Input <- 10
	}
	go sendInput()

	computer.Run()

	assert.Equal(t, []AddressValue{3, 3, 99, 10}, computer.Memory.rawMemory)
}

func TestOutputChannel(t *testing.T) {
	computer := NewComputer([]AddressValue{104, 10, 99})

	wait := make(chan bool)

	// Listen for the output
	listenForOutput := func() {
		output := <-computer.Output
		assert.Equal(t, AddressValue(10), output)
		wait <- true
	}
	go listenForOutput()

	computer.Run()

	assert.Equal(
		t,
		[]AddressValue{104, 10, 99},
		computer.Memory.rawMemory,
	)

	// Make sure the output has been read
	<-wait
}

// Some test programs

// Add two numbers.
func TestAddTwoNumber(t *testing.T) {
	computer := NewComputer([]AddressValue{1101, 11, 22, 0, 99})

	computer.Run()

	assert.Equal(
		t,
		[]AddressValue{33, 11, 22, 0, 99},
		computer.Memory.rawMemory,
	)
}

// Take an input, double it and output it.
func TestDoubleInput(t *testing.T) {
	computer := NewComputer([]AddressValue{3, 0, 2, 2, 0, 0, 4, 0, 99})

	// Number to be doubled
	sendInput := func() {
		computer.Input <- 11
	}
	go sendInput()

	wait := make(chan bool)

	// Listen for the output
	listenForOutput := func() {
		output := <-computer.Output
		assert.Equal(t, AddressValue(22), output)
		wait <- true
	}
	go listenForOutput()

	computer.Run()

	assert.Equal(
		t,
		[]AddressValue{22, 0, 2, 2, 0, 0, 4, 0, 99},
		computer.Memory.rawMemory,
	)

	// Make sure the output has been read
	<-wait
}

// Test if the input is greater then zero.
func TestIsGreaterThenZero(t *testing.T) {
	computer := NewComputer([]AddressValue{
		// Program
		3, 12, //           INPUT					Read input to address 12
		6, 12, 15, //       JUMP-IF-FALSE	If the contents of address 12 are zero
		//                                jump to the location in address 15 (address 9)
		1, 13, 14, 13, //   ADD						Add the values from addresses 13 and 14 and put them in address 13
		4, 13, //           OUTPUT				Output the the value of address 13
		99, //              HALT

		// Data
		-1, // Address 12
		0,  //         13
		1,  //         14
		9,  //         15
	})

	sendInput := func() {
		computer.Input <- 22
	}
	go sendInput()

	wait := make(chan bool)

	// Listen for the output
	listenForOutput := func() {
		output := <-computer.Output
		assert.Equal(t, AddressValue(1), output)
		wait <- true
	}
	go listenForOutput()

	computer.Run()

	// Make sure the output has been read
	<-wait
}
