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

func TestParseInput(t *testing.T) {
	parsedInput, err := ParseInput("1,2,3")

	assert.Nil(t, err, "Parsing returned an error")
	assert.Equal(t, []int{1, 2, 3}, parsedInput)
}

func TestParseInputNonInteger(t *testing.T) {
	parsedInput, err := ParseInput("1,haha,3")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "strconv.Atoi: parsing \"haha\": invalid syntax")
	assert.Nil(t, parsedInput)
}
