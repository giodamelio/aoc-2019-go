package assembler

import (
	"os"
	"testing"

	"github.com/giodamelio/aoc-2020-go/intcode"
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

func TestPowInt(t *testing.T) {
	assert.Equal(t, 256, powInt(4, 4))
	assert.Equal(t, 4, powInt(2, 2))
	assert.Equal(t, 144, powInt(12, 2))
	assert.Equal(t, 1, powInt(2, 0))
}

func TestParseArgument(t *testing.T) {
	argument, mode := parseArgument("10")

	assert.Equal(t, 10, argument)
	assert.Equal(t, byte('0'), mode)

	argument, mode = parseArgument("i10")

	assert.Equal(t, 10, argument)
	assert.Equal(t, byte('1'), mode)
}

func TestSimpleHaltProgram(t *testing.T) {
	program := Assemble("HALT")

	assert.Equal(t, []intcode.AddressValue{99}, program)
}

func TestOpcodeWithParameters(t *testing.T) {
	program := Assemble("ADD	0	0	0")

	assert.Equal(t, []intcode.AddressValue{1, 0, 0, 0}, program)
}

func TestMultipleInstructions(t *testing.T) {
	program := Assemble(`
	ADD	0	0	0
	HALT
	`)

	assert.Equal(t, []intcode.AddressValue{1, 0, 0, 0, 99}, program)
}

func TestMultipleTabs(t *testing.T) {
	program := Assemble(`
	ADD					0	0	0
	MULTIPLY		0	0	0
	HALT
	`)

	assert.Equal(t, []intcode.AddressValue{1, 0, 0, 0, 2, 0, 0, 0, 99}, program)
}

func TestArgumentModes(t *testing.T) {
	program := Assemble(`
	ADD	i10	i10	0
	HALT
	`)

	assert.Equal(t, []intcode.AddressValue{1101, 10, 10, 0, 99}, program)
}

func TestData(t *testing.T) {
	program := Assemble(`
	ADD	i10	i10	0
	HALT
	DATA	10
	`)

	assert.Equal(t, []intcode.AddressValue{1101, 10, 10, 0, 99, 10}, program)
}

// Test some more complicated programs.
func TestAddTwoNumber(t *testing.T) {
	computer := intcode.NewComputer(Assemble(`
	ADD		i11	i22	0
	HALT
	`))

	computer.Run()

	assert.Equal(
		t,
		intcode.AddressValue(33),
		computer.Memory.Get(0),
	)
}

func TestIsGreaterThenZero(t *testing.T) {
	computer := intcode.NewComputer(Assemble(`
	INPUT		12
	JUMP-IF-FALSE	12	15
	ADD			13	14	13
	OUTPUT	13
	HALT
	DATA	-1
	DATA	0
	DATA	1
	DATA	9
	`))

	sendInput := func() {
		computer.Input <- 22
	}
	go sendInput()

	wait := make(chan bool)

	// Listen for the output
	listenForOutput := func() {
		output := <-computer.Output
		assert.Equal(t, intcode.AddressValue(1), output)
		wait <- true
	}
	go listenForOutput()

	computer.Run()

	// Make sure the output has been read
	<-wait
}

// Take an input, double it and output it.
func TestDoubleInput(t *testing.T) {
	computer := intcode.NewComputer(Assemble(`
	INPUT	0
	MULTIPLY	0	i2	0
	OUTPUT	0
	HALT
	`))

	// Number to be doubled
	sendInput := func() {
		computer.Input <- 11
	}
	go sendInput()

	wait := make(chan bool)

	// Listen for the output
	listenForOutput := func() {
		output := <-computer.Output
		assert.Equal(t, intcode.AddressValue(22), output)
		wait <- true
	}
	go listenForOutput()

	computer.Run()

	// Make sure the output has been read
	<-wait
}
