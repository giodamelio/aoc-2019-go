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

func TestLength(t *testing.T) {
	opcodeAdd := Opcodes[1]
	opcodeHalt := Opcodes[99]

	assert.Equal(t, 4, opcodeAdd.length())
	assert.Equal(t, 1, opcodeHalt.length())
}

func TestIncrementInstructionPointer(t *testing.T) {
	opcodeAdd := Opcodes[1]
	computer := NewComputer([]int{1, 1, 1, 0})

	assert.Equal(t, 0, computer.instructionPointer)

	opcodeAdd.incrementInstructionPointer(computer)

	assert.Equal(t, 4, computer.instructionPointer)
}

func TestAdd(t *testing.T) {
	opcodeAdd := Opcodes[1]
	computer := NewComputer([]int{1, 1, 1, 0})

	opcodeAdd.execute(computer, opcodeAdd, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{2, 1, 1, 0}, computer.Memory.rawMemory)
}

func TestMultiply(t *testing.T) {
	opcodeMultiply := Opcodes[2]
	computer := NewComputer([]int{2, 2, 2, 0})

	opcodeMultiply.execute(computer, opcodeMultiply, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{4, 2, 2, 0}, computer.Memory.rawMemory)
}

func TestInput(t *testing.T) {
	opcodeInput := Opcodes[3]
	computer := NewComputer([]int{3, 2, 0})

	input := make(chan int)
	computer.Input = input

	go func() {
		input <- 10
	}()

	opcodeInput.execute(computer, opcodeInput, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{3, 2, 10}, computer.Memory.rawMemory)
}

func TestOutput(t *testing.T) {
	opcodeOutput := Opcodes[4]
	computer := NewComputer([]int{4, 10, 10})

	output := make(chan int, 1)
	computer.Output = output

	opcodeOutput.execute(computer, opcodeOutput, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{4, 10, 10}, computer.Memory.rawMemory)
	assert.Equal(t, 10, <-output)
}

func TestJumpIfTrue(t *testing.T) {
	opcodeJumpIfTrue := Opcodes[5]

	// Jump case
	computer := NewComputer([]int{1105, 1, 22, 99})

	opcodeJumpIfTrue.execute(computer, opcodeJumpIfTrue, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{1105, 1, 22, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 22, computer.instructionPointer)

	// Continue case
	computer = NewComputer([]int{1105, 0, 22, 99})

	opcodeJumpIfTrue.execute(computer, opcodeJumpIfTrue, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{1105, 0, 22, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 3, computer.instructionPointer)
}

func TestJumpIfFalse(t *testing.T) {
	opcodeJumpIfFalse := Opcodes[6]

	// Jump case
	computer := NewComputer([]int{1106, 0, 22, 99})

	opcodeJumpIfFalse.execute(computer, opcodeJumpIfFalse, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{1106, 0, 22, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 22, computer.instructionPointer)

	// Continue case
	computer = NewComputer([]int{1106, 1, 22, 99})

	opcodeJumpIfFalse.execute(computer, opcodeJumpIfFalse, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{1106, 1, 22, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 3, computer.instructionPointer)
}

func TestLessThan(t *testing.T) {
	opcodeLessThan := Opcodes[7]

	// Test less then case
	computer := NewComputer([]int{1107, 10, 20, 0, 99})

	opcodeLessThan.execute(computer, opcodeLessThan, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{1, 10, 20, 0, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 4, computer.instructionPointer)

	// Test not less then case
	computer = NewComputer([]int{1107, 20, 10, 0, 99})

	opcodeLessThan.execute(computer, opcodeLessThan, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{0, 20, 10, 0, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 4, computer.instructionPointer)
}

func TestEquals(t *testing.T) {
	opcodeEquals := Opcodes[8]

	// Test equals
	computer := NewComputer([]int{1108, 10, 10, 0, 99})

	opcodeEquals.execute(computer, opcodeEquals, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{1, 10, 10, 0, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 4, computer.instructionPointer)

	// Test not equals
	computer = NewComputer([]int{1108, 10, 20, 0, 99})

	opcodeEquals.execute(computer, opcodeEquals, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{0, 10, 20, 0, 99}, computer.Memory.rawMemory)
	assert.Equal(t, 4, computer.instructionPointer)
}

func TestHalt(t *testing.T) {
	opcodeHalt := Opcodes[99]
	computer := NewComputer([]int{99})

	opcodeHalt.execute(computer, opcodeHalt, []int{})

	assert.Equal(t, []int{99}, computer.Memory.rawMemory)
}
