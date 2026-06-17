//go:build js && wasm

package main

import (
	"encoding/base64"
	"encoding/json"
	"syscall/js"

	"github.com/alttpr-mudora/mudora/internal"
	"github.com/alttpr-mudora/mudora/internal/alttp"
	"github.com/alttpr-mudora/mudora/internal/icons"
	"github.com/alttpr-mudora/mudora/internal/rom"
)

type placement struct {
	Location    string `json:"location"`
	Item        string `json:"item"`
	Icon        string `json:"icon,omitempty"`
	Progression bool   `json:"progression"`
}

type group struct {
	Region    string      `json:"region"`
	Locations []placement `json:"locations"`
}

func main() {
	js.Global().Set("mudoraInspect", js.FuncOf(inspect))
	js.Global().Set("mudoraVersion", js.ValueOf(internal.Version))
	select {}
}

func inspect(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errResult("missing ROM bytes")
	}

	romBytes := args[0]
	data := make([]byte, romBytes.Get("length").Int())
	js.CopyBytesToGo(data, romBytes)

	query := ""
	if len(args) > 1 {
		query = args[1].String()
	}

	groups := alttp.Filter(alttp.Grouped(rom.Inspect(data)), query)
	out := make([]group, 0, len(groups))
	for _, g := range groups {
		og := group{Region: g.Region, Locations: make([]placement, 0, len(g.Locations))}
		for _, p := range g.Locations {
			pl := placement{Location: p.Location, Item: p.Item, Progression: alttp.IsProgression(p.Item)}
			if png, ok := icons.PNG(p.Item); ok {
				pl.Icon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
			}
			og.Locations = append(og.Locations, pl)
		}
		out = append(out, og)
	}

	b, err := json.Marshal(out)
	if err != nil {
		return errResult(err.Error())
	}
	return string(b)
}

func errResult(msg string) any {
	b, _ := json.Marshal(map[string]string{"error": msg})
	return string(b)
}
