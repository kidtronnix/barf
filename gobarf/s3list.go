package barf

import (
	"errors"

	"launchpad.net/goamz/s3"
)

// type s3Spewer struct {
// 	bucket  *s3.Bucket
// 	path    string
// 	marker  string
// 	limiter chan struct{}
// }
//
//
//
//
// func (s *s3Spewer) Spew() chan []byte {
// }
//
// // list is our listing gopher
// func (s *s3Spewer) list(path string) chan string {
//
// }

const MaxFileListings = 1000

type s3bucketLister interface {
	List(prefix, delim, marker string, max int) (*s3.ListResp, error)
}

type S3Lister struct {
	Bucket   s3bucketLister
	Path     string
	max      int
	Limiter  <-chan struct{}
	DoneChan chan struct{}
}

func (s *S3Lister) list(done chan struct{}) chan string {
	output := make(chan string, s.max)
	delim := ""
	// apply sensible default to max if not specified
	if s.max == 0 {
		s.max = MaxFileListings
	}
	// TODO: think could make work as team of concurrent gophers maybe
	go func() {
		defer close(output)
		var marker string
		for {
			// let's make sure we aren't looping too fast
			if s.Limiter != nil {
				<-s.Limiter
			}

			// call S3 API to fetch our listing
			resp, err := s.Bucket.List(s.Path, delim, marker, s.max)
			if err != nil {
				panic(errors.New("Error listing s3 files: " + err.Error()))
			}
			// loop through contents of resp and send on output chan
			for _, key := range resp.Contents {
				// fmt.Println("key:", key.Key) // debug

				select {
				case output <- key.Key:
				case <-done:
					return
				}
			}

			if resp.IsTruncated {
				// response is trunated so we need to keep looping to get rest of files
				// update local marker so we know which file we got to for next call
				marker = resp.Contents[s.max-1].Key
				continue
			} else {
				// response is not truncated so we have all the files
				break
			}
		}
	}()
	return output
}
