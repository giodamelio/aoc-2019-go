package assembler

import (
	"testing"

	"github.com/giodamelio/aoc-2020-go/intcode"
	"github.com/stretchr/testify/assert"
)

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

	t.Fail()
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
