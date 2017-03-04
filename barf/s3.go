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

type s3bucketLister interface {
	List(prefix, delim, marker string, max int) (*s3.ListResp, error)
}

type s3lister struct {
	bucket    s3bucketLister
	path      string
	errorMode int
	limiter   chan struct{}
}

func (s *s3lister) list() chan string {
	output := make(chan string)
	max := 1000
	delim := ""
	// TODO: think could make work as team of concurrent gophers maybe
	go func() {

		var marker string
		for {
			// let's make sure we aren't looping too fast
			if s.limiter != nil {
				<-s.limiter
			}

			// call S3 API to fetch our listing
			resp, err := s.bucket.List(s.path, delim, marker, max)
			if err != nil {
				panic(errors.New("Error listing s3 files: " + err.Error()))
			}

			// loop through contents of resp and send on output chan
			for _, key := range resp.Contents {
				output <- key.Key
			}

			if resp.IsTruncated {
				// response is trunated so we need to keep looping to get rest of files
				// update local marker so we know which file we got to for next call
				marker = resp.Contents[max-1].Key
				continue
			} else {
				// response is not truncated so we have all the files
				break
			}
		}
		//
		close(output)
	}()
	return output
}
