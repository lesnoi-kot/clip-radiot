package mpeg

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/bogem/id3v2"
)

func ParseTagFullSize(stream io.Reader) (int64, error) {
	// The first part of the ID3v2 tag is the 10 byte tag header.
	var buff [ID3TagHeaderSize]byte
	n, err := stream.Read(buff[:])

	if err != nil && err != io.EOF {
		return 0, err
	}

	if n != ID3TagHeaderSize {
		return 0, errors.New("Incomplete id3 header data")
	}

	// The first three bytes of the tag are always "ID3".
	if !bytes.Equal(buff[:3], []byte("ID3")) {
		return 0, errors.New("Input data is not an idv3 header")
	}

	size, err := parseSize(buff[6:])
	if err != nil {
		return 0, err
	}

	return size + ID3TagHeaderSize, nil
}

// "Inspired" by https://github.com/n10v/id3v2/blob/v1.2.0/size.go#L80
// Unfortunately it's private function in the lib.
func parseSize(data []byte) (int64, error) {
	if len(data) != 4 {
		return 0, fmt.Errorf("ID3 size invalid size: %d, should be 4", len(data))
	}

	var size int64

	for _, b := range data {
		if b&128 > 0 { // 128 = 0b1000_0000
			return 0, errors.New("ID3 invalid size format")
		}

		size = (size << 7) | int64(b)
	}

	return size, nil
}

func WriteRadioTArtistInfo(w io.Writer, episodeNumber int) error {
	tag := id3v2.NewEmptyTag()
	tag.SetArtist("Umputun, Bobuk, Gray, Ksenks, Alek.sys")
	tag.SetTitle(fmt.Sprintf("Радио-Т %d", episodeNumber))

	_, err := tag.WriteTo(w)
	return err
}
