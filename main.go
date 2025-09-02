package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: minitorrent <torrent-file>")
		os.Exit(1)
	}
	torrentFile := os.Args[1]
	fmt.Printf("Starting download for: %s\n", torrentFile)
}
