package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/anishnarang9/eth-indexer/internal/event"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: eventcheck <events.jsonl>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, "open file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxLineSize = 1024 * 1024
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	var processed, valid, invalid int

	for scanner.Scan() {
		processed++
		e, err := event.Decode(scanner.Bytes())
		if err == nil {
			err = event.Validate(e)
		}
		if err != nil {
			invalid++
			fmt.Printf("Line %d: %v\n", processed, err)
			continue
		}
		valid++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read file:", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records: %d valid, %d invalid\n",
		processed, valid, invalid)

	if invalid > 0 {
		os.Exit(1)
	}
}
