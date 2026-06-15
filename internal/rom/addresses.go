package rom

import (
	"fmt"
	"strconv"
	"strings"
)

func parseAddresses(filename, data string) map[uint32]string {
	mapping := make(map[uint32]string)
	for i, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		sep := strings.IndexFunc(line, func(r rune) bool { return r == ' ' || r == '\t' })
		if sep < 0 {
			panic(fmt.Sprintf("%s line %d: missing location name: %q", filename, i+1, line))
		}
		hex := strings.TrimPrefix(line[:sep], "0x")
		addr, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			panic(fmt.Sprintf("%s line %d: bad address %q: %v", filename, i+1, line[:sep], err))
		}
		name := strings.TrimSpace(line[sep:])
		mapping[uint32(addr)] = name
	}
	return mapping
}
