package torrent

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bncode"
	"github.com/jackpal/bencode-go"
)

type MetaInfo struct {
	Name        string `bencode:"name"`
	Pieces      string `bencode:"pieces"`
	Length      int64  `bencode:"length"`
	PieceLength int64  `bencode:"piece length"`
}
type Meta struct {
	Announce string   `bencode:"announce"`
	Info     MetaInfo `bencode:"info"`
}

func Read(fn string) (map[string]interface{}, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := make([]byte, 1024)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	d, err := bncode.Decode(string(b))
	if err != nil {
		return nil, err
	}
	return d.(map[string]interface{}), nil
}

func ReadMetaData(fn string) (Meta, string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return Meta{}, "", err
	}
	defer f.Close()
	var meta Meta
	if err := bencode.Unmarshal(f, &meta); err != nil {
		return meta, "", err
	}
	h := sha1.New()
	if err := bencode.Marshal(h, meta.Info); err != nil {
		return meta, "", err
	}
	return meta, fmt.Sprintf("%x", h.Sum(nil)), nil
}
