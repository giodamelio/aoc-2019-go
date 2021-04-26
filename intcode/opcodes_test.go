package intcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	opcodeAdd := Opcodes[1]
	memory := newMemory([]int{1, 0, 0, 0})

	opcodeAdd.execute(memory, memory.rawMemory[1:])

	assert.Equal(t, []int{2, 0, 0, 0}, memory.rawMemory)
}

func TestMultiply(t *testing.T) {
	opcodeMultiply := Opcodes[2]
	memory := newMemory([]int{2, 0, 0, 0})

	opcodeMultiply.execute(memory, memory.rawMemory[1:])

	assert.Equal(t, []int{4, 0, 0, 0}, memory.rawMemory)
}

func TestHalt(t *testing.T) {
	opcodeHalt := Opcodes[99]
	memory := newMemory([]int{99})

	opcodeHalt.execute(memory, []int{})

	assert.Equal(t, []int{99}, memory.rawMemory)
}
