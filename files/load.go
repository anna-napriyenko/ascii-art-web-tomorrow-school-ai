package files

import (
	"fmt"
	"os"
	"strings"

	"ascii-art/printing"
)

func LoadBanner(filename string) (map[rune][]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read banner file %q: %w", filename, err)
	}

lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")

	bannerMap := make(map[rune][]string)
	currentChar := ' '
	for i := 1; currentChar <= '~'; i += (printing.ArtHeight + 1) {
		charArt := make([]string, printing.ArtHeight)
		for j := 0; j < printing.ArtHeight; j++ {
			if i+j < len(lines) {
				charArt[j] = lines[i+j]
			} else {
				charArt[j] = ""
			}
		}
		bannerMap[currentChar] = charArt
		currentChar++
		if i+printing.ArtHeight >= len(lines) {
			break
		}
	}
	return bannerMap, nil
}
