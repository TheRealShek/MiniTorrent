package main

import (
	"fmt"
	"os"

	"minitorrent/bencode"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: minitorrent <torrent-file>")
		os.Exit(1)
	}
	torrentFile := os.Args[1]
	fmt.Printf("Starting download for: %s\n", torrentFile)

	// Parse the torrent file
	torrent, err := bencode.ParseTorrentFile(torrentFile)
	if err != nil {
		fmt.Printf("Error parsing torrent file: %v\n", err)
		os.Exit(1)
	}

	// Display torrent information
	fmt.Printf("Name: %s\n", torrent.Name)
	fmt.Printf("Announce URL: %s\n", torrent.Announce)
	fmt.Printf("File Size: %d bytes\n", torrent.Length)
	fmt.Printf("Piece Length: %d bytes\n", torrent.PieceLength)
	fmt.Printf("Number of Pieces: %d\n", len(torrent.PieceHashes))
	fmt.Printf("Info Hash: %x\n", torrent.InfoHash)
}
