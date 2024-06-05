package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bncode"
	"github.com/codecrafters-io/bittorrent-starter-go/internal/torrent"
)

func main() {
	command := os.Args[1]

	if command == "decode" {

		bencodedValue := os.Args[2]

		decoded, err := bncode.Decode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else if command == "info" {
		f := os.Args[2]
		m, h, err := torrent.ReadMetaData(f)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Tracker URL: %s\n", m.Announce)
		fmt.Printf("Length: %d\n", m.Info.Length)
		fmt.Printf("Info Hash: %s", h)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
