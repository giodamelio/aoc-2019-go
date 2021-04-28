package intcode

import (
	"strconv"
	"strings"
)

func copyMemory(input []int) []int {
	inputCopy := make([]int, len(input))
	copy(inputCopy, input)

	return inputCopy
}

func ParseInput(input string) ([]int, error) {
	bytes := strings.Split(strings.TrimSpace(input), ",")
	numbers := make([]int, len(bytes))

	for index, programByte := range bytes {
		number, err := strconv.Atoi(programByte)
		if err != nil {
			return nil, err
		}

		numbers[index] = number
	}

	return numbers, nil
}
