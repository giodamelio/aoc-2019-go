package assembler

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/giodamelio/aoc-2020-go/intcode"
)

func powInt(x int, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

// Return the raw argument and it's mode.
func parseArgument(argument string) (int, byte) {
	if argument[0] == 'i' {
		arg, err := strconv.Atoi(argument[1:])
		if err != nil {
			panic(err)
		}

		return arg, '1'
	}

	arg, err := strconv.Atoi(argument)
	if err != nil {
		panic(err)
	}

	return arg, '0'
}

func Assemble(programRaw string) []intcode.AddressValue {
	lines := strings.Split(strings.TrimSpace(programRaw), "\n")
	program := make([]intcode.AddressValue, 0, len(lines))

	// Build a map of TEXT to OPCODE mappings from the OPCODE to TEXT map
	opcodeMap := make(map[string]intcode.Opcode)
	for _, opcodeDetails := range intcode.Opcodes {
		opcodeMap[opcodeDetails.Name] = opcodeDetails
	}

	for _, line := range lines {
		// Compact multiple tabs into one tab
		re := regexp.MustCompile(`\t+`)
		lineSingleTabs := re.ReplaceAllString(line, "\t")

		// Split the line up by tabs
		sections := strings.Split(strings.TrimSpace(lineSingleTabs), "\t")

		// Get opcode details using the first section
		opcode := opcodeMap[sections[0]]

		// Get the count of arguments
		argumentCount := len(sections[1:])

		// Parse the args and their modes
		args := make([]intcode.AddressValue, argumentCount)
		modes := make([]byte, argumentCount)

		for index, argument := range sections[1:] {
			arg, mode := parseArgument(argument)

			args[index] = intcode.AddressValue(arg)
			modes[argumentCount-index-1] = mode
		}

		// Add argument modes to opcode string
		stringOpcode := fmt.Sprintf("%s%02d", string(modes), opcode.Opcode)

		// Convert opcode string to int
		numOpcode, err := strconv.Atoi(stringOpcode)
		if err != nil {
			panic(err)
		}

		// Add the opcode to the program
		program = append(program, intcode.AddressValue(numOpcode))

		// Add the arguments
		program = append(program, args...)
	}

	return program
}
