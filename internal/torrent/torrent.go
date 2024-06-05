package torrent

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	Announce string    `bencode:"announce"`
	Info     MetaInfo  `bencode:"info"`
	InfoHash hash.Hash `bencode:"-"`
}

type Peer struct {
	IP   [4]uint8
	Port uint16
}
type TrackerResponse struct {
	Interval int
	Peers    []Peer
}

func DecodeTrackerResponse(resp string) (*TrackerResponse, error) {
	r, err := bncode.Decode(resp)
	if err != nil {
		return nil, err
	}
	Interval := r.(map[string]interface{})["interval"].(int)
	Peers := make([]Peer, 0)
	peers := r.(map[string]interface{})["peers"].(string)
	ipbuf := make([]byte, 4)
	portbuf := make([]byte, 2)
	reader := strings.NewReader(peers)
	for {
		_, err = reader.Read(ipbuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		_, err = reader.Read(portbuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		Peers = append(Peers, Peer{IP: [4]uint8(ipbuf), Port: binary.BigEndian.Uint16(portbuf)})
	}
	return &TrackerResponse{
		Interval,
		Peers,
	}, nil
}

func (m *Meta) GetInfoHash() string {
	return fmt.Sprintf("%x", m.InfoHash.Sum(nil))
}

func (m *Meta) GetInfoHashShort() (string, error) {
	b, err := hex.DecodeString(m.GetInfoHash())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m *MetaInfo) Hashes() ([]string, error) {
	r := strings.NewReader(m.Pieces)
	b := make([]byte, 20)
	h := make([]string, 0)
	for {
		_, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		h = append(h, fmt.Sprintf("%x\n", b))
	}
	return h, nil
}

func DiscoverPeers(m *Meta) (*TrackerResponse, error) {
	p := url.Values{}
	infoHash, err := m.GetInfoHashShort()
	if err != nil {
		return nil, err
	}
	p.Add("info_hash", infoHash)
	p.Add("peer_id", "12345123451234512345")
	p.Add("port", "6881")
	p.Add("uploaded", "0")
	p.Add("downloaded", "0")
	p.Add("left", fmt.Sprintf("%d", m.Info.Length))
	p.Add("compact", "1")

	url := fmt.Sprintf("%s?%s", m.Announce, p.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r, err := DecodeTrackerResponse(string(body))
	if err != nil {
		return nil, err
	}
	return r, nil
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

func ReadMetaData(fn string) (Meta, error) {
	f, err := os.Open(fn)
	if err != nil {
		return Meta{}, err
	}
	defer f.Close()
	var meta Meta
	if err := bencode.Unmarshal(f, &meta); err != nil {
		return meta, err
	}
	h := sha1.New()
	if err := bencode.Marshal(h, meta.Info); err != nil {
		return meta, err
	}
	meta.InfoHash = h
	return meta, nil
}
