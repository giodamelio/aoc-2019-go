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
	memory := newMemory([]int{1, 0, 0, 0})

	opcodeAdd.execute(memory, memory.rawMemory[1:], nil, nil)

	assert.Equal(t, []int{2, 0, 0, 0}, memory.rawMemory)
}

func TestMultiply(t *testing.T) {
	opcodeMultiply := Opcodes[2]
	memory := newMemory([]int{2, 0, 0, 0})

	opcodeMultiply.execute(memory, memory.rawMemory[1:], nil, nil)

	assert.Equal(t, []int{4, 0, 0, 0}, memory.rawMemory)
}

func TestInput(t *testing.T) {
	opcodeInput := Opcodes[3]
	memory := newMemory([]int{3, 2, 0})

	input := make(chan int)

	go func() {
		input <- 10
	}()

	opcodeInput.execute(memory, memory.rawMemory[1:], input, nil)

	assert.Equal(t, []int{3, 2, 10}, memory.rawMemory)
}

func TestOutput(t *testing.T) {
	opcodeOutput := Opcodes[4]
	memory := newMemory([]int{4, 2, 10})

	output := make(chan int, 1)

	opcodeOutput.execute(memory, memory.rawMemory[1:], nil, output)

	assert.Equal(t, []int{4, 2, 10}, memory.rawMemory)
	assert.Equal(t, 10, <-output)
}

func TestHalt(t *testing.T) {
	opcodeHalt := Opcodes[99]
	memory := newMemory([]int{99})

	opcodeHalt.execute(memory, []int{}, nil, nil)

	assert.Equal(t, []int{99}, memory.rawMemory)
}
