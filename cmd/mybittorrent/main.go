package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bncode"
	"github.com/codecrafters-io/bittorrent-starter-go/internal/torrent"
)

const (
	PEER_ID = "00112233445566778899"
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
		fn := os.Args[2]
		f, err := os.Open(fn)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
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
		fn := os.Args[2]
		f, err := os.Open(fn)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		m, err := torrent.ReadMetaData(f)
		if err != nil {
			fmt.Println(err)
		}
		resp, err := torrent.DiscoverPeers(m)
		if err != nil {
			fmt.Println(err)
		}
		for _, p := range resp.Peers {
			fmt.Printf("%d.%d.%d.%d:%d\n", p.IP[0], p.IP[1], p.IP[2], p.IP[3], p.Port)
		}
	} else if command == "handshake" {
		if len(os.Args) < 4 {
			fmt.Println("invalid number of arguments for handshake")
		}
		fn := os.Args[2]
		socket := os.Args[3]
		f, err := os.Open(fn)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		m, err := torrent.ReadMetaData(f)
		if err != nil {
			fmt.Println(err)
		}
		infohash, err := m.GetInfoHashShort()
		if err != nil {
			fmt.Println(err)
		}

		hs := torrent.NewPeerHandshake(infohash, PEER_ID)
		conn, err := net.Dial("tcp", socket)
		if err != nil {
			fmt.Println(err)
		}
		defer conn.Close()
		_, err = conn.Write(hs.Bytes())
		if err != nil {
			fmt.Println(err)
		}
		respBuf := make([]byte, len(hs.Bytes()))
		_, err = conn.Read(respBuf)
		if err != nil {
			fmt.Println(err)
		}
		phs, err := torrent.PeerHandshakeFromBytes(respBuf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Peer ID: %s\n", hex.EncodeToString([]byte(phs.PeerID)))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
