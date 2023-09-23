package mpeg_test

import (
	"bytes"
	"testing"

	"github.com/lesnoi-kot/clip-radiot/mpeg"
)

func TestRewindToSyncMark(t *testing.T) {
	t.Run("Frame at the very start", func(t *testing.T) {
		index, err := mpeg.SkipToSyncMark([]byte{
			0xFF, 0xFB, 0x90, 0xC4, 0, 0, 0x0B,
		}, 1)

		if err != nil {
			t.Errorf("Got error: %s", err)
		}

		if index[0] != 0xFF || index[1] != 0xFB || len(index) != 7 {
			t.Error("Got invalid start of a frame")
		}
	})

	t.Run("Some padding before frame", func(t *testing.T) {
		index, err := mpeg.SkipToSyncMark([]byte{
			0xDD, 0, 0xFF, 0, 0, 0, 0, 0xFF, 0xFB, 0x90,
		}, 1)

		if err != nil {
			t.Errorf("Got error: %s", err)
		}

		if index[0] != 0xFF || index[1] != 0xFB || len(index) != 3 {
			t.Error("Got invalid start of a frame")
		}
	})

	t.Run("No sync mark", func(t *testing.T) {
		_, err := mpeg.SkipToSyncMark([]byte{
			0xDD, 0, 0xFF, 0, 1, 2, 3, 0xFF, 0xFF, 0xFF,
		}, 1)

		if err == nil {
			t.Errorf("Expected error but got nil")
		}
	})

	t.Run("No sync mark", func(t *testing.T) {
		_, err := mpeg.SkipToSyncMark([]byte{
			0xDD, 0, 0xFF, 0, 1, 2, 3, 0xFB, 0xFF, 0x90,
		}, 1)

		if err == nil {
			t.Errorf("Expected error but got nil")
		}
	})
}

func TestShrinkToSyncMark(t *testing.T) {
	t.Run("Should cut last frame", func(t *testing.T) {
		index := mpeg.ShrinkToSyncMark([]byte{
			0xFF, 0xFB, 1, 2, 3, 4, 5,
			0xFF, 0xFB, 0, 0, 0, 0,
		})

		if !bytes.Equal(index, []byte{0xFF, 0xFB, 1, 2, 3, 4, 5}) {
			t.Error("Got invalid start of a frame")
		}
	})

	t.Run("Multiple frames", func(t *testing.T) {
		index := mpeg.ShrinkToSyncMark([]byte{
			0xFF, 0xFB, 1, 2, 3, 4, 5,
			0xFF, 0xFB, 0, 0, 0, 0, 0,
			0xFF, 0xFB, 6, 6, 6,
		})

		if !bytes.Equal(index, []byte{0xFF, 0xFB, 1, 2, 3, 4, 5, 0xFF, 0xFB, 0, 0, 0, 0, 0}) {
			t.Error("Got invalid start of a frame")
		}
	})
}
