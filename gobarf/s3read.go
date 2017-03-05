package barf

import (
	"bufio"
	"compress/gzip"
	"io"
	"strings"
)

type s3bucketReader interface {
	GetReader(key string) (io.ReadCloser, error)
}

type S3Reader struct {
	Bucket  s3bucketReader
	Limiter <-chan struct{}
}

func (s *S3Reader) read(done chan struct{}, in chan string) chan []byte {
	out := make(chan []byte)
	go func() {
		defer close(out)
		for key := range in {
			// let's make sure we aren't looping too fast
			if s.Limiter != nil {
				<-s.Limiter
			}
			// fetch s3 object
			rc, err := s.Bucket.GetReader(key)
			if err != nil {
				panic("error fetching s3 object: " + err.Error())
			}

			comp := parseCompression(key)

			// decompress reader depending on setting
			switch comp {
			case "":
				// no compression, do nothing
			case "gzip":
				rc, err = gzip.NewReader(rc)
				if err != nil {
					panic("error ungzipping s3 object: " + err.Error())
				}
				defer rc.Close()
			default:
				panic("unsupported compression")
			}

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				select {
				case out <- scanner.Bytes():
				case <-done:
					return
				}

			}
			if err := scanner.Err(); err != nil {
				panic("error scanning through file: " + err.Error())
			}
		}
	}()

	return out
}

func parseCompression(key string) string {
	comp := ""
	// does key end in .gz?
	i := strings.Index(key, ".gz")
	if i == len(key)-len(".gz") {
		comp = "gzip"
	}
	return comp
}
