package torrent

import (
	"fmt"
	"os"
	"testing"
)

func TestTorrent(t *testing.T) {
	t.Run("torrent - read metadata", func(t *testing.T) {
		f, err := os.Open("../../sample.torrent")
		if err != nil {
			t.Errorf("didn't expect an error, got one: %v", err)
		}
		defer f.Close()
		m, err := ReadMetaData(f)
		if err != nil {
			t.Errorf("didn't expect an error, got one: %v", err)
		}
		got := fmt.Sprintf("%x", m.InfoHash.Sum(nil))
		want := "d69f91e6b2ae4c542468d1073a71d4ea13879a7f"
		// fmt.Sprintf("%x", h.Sum(nil))
		if want != got {
			t.Errorf("got %s, want %s", want, got)
		}
	})
}
