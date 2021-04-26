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

	opcodeAdd.execute(memory, memory.rawMemory[1:], nil)

	assert.Equal(t, []int{2, 0, 0, 0}, memory.rawMemory)
}

func TestMultiply(t *testing.T) {
	opcodeMultiply := Opcodes[2]
	memory := newMemory([]int{2, 0, 0, 0})

	opcodeMultiply.execute(memory, memory.rawMemory[1:], nil)

	assert.Equal(t, []int{4, 0, 0, 0}, memory.rawMemory)
}

func TestInput(t *testing.T) {
	opcodeInput := Opcodes[3]
	memory := newMemory([]int{3, 2, 0})

	input := make(chan int)

	go func() {
		input <- 10
	}()

	opcodeInput.execute(memory, memory.rawMemory[1:], input)

	assert.Equal(t, []int{3, 2, 10}, memory.rawMemory)
}

func TestHalt(t *testing.T) {
	opcodeHalt := Opcodes[99]
	memory := newMemory([]int{99})

	opcodeHalt.execute(memory, []int{}, nil)

	assert.Equal(t, []int{99}, memory.rawMemory)
}
