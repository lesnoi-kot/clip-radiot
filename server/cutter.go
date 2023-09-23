package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/lesnoi-kot/clip-radiot/mpeg"
)

type probeAudioResult struct {
	fullContentSize int64
	fullTagSize     int64
}

func probeAudio(ctx context.Context, url string) (*probeAudioResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Request only tag header info.
	req.Header.Add("Range", fmt.Sprintf("bytes=0-%d", mpeg.ID3TagHeaderSize-1))

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	tagFullSize, err := mpeg.ParseTagFullSize(res.Body)
	if err != nil {
		return nil, err
	}

	contentFullSize, err := parseContentRange(res.Header.Get("Content-Range"))
	if err != nil {
		return nil, err
	}

	return &probeAudioResult{
		fullContentSize: contentFullSize,
		fullTagSize:     tagFullSize,
	}, nil
}

type cutAudioParams struct {
	ctx         context.Context
	url         string
	offsetBytes int64
	fromMs      int64
	toMs        int64
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

	fromBytes := params.offsetBytes + mpeg.TimestampToByteOffset(params.fromMs)
	toBytes := params.offsetBytes + mpeg.TimestampToByteOffset(params.toMs)
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

	// Go to LAME Info-frame first, then skip to actual frames.
	blob, err = mpeg.SkipToSyncMark(blob, 2)
	if err != nil {
		return nil, err
	}

	return blob, nil
}
