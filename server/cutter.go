package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lesnoi-kot/clip-radiot/mpeg"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

// Retrieve firstChunkSize bytes from an audio to inspect it.
const firstChunkSize = 500 * 1024

type probeAudioResult struct {
	fullContentSize int64
	audioInfo       mpeg.AudioInfo
}

func probeAudio(ctx context.Context, url string) (*probeAudioResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=0-%d", firstChunkSize))

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	audioInfo, err := mpeg.InpsectAudio(res.Body)
	if err != nil {
		return nil, err
	}

	contentFullSize, err := parseContentRange(res.Header.Get("Content-Range"))
	if err != nil {
		return nil, err
	}

	return &probeAudioResult{
		fullContentSize: contentFullSize,
		audioInfo:       audioInfo,
	}, nil
}

type cutAudioParams struct {
	ctx       context.Context
	url       string
	audioInfo mpeg.AudioInfo
	fromMs    int64
	toMs      int64
}

func cutAudio(params cutAudioParams) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		params.ctx,
		http.MethodGet,
		params.url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	fromBytes := params.audioInfo.GetTimestampOffset(params.fromMs)
	toBytes := params.audioInfo.GetTimestampOffset(params.toMs)
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", fromBytes, toBytes))

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	blob, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Skip anything before the first occured actual frame.
	var skipTimes uint = 1

	// If a clip requested from the very beginning,
	// make an extra skip of the LAME-header with uncropped audio info.
	if params.fromMs == 0 {
		skipTimes++
	}

	blob, err = mpeg.SkipToSyncMark(mpeg.ShrinkToSyncMark(blob), skipTimes)
	if err != nil {
		return nil, err
	}

	return blob, nil
}
