package barf

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/smaxwellstewart/go-resiliency/limiter"
	"github.com/stretchr/testify/assert"
)

type mock3BucketReader struct {
	rc  io.ReadCloser
	err error
}

func (b *mock3BucketReader) GetReader(key string) (io.ReadCloser, error) {
	if b.err != nil {
		return nil, b.err
	}

	switch key {
	case "1":
		// get file 1
		b.rc = ioutil.NopCloser(bytes.NewReader([]byte("a\nb\nc\nd")))
	case "2":
		// get file 2
		b.rc = ioutil.NopCloser(bytes.NewReader([]byte("e\nf\ng\nh")))
	default:
		return nil, errors.New("No such file!")
	}

	return b.rc, nil
}

func TestS3Read(t *testing.T) {
	assert := assert.New(t)

	b := []byte("a\nb")
	rc := ioutil.NopCloser(bytes.NewReader(b))

	r := &S3Reader{
		Bucket: &mock3BucketReader{
			rc: rc,
		},
	}

	in := make(chan string, 1)
	in <- "1"
	close(in)

	ch := r.read(make(chan struct{}), in)

	assert.Equal("a", string(<-ch))
	assert.Equal("b", string(<-ch))
}

func TestS3ReadRateLimit(t *testing.T) {
	assert := assert.New(t)

	rl := limiter.New(0, 10.0)

	b := []byte("a\nb\nc\nd")
	rc := ioutil.NopCloser(bytes.NewReader(b))

	r := &S3Reader{
		Bucket: &mock3BucketReader{
			rc: rc,
		},
		Limiter: rl.Limiter(),
	}

	in := make(chan string, 2)
	in <- "1"
	in <- "2"
	close(in)

	ch := r.read(make(chan struct{}), in)

	quit := time.After(150 * time.Millisecond)

	res := []string{}

	go func() {
		for l := range ch {
			res = append(res, string(l))
		}
	}()

	<-quit

	assert.Equal(4, len(res))
	assert.Equal("a", res[0])
	assert.Equal("b", res[1])
	assert.Equal("c", res[2])
	assert.Equal("d", res[3])

	rl.Close()
}
