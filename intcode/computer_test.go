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

func TestSendInput(t *testing.T) {
	computer := NewComputer([]int{3, 3, 99, 0})

	computer.SendInput(10)
	computer.Run()

	assert.Equal(t, []int{3, 3, 99, 10}, computer.Memory.rawMemory)
}

// Test GetOutputChannel by taking an input, doubling it and outputing the result
func TestGetOutputChannel(t *testing.T) {
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
