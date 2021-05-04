package intcode

import (
	"strconv"
	"strings"
)

func copyMemory(input []AddressValue) []AddressValue {
	inputCopy := make([]AddressValue, len(input))
	copy(inputCopy, input)

	return inputCopy
}

func ParseInput(input string) ([]AddressValue, error) {
	bytes := strings.Split(strings.TrimSpace(input), ",")
	numbers := make([]AddressValue, len(bytes))

	for index, programByte := range bytes {
		number, err := strconv.Atoi(programByte)
		if err != nil {
			return nil, err
		}

		numbers[index] = AddressValue(number)
	}

	return numbers, nil
}
