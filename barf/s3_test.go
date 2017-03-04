package barf

import (
	"testing"

	"launchpad.net/goamz/s3"

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
	} else {
		return b.resp, nil
	}
}


// NEED TO WRITE TEST
func TestS3List(t *testing.T) {
	assert := assert.New(t)

	s := *s3Lister{
    &mock3Bucket{&s3.ListResp{

    }}
  }
}
