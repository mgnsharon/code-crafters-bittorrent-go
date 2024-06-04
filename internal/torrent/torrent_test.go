package torrent

import (
	"encoding/json"
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
				"pieces": `�v�z*����k\x13g&�\x0f��\x03\x02-n"u�\x04�vfVsn��\x10�R\x04��5�
				�z\x02\x13�\x19���	r'�����\x17ee`,
			},
		}
		jsonOutput, _ := json.Marshal(got)
		expectedOutput, _ := json.Marshal(want)
		if string(jsonOutput) != string(expectedOutput) {
			t.Errorf("got %s, want %s", jsonOutput, expectedOutput)
		}
	})
}
