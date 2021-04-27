package intcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyMemory(t *testing.T) {
	oldSlice := []int{1, 2, 3}
	newSlice := copyMemory(oldSlice)

	newSlice[0] = 4

	assert.Equal(t, []int{1, 2, 3}, oldSlice)
	assert.Equal(t, []int{4, 2, 3}, newSlice)
}
