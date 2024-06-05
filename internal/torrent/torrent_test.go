package torrent

import (
	"testing"
)

func TestTorrent(t *testing.T) {
	t.Run("torrent", func(t *testing.T) {
		got, err := Read("../../sample.torrent")
		if err != nil {
			t.Errorf("didn't expect an error, got one: %v", err)
		}
		want := map[string]interface{}{
			"announce":   "http://bittorrent-test-tracker.codecrafters.io/announce",
			"created by": "mktorrent 1.1",
			"info": map[string]interface{}{
				"length":       92063,
				"name":         "sample.txt",
				"piece length": 32768,
				"pieces":       "\ufffdv\ufffdz*\ufffd\ufffd\ufffd\ufffdk\u0013g\u0026\ufffd\u000f\ufffd\ufffd\u0003\u0002-n\"u\ufffd\u0004\ufffdvfVsn\ufffd\ufffd\u0010\ufffdR\u0004\ufffd\ufffd5\ufffd\r\ufffdz\u0002\u0013\ufffd\u0019\ufffd\ufffd\ufffd\tr'\ufffd\ufffd\ufffd\ufffd\ufffd\u0017",
			},
		}

		if want["announce"] != got["announce"] {
			t.Errorf("got %s, want %s", want["announce"], got["announce"])
		}
	})

	t.Run("torrent - read meta data", func(t *testing.T) {
		_, got, err := ReadMetaData("../../sample.torrent")
		if err != nil {
			t.Errorf("didn't expect an error, got one: %v", err)
		}
		want := "d69f91e6b2ae4c542468d1073a71d4ea13879a7f"

		if want != got {
			t.Errorf("got %s, want %s", want, got)
		}
	})
}
