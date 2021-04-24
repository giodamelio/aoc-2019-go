package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	log := log.New(os.Stdout, "", 0)

	// Parse input flags
	dayPtr := flag.Int("day", 1, "The day you want to run")
	partPtr := flag.String("part", "both", "The part you want to run \"1\", \"2\" or \"both\"")

	flag.Parse()

	// Ensure a valid day
	if *dayPtr < 1 || *dayPtr > 25 {
		log.Fatal("Day must be between 1 and 25")
	}

	// Ensure a valid part
	if *partPtr != "1" && *partPtr != "2" && *partPtr != "both" {
		log.Fatal("Part must be 1, 2, or both")
	}

	fmt.Println("day:", *dayPtr)
	fmt.Println("part:", *partPtr)
}
