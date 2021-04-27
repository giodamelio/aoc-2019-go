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
	computer.input = input

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
	computer.output = output

	opcodeOutput.execute(computer, opcodeOutput, computer.Memory.rawMemory[1:])

	assert.Equal(t, []int{4, 10, 10}, computer.Memory.rawMemory)
	assert.Equal(t, 10, <-output)
}

func TestHalt(t *testing.T) {
	opcodeHalt := Opcodes[99]
	computer := NewComputer([]int{99})

	opcodeHalt.execute(computer, opcodeHalt, []int{})

	assert.Equal(t, []int{99}, computer.Memory.rawMemory)
}
