package main

import (
	"hangman/hangman"
	"os"
)

func main() {
	var asciiOk bool
	if len(os.Args) == 4 {
		asciiOk = hangman.AsciiArtok()
	}
	if len(os.Args) == 3 {
		hangman.Load()
	}
	hangman.IntroHang(asciiOk)
}
