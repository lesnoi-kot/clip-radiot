package server

import (
	"errors"
	"strconv"
	"strings"
)

// Parse total size from the Content-Range header value.
// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Range
func parseContentRange(rangeValue string) (int64, error) {
	_, totalSize, found := strings.Cut(rangeValue, "/")

	if !found || totalSize == "*" {
		return 0, errors.New("Total size of the audio is unknown")
	}

	return strconv.ParseInt(totalSize, 10, 64)
}
