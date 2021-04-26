package intcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
