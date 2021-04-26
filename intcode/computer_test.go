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

	assert.Equal(t, 0, computer.programCounter)
	assert.Equal(t, []int{1, 2, 3}, computer.Memory.rawMemory)
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
