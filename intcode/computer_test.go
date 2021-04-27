package intcode

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func TestNewComputer(t *testing.T) {
	computer := NewComputer([]int{1, 2, 3})

	assert.Equal(t, 0, computer.instructionPointer)
	assert.Equal(t, []int{1, 2, 3}, computer.Memory.rawMemory)
}

func TestReverOutputModes(t *testing.T) {
	out := reverseOutputModes([]mode{Immediate, Immediate, Position})
	assert.Equal(t, []mode{Position, Immediate, Immediate}, out)
}

func TestParseOpcode(t *testing.T) {
	computer := NewComputer([]int{})

	opcodes := []int{2, 1002, 4, 99}
	for _, opcode := range opcodes {
		opcode, parameterModes, err := computer.parseOpcode(opcode)

		assert.Nil(t, err)
		assert.Equal(t, len(computer.opcodes[opcode].parameters), len(parameterModes))
	}
}

func TestParseOpcodeProperlyReversed(t *testing.T) {
	computer := NewComputer([]int{})

	_, parameterModes, err := computer.parseOpcode(1101)
	assert.Nil(t, err)
	assert.Equal(t, []mode{Immediate, Immediate, Position}, parameterModes)
}

func TestParseOpcodeInvalid(t *testing.T) {
	computer := NewComputer([]int{})

	_, _, err := computer.parseOpcode(1050)
	assert.Equal(t, "invalid opcode: 50", err.Error())
}

func TestResolveParameters(t *testing.T) {
	computer := NewComputer([]int{11002, 11, 11, 0})
	opcodeParameters, err := computer.resolveParameters(
		computer.Memory,
		2,
		[]int{11, 11, 0},
		[]mode{Immediate, Immediate, Position},
	)

	assert.Nil(t, err)
	assert.Equal(t, []int{11, 11, 0}, opcodeParameters)
}

func TestResolveParametersErrors(t *testing.T) {
	computer := NewComputer([]int{11102, 11, 11, 0})

	opcodeParameters, err := computer.resolveParameters(
		computer.Memory,
		2,
		[]int{11, 11, 0},
		[]mode{Immediate, Immediate, Immediate},
	)

	assert.Equal(t, "write parameter cannot be in immediate mode: 0", err.Error())
	assert.Nil(t, opcodeParameters)

	opcodeParameters, err = computer.resolveParameters(
		computer.Memory,
		2,
		[]int{11, 11, 0},
		[]mode{Immediate, Immediate, 2},
	)

	assert.Equal(t, "invalid mode: 2", err.Error())
	assert.Nil(t, opcodeParameters)

	// Create new opcode with bad parameter mode
	computer.opcodes[98] = Opcode{
		name:       "FAKE",
		opcode:     98,
		parameters: []readWrite{Read, Read, 2},
		execute:    func(computer *Computer, operation Opcode, parameters []int) {},
	}
	opcodeParameters, err = computer.resolveParameters(
		computer.Memory,
		98,
		[]int{11, 11, 0},
		[]mode{Immediate, Immediate, Position},
	)

	assert.Equal(t, "invalid parameter mode: 0", err.Error())
	assert.Nil(t, opcodeParameters)
}

func TestSetInstructionPointer(t *testing.T) {
	computer := NewComputer([]int{1, 0, 0, 0})

	assert.Equal(t, 0, computer.instructionPointer)

	computer.SetInstructionPointer(3)

	assert.Equal(t, 3, computer.instructionPointer)
}

func TestStep(t *testing.T) {
	computer := NewComputer([]int{1, 0, 0, 0})

	opcode, err := computer.Step()

	assert.Nil(t, err)
	assert.Equal(t, 1, opcode)
	assert.Equal(t, []int{2, 0, 0, 0}, computer.Memory.rawMemory)
}

func TestStepInvalidOpcode(t *testing.T) {
	computer := NewComputer([]int{-1, 0, 0, 0})

	opcode, err := computer.Step()

	assert.Equal(t, -1, opcode)
	assert.Equal(t, "invalid opcode: -1", err.Error())
	assert.Equal(t, []int{-1, 0, 0, 0}, computer.Memory.rawMemory)
}

func TestStepNoIncrementInstructionPointer(t *testing.T) {
	computer := NewComputer([]int{99})

	opcode, err := computer.Step()

	assert.Nil(t, err)
	assert.Equal(t, 99, opcode)
	assert.Equal(t, []int{99}, computer.Memory.rawMemory)
	assert.Equal(t, 0, computer.instructionPointer)
}

func TestRun(t *testing.T) {
	computer := NewComputer([]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50})

	computer.Run()

	assert.Equal(t, 3500, computer.Memory.rawMemory[0])
}

func TestRunInvalidOpcode(t *testing.T) {
	computer := NewComputer([]int{-1})

	assert.PanicsWithError(t, "invalid opcode: -1", func() {
		computer.Run()
	})
}

func TestSendInput(t *testing.T) {
	computer := NewComputer([]int{3, 3, 99, 0})

	computer.SendInput(10)
	computer.Run()

	assert.Equal(t, []int{3, 3, 99, 10}, computer.Memory.rawMemory)
}

func TestGetOutputChannel(t *testing.T) {
	computer := NewComputer([]int{104, 10, 99})

	// Listen for the output
	wait := make(chan bool)
	go func() {
		output := <-computer.GetOutputChannel()
		assert.Equal(t, 10, output)
		wait <- true
	}()

	computer.Run()

	assert.Equal(
		t,
		[]int{104, 10, 99},
		computer.Memory.rawMemory,
	)

	// Make sure the output has been read
	<-wait
}

// Some test programs

// Add two numbers
func TestAddTwoNumber(t *testing.T) {
	computer := NewComputer([]int{1101, 11, 22, 0, 99})

	computer.Run()

	assert.Equal(
		t,
		[]int{33, 11, 22, 0, 99},
		computer.Memory.rawMemory,
	)
}

// Take an input, double it and output it
func TestDoubleInput(t *testing.T) {
	computer := NewComputer([]int{3, 0, 2, 2, 0, 0, 4, 0, 99})

	// Number to be doubled
	computer.SendInput(11)

	// Listen for the output
	wait := make(chan bool)
	go func() {
		output := <-computer.GetOutputChannel()
		assert.Equal(t, 22, output)
		wait <- true
	}()

	computer.Run()

	assert.Equal(
		t,
		[]int{22, 0, 2, 2, 0, 0, 4, 0, 99},
		computer.Memory.rawMemory,
	)

	// Make sure the output has been read
	<-wait
}

// Test if the input is greater then zero
func TestIsGreaterThenZero(t *testing.T) {
	computer := NewComputer([]int{
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

	computer.SendInput(22)

	// Listen for the output
	wait := make(chan bool)
	go func() {
		output := <-computer.GetOutputChannel()
		assert.Equal(t, 1, output)
		wait <- true
	}()

	computer.Run()

	// Make sure the output has been read
	<-wait
}
