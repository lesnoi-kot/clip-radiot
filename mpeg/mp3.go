package mpeg

import (
	"errors"
	"io"

	"github.com/tcolgate/mp3"
)

const (
	ID3TagHeaderSize     = 10
	framesCountToAnalyze = 500

	// http://www.mp3-tech.org/programmer/frame_header.html
	// http://www.multiweb.cz/twoinches/mp3inside.htm
	syncMark uint16 = 0b11111111111_11011
)

type AudioInfo struct {
	FramesOffset  int64
	FrameSize     float64 // Average frame size
	FrameDuration float64 // ms
}

func (audioInfo AudioInfo) GetTimestampOffset(ms int64) int64 {
	return audioInfo.FramesOffset + int64(float64(ms)/float64(audioInfo.FrameDuration)*float64(audioInfo.FrameSize))
}

func InpsectAudio(stream io.Reader) (AudioInfo, error) {
	var info AudioInfo
	var frame mp3.Frame
	var skipped int
	var avgFrameSize float64

	info.FramesOffset, _ = ParseTagFullSize(stream)
	io.CopyN(io.Discard, stream, info.FramesOffset-ID3TagHeaderSize)

	decoder := mp3.NewDecoder(stream)

	for i := 1; i <= framesCountToAnalyze; i++ {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			return info, err
		}
		avgFrameSize += float64(frame.Size())
	}

	avgFrameSize = avgFrameSize / framesCountToAnalyze

	info.FramesOffset += int64(skipped)
	info.FrameSize = avgFrameSize
	info.FrameDuration = (1000 / float64(frame.Header().SampleRate())) * float64(frame.Samples())

	return info, nil
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
