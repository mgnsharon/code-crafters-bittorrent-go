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
		m, err := torrent.ReadMetaData(f)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Tracker URL: %s\n", m.Announce)
		fmt.Printf("Length: %d\n", m.Info.Length)
		fmt.Printf("Info Hash: %s\n", m.GetInfoHash())
		fmt.Printf("Piece Length: %d\n", m.Info.PieceLength)
		fmt.Println("Piece Hashes:")
		hashes, err := m.Info.Hashes()
		if err != nil {
			fmt.Println(err)
		}
		for _, h := range hashes {
			fmt.Println(h)
		}
	} else if command == "peers" {
		f := os.Args[2]
		m, err := torrent.ReadMetaData(f)
		if err != nil {
			fmt.Println(err)
		}
		resp, err := torrent.DiscoverPeers(&m)
		if err != nil {
			fmt.Println(err)
		}
		for _, p := range resp.Peers {
			fmt.Printf("%d.%d.%d.%d:%d\n", p.IP[0], p.IP[1], p.IP[2], p.IP[3], p.Port)
		}
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
