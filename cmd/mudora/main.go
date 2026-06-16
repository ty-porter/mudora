package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ty-porter/mudora/internal"
	"github.com/ty-porter/mudora/internal/alttp"
	"github.com/ty-porter/mudora/internal/rom"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("mudora", internal.Version)
		return
	}
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mudora <rom.sfc> [item-query]")
		fmt.Fprintln(os.Stderr, "       mudora --version")
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "mudora:", err)
		os.Exit(1)
	}

	query := strings.Join(os.Args[2:], " ")
	for _, g := range alttp.Filter(alttp.Grouped(rom.Inspect(data)), query) {
		fmt.Println(g.Region)
		for _, p := range g.Locations {
			fmt.Printf("  %-56s %s\n", p.Location, p.Item)
		}
		fmt.Println()
	}
}
