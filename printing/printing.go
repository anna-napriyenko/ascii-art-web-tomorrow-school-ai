package printing

import "fmt"

const ArtHeight = 8

func PrintAsciiArt(text string, banner map[rune][]string) {
	outputLines := make([]string, ArtHeight)

	for _, char := range text {
		if art, ok := banner[char]; ok {
			for i := 0; i < ArtHeight; i++ {
				if i < len(art) {
					outputLines[i] += art[i]
				}
			}
		}
	}

	for _, line := range outputLines {
		fmt.Println(line)
	}
}
