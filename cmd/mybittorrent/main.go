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
		info, err := torrent.Read(f)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Tracker URL: %s\n", info["announce"].(string))
		fmt.Printf("Length: %d", info["info"].(map[string]interface{})["length"].(int))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
