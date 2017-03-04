package barf

import (
	"errors"
	"strings"
)

var (
	errorBadPath             = errors.New("Unable to parse provided path.")
	errorUnsupportedProtocol = errors.New("Unsupported filesystem.")
)

func parsePath(srcPath string) (string, string, error) {
	// first check if empty
	if srcPath == "" {
		return "", "", errorBadPath
	}
	parts := strings.Split(srcPath, "://")
	if len(parts) == 2 {
		// ok, we have a supplied protocol!
		// now need to check if protocol we support.
		switch parts[0] {
		case "s3":
			return "s3", parts[1], nil
		case "fs":
			return "fs", parts[1], nil
		default:
			return "", "", errorUnsupportedProtocol
		}
	} else if len(parts) < 2 {
		// so no protocol specified, so assume filesystem
		return "fs", parts[0], nil
	} else {
		return "", "", errorBadPath
	}

}
