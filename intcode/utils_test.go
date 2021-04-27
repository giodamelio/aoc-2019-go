package intcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyMemory(t *testing.T) {
	old := []int{1, 2, 3}
	new := copyMemory(old)

	new[0] = 4

	assert.Equal(t, []int{1, 2, 3}, old)
	assert.Equal(t, []int{4, 2, 3}, new)
}
