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

func TestGet(t *testing.T) {
	computer := NewComputer([]int{1, 2, 3})

	assert.Equal(t, 2, computer.Memory.Get(1))
}

func TestGetRange(t *testing.T) {
	computer := NewComputer([]int{1, 2, 3, 4, 5, 6})

	assert.Equal(t, []int{3, 4, 5}, computer.Memory.GetRange(2, 3))
}

func TestSet(t *testing.T) {
	computer := NewComputer([]int{1, 2, 3})

	computer.Memory.Set(1, 10)
	assert.Equal(t, 10, computer.Memory.Get(1))
}
