package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/clip-radiot/mpeg"
)

const (
	maxDuration = 1 * 60 * 1000 // 1 min
	minDuration = 5_000         // 5 sec
)

func cutAudioHandler(c echo.Context) error {
	episode, err := strconv.ParseUint(c.QueryParam("episode"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Некорректный номер эпизода")
	}

	fromMs, err := strconv.ParseInt(c.QueryParam("from"), 10, 64)
	if err != nil || fromMs < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Некорректно задан отрезок времени")
	}

	toMs, err := strconv.ParseInt(c.QueryParam("to"), 10, 64)
	if err != nil || toMs < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Некорректно задан отрезок времени")
	}

	if toMs <= fromMs {
		return echo.NewHTTPError(http.StatusBadRequest, "Некорректно задан отрезок времени")
	}

	duration := toMs - fromMs
	if duration > maxDuration || duration < minDuration {
		return echo.NewHTTPError(http.StatusBadRequest, "Длина клипа не должна быть меньше 5 секунд или больше 1-ой минуты")
	}

	audioURL := fmt.Sprintf("https://cdn.radio-t.com/rt_podcast%d.mp3", episode)

	log.Printf(
		"Requested cut for Radio-T #%d: from=%dms, to=%dms, url=%s",
		episode, fromMs, toMs, audioURL,
	)

	ctx := c.Request().Context()

	audioInfo, err := probeAudio(ctx, audioURL)
	if err != nil {
		return fmt.Errorf("Audio probing error: %w", err)
	}

	log.Printf(
		"Audio probe info for Radio-T #%d: size=%d bytes, tag=%d bytes",
		episode,
		audioInfo.fullContentSize,
		audioInfo.fullTagSize,
	)

	audioFrames, err := cutAudio(cutAudioParams{
		ctx:         ctx,
		url:         audioURL,
		offsetBytes: audioInfo.fullTagSize,
		fromMs:      fromMs,
		toMs:        toMs,
	})
	if err != nil {
		return fmt.Errorf("Audio cut error: %w", err)
	}

	log.Printf("Successfully extracted audio data for Radio-T #%d: size=%d bytes", episode, len(audioFrames))

	buff := bytes.NewBuffer(nil)

	// Write the tag.
	if err = mpeg.WriteRadioTArtistInfo(buff, int(episode)); err != nil {
		return err
	}

	// Write mp3 audio data.
	if _, err = buff.Write(audioFrames); err != nil {
		return err
	}

	c.Blob(http.StatusOK, "audio/mpeg", buff.Bytes())
	return nil
}
