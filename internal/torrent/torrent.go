package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bncode"
	"github.com/jackpal/bencode-go"
)

type ProtocolString string

const (
	BitTorrentProtocolString ProtocolString = "BitTorrent protocol"
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

type PeerHandshake struct {
	ProtocolString ProtocolString
	reservedBytes  string
	InfoHash       string
	PeerID         string
	position       int
}

func (p *PeerHandshake) Bytes() []byte {
	l := []byte{uint8(len(p.ProtocolString))}
	hs := []byte(string(p.ProtocolString) + p.reservedBytes + p.InfoHash + p.PeerID)
	return bytes.Join([][]byte{l, hs}, []byte(""))
}

func PeerHandshakeFromBytes(b []byte) (*PeerHandshake, error) {
	protocolStringLen := b[0]
	i := 1
	ProtocolString := ProtocolString(string(b[:protocolStringLen]))
	i += len(BitTorrentProtocolString)
	reservedBytes := string(b[i : i+8])
	i += 8
	InfoHash := string(b[i : i+20])
	i += 20
	PeerID := string(b[i : i+20])
	position := 0
	return &PeerHandshake{
		ProtocolString,
		reservedBytes,
		InfoHash,
		PeerID,
		position,
	}, nil
}

func NewPeerHandshake(h string, pid string) *PeerHandshake {
	return &PeerHandshake{
		ProtocolString: BitTorrentProtocolString,
		reservedBytes:  "00000000",
		InfoHash:       h,
		PeerID:         pid,
		position:       0,
	}
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

func ReadMetaData(r io.Reader) (*Meta, error) {
	var meta Meta
	if err := bencode.Unmarshal(r, &meta); err != nil {
		return nil, err
	}
	h := sha1.New()
	if err := bencode.Marshal(h, meta.Info); err != nil {
		return nil, err
	}
	meta.InfoHash = h
	return &meta, nil
}

func SendHandShake(m *Meta, p *Peer) {
}
