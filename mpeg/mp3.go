package mpeg

import "errors"

const (
	ID3TagHeaderSize = 10

	// Hardcoded params for the Radio-T audio files.
	sampleRate         = 44100
	bitrate            = 128000
	padding            = 1
	frameSize          = (144*bitrate)/sampleRate + padding
	msPerFrame float64 = 1000 / ((float64(bitrate) / 8) / frameSize)

	// http://www.mp3-tech.org/programmer/frame_header.html
	// http://www.multiweb.cz/twoinches/mp3inside.htm
	syncMark uint16 = 0b11111111111_11011
)

func TimestampToByteOffset(ms int64) int64 {
	return int64(float64(ms) / msPerFrame * float64(frameSize))
}

func SkipToSyncMark(body []byte, n uint) ([]byte, error) {
	if n == 0 {
		return body, nil
	}

	for i := 0; i < len(body)-1; i++ {
		if ((uint16(body[i]) << 8) | uint16(body[i+1])) == syncMark {
			if n == 1 {
				return body[i:], nil
			}

			n--
		}
	}

	return body, errors.New("No sync mark found")
}

func ShrinkToSyncMark(body []byte) []byte {
	i := len(body) - 2

	for ; i >= 0; i-- {
		if ((uint16(body[i]) << 8) | uint16(body[i+1])) == syncMark {
			return body[:i]
		}
	}

	return body
}
