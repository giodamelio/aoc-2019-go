package intcode

type Memory struct {
	rawMemory []int
}

func newMemory(initialMemory []int) *Memory {
	mem := new(Memory)
	mem.rawMemory = initialMemory
	return mem
}

// Get the value of an address
func (im Memory) Get(location int) int {
	return im.rawMemory[location]
}

// Get the values from a range of addresses
func (im Memory) GetRange(location int, length int) []int {
	return im.rawMemory[location : location+length]
}

// Set the value of an address
func (im Memory) Set(location int, value int) {
	im.rawMemory[location] = value
}
