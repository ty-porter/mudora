package main

import (
	"fmt"
	"os"

	"github.com/ty-porter/mudora/internal/rom"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mudora <rom.sfc>")
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "mudora:", err)
		os.Exit(1)
	}

	for _, e := range rom.Inspect(data) {
		fmt.Printf("%-60s %s\n", e.Location, e.Item)
	}
}
