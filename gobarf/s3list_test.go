package barf

import (
	"testing"
	"time"

	"launchpad.net/goamz/s3"

	"github.com/smaxwellstewart/go-resiliency/limiter"
	"github.com/stretchr/testify/assert"
)

// mock s3 bucket is a struct to make it easy
// to mock responses from s3 bucket listing.
type mock3Bucket struct {
	resp *s3.ListResp
	err  error
}

func (b *mock3Bucket) List(prefix, delim, marker string, max int) (*s3.ListResp, error) {
	// will default to return error if set
	if b.err != nil {
		return nil, b.err
	}

	// if response has been set to be truncated will return response one
	// more time but will say it's not truncated on the next iteration.
	// in effect this means that the response contents can be read twice when using
	// mocks3bucket with IsTruncated set to true initially.
	if b.resp.IsTruncated == true {
		_resp := *b.resp
		b.resp.IsTruncated = false

		return &_resp, nil
	}
	return b.resp, nil

}

// NEED TO WRITE TEST
func TestS3List(t *testing.T) {
	assert := assert.New(t)

	s := &S3Lister{
		Bucket: &mock3Bucket{resp: &s3.ListResp{
			Contents:    []s3.Key{s3.Key{Key: "a"}, s3.Key{Key: "b"}},
			IsTruncated: true,
		}},
		max: 2,
	}

	ch := s.list(make(chan struct{}))

	assert.Equal("a", <-ch)
	assert.Equal("b", <-ch)
	assert.Equal("a", <-ch)
	assert.Equal("b", <-ch)
}

//
// // NEED TO WRITE TEST
func TestS3ListNotTruncated(t *testing.T) {
	assert := assert.New(t)

	s := &S3Lister{
		Bucket: &mock3Bucket{resp: &s3.ListResp{
			Contents:    []s3.Key{s3.Key{Key: "a"}, s3.Key{Key: "b"}},
			IsTruncated: false,
		}},
	}

	ch := s.list(make(chan struct{}))

	assert.Equal("a", <-ch)
	assert.Equal("b", <-ch)
}

func TestS3ListRateLimit(t *testing.T) {
	assert := assert.New(t)

	rl := limiter.New(0, 10.0)

	s := &S3Lister{
		Bucket: &mock3Bucket{resp: &s3.ListResp{
			Contents:    []s3.Key{s3.Key{Key: "a"}, s3.Key{Key: "b"}},
			IsTruncated: true,
		}},
		max:     2,
		Limiter: rl.Limiter(), // limit 10 per second
	}

	ch := s.list(make(chan struct{}))

	quit := time.After(150 * time.Millisecond)

	res := []string{}

	go func() {
		for l := range ch {
			res = append(res, l)
		}
	}()

	<-quit

	assert.Equal(2, len(res))
	rl.Close()
}
