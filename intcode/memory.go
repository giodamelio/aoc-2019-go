package intcode

import "github.com/rs/zerolog/log"

type AddressLocation int64
type AddressValue int64

type Memory struct {
	rawMemory []AddressValue
}

func newMemory(initialMemory []AddressValue) *Memory {
	log.Trace().Msg("[MEMORY] Memory Created")

	mem := new(Memory)
	mem.rawMemory = initialMemory

	return mem
}

// Get the value of an address.
func (im Memory) Get(address AddressLocation) AddressValue {
	value := im.rawMemory[address]

	log.
		Trace().
		Int64("address", int64(address)).
		Int64("value", int64(value)).
		Msg("[MEMORY] Get")

	return value
}

// Get the values from a range of addresses.
func (im Memory) GetRange(address AddressLocation, length int64) []AddressValue {
	value := im.rawMemory[address : int64(address)+int64(length)]

	log.
		Trace().
		Int64("address", int64(address)).
		Int64("length", int64(length)).
		// TODO: fix this log
		// Ints("value", value).
		Msg("[MEMORY] GetRange")

	return value
}

// Set the value of an address.
func (im Memory) Set(address AddressValue, value AddressValue) {
	log.
		Trace().
		Int64("address", int64(address)).
		Int64("value", int64(value)).
		Int64("oldvalue", int64(im.rawMemory[address])).
		Msg("[MEMORY] Set")

	im.rawMemory[address] = value
}
