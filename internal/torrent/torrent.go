package torrent

import (
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bncode"
)

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
