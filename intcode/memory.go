package intcode

import "github.com/rs/zerolog/log"

type Memory struct {
	rawMemory []int
}

func newMemory(initialMemory []int) *Memory {
	log.Trace().Msg("[MEMORY] Memory Created")

	mem := new(Memory)
	mem.rawMemory = initialMemory
	return mem
}

// Get the value of an address
func (im Memory) Get(address int) int {
	value := im.rawMemory[address]

	log.
		Trace().
		Int("address", address).
		Int("value", value).
		Msg("[MEMORY] Get")

	return value
}

// Get the values from a range of addresses
func (im Memory) GetRange(address int, length int) []int {
	value := im.rawMemory[address : address+length]

	log.
		Trace().
		Int("address", address).
		Ints("value", value).
		Msg("[MEMORY] GetRange")

	return value
}

// Set the value of an address
func (im Memory) Set(address int, value int) {
	log.
		Trace().
		Int("address", address).
		Int("value", value).
		Int("oldvalue", im.rawMemory[address]).
		Msg("[MEMORY] Set")

	im.rawMemory[address] = value
}
