package barf

import (
	"errors"
	"strings"
)

var (
	errorBadPath             = errors.New("Unable to parse provided path.")
	errorUnsupportedProtocol = errors.New("Unsupported filesystem.")
	errorNoBucket            = errors.New("No bucket specified.")
)

func parseSrc(srcPath string) (string, string, error) {
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

func parseBucket(bucketPath string) (string, string, error) {
	if bucketPath == "" {
		return "", "", errorNoBucket
	}
	parts := strings.Split(bucketPath, "/")
	if parts[0] == "" {
		return "", "", errorNoBucket
	}
	// return prefix if we have it
	if len(parts) > 1 {
		return parts[0], strings.Join(parts[1:], "/"), nil
	}

	return parts[0], "", nil

}

// parses compression from s3 file keys
func parseCompression(key string) string {
	possible := []string{".gz", ".bz2"}
	for _, ending := range possible {
		// does key end in one of possible endings
		i := strings.Index(key, ending)
		if i == len(key)-len(ending) {
			return strings.TrimPrefix(ending, ".") // return file ending but without "."
		}
	}
	// no compression detected
	return ""
}
