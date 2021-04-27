package intcode

func copyMemory(input []int) []int {
	inputCopy := make([]int, len(input))
	copy(inputCopy, input)
	return inputCopy
}
