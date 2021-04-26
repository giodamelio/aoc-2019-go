package intcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComputer(t *testing.T) {
	computer := NewComputer([]int{1, 2, 3})

	assert.Equal(t, 0, computer.programCounter)
	assert.Equal(t, []int{1, 2, 3}, computer.memory.rawMemory)
}

func TestStep(t *testing.T) {
	computer := NewComputer([]int{1, 0, 0, 0})

	opcode, err := computer.Step()

	assert.Nil(t, err)
	assert.Equal(t, 1, opcode)
	assert.Equal(t, []int{2, 0, 0, 0}, computer.memory.rawMemory)
}

func TestStepInvalidOpcode(t *testing.T) {
	computer := NewComputer([]int{-1, 0, 0, 0})

	opcode, err := computer.Step()

	assert.Equal(t, -1, opcode)
	assert.Equal(t, "invalid opcode: -1", err.Error())
	assert.Equal(t, []int{-1, 0, 0, 0}, computer.memory.rawMemory)
}

func TestRun(t *testing.T) {
	computer := NewComputer([]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50})

	computer.Run()

	assert.Equal(t, 3500, computer.memory.rawMemory[0])
}

func TestRunInvalidOpcode(t *testing.T) {
	computer := NewComputer([]int{-1})

	assert.PanicsWithError(t, "invalid opcode: -1", func() {
		computer.Run()
	})
}
