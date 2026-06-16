package main

import (
	"fmt"
	"os"

	"github.com/ty-porter/mudora/internal/alttp"
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

	for _, g := range alttp.Grouped(rom.Inspect(data)) {
		fmt.Println(g.Region)
		for _, p := range g.Locations {
			fmt.Printf("  %-56s %s\n", p.Location, p.Item)
		}
		fmt.Println()
	}
}
