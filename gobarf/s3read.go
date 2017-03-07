package barf

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
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

			var r io.Reader
			// decompress reader depending on setting
			switch parseCompression(key) {
			case "":
				r = rc
			case "bz2":
				r = bzip2.NewReader(rc)
				if err != nil {
					panic("error in bzip2 decompression: " + err.Error())
				}
			case "gz":
				r, err = gzip.NewReader(rc)
				if err != nil {
					panic("error in gzip decompression: " + err.Error())
				}
			default:
				panic("unsupported compression")
			}

			scanner := bufio.NewScanner(r)
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

			rc.Close()
		}
	}()

	return out
}
