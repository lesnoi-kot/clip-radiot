package server

import (
	"strconv"
	"strings"
)

func parseContentRange(rangeValue string) (int64, error) {
	_, contentSizeStr, _ := strings.Cut(rangeValue, "/")

	contentSize, err := strconv.ParseInt(contentSizeStr, 10, 64)
	return contentSize, err
}
