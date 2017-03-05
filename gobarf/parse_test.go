package barf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSrcEmptyPath(t *testing.T) {
	assert := assert.New(t)

	url := ""
	_, _, err := parseSrc(url)

	assert.Error(err)
	assert.Equal(errorBadPath, err)
}

func TestParseSrcFSPath(t *testing.T) {
	assert := assert.New(t)

	url := "fs://yoyoyo/asd"
	src, path, err := parseSrc(url)

	assert.NoError(err)
	assert.Equal(src, "fs")
	assert.Equal(path, "yoyoyo/asd")
}

func TestParseSrcS3Path(t *testing.T) {
	assert := assert.New(t)

	url := "s3://yoyoyo/asd"
	src, path, err := parseSrc(url)

	assert.NoError(err)
	assert.Equal(src, "s3")
	assert.Equal(path, "yoyoyo/asd")
}

func TestParseSrcDefaultPath(t *testing.T) {
	assert := assert.New(t)

	url := "yoyoyo/asd"
	src, path, err := parseSrc(url)

	assert.NoError(err)
	assert.Equal(src, "fs")
	assert.Equal(path, "yoyoyo/asd")
}

func TestParseSrcUnssuportedSource(t *testing.T) {
	assert := assert.New(t)

	url := "hdfs://yoyoyo/asd"
	_, _, err := parseSrc(url)

	assert.Error(err)
	assert.Equal(errorUnsupportedProtocol, err)

}

func TestParseSrcCrazyPath(t *testing.T) {
	assert := assert.New(t)

	url := "s3://yoyoyo/asds3://foo"
	_, _, err := parseSrc(url)

	assert.Error(err)
	assert.Equal(errorBadPath, err)

}

func TestParseBucket(t *testing.T) {
	assert := assert.New(t)

	path := "yoyoyo"
	bucket, prefix, err := parseBucket(path)

	assert.NoError(err)
	assert.Equal(path, bucket)
	assert.Equal(prefix, "")

}

func TestParseBucketPrefix(t *testing.T) {
	assert := assert.New(t)

	path := "mybucket/x"
	bucket, prefix, err := parseBucket(path)

	assert.NoError(err)
	assert.Equal("mybucket", bucket)
	assert.Equal(prefix, "x")

}

func TestParseBucketError(t *testing.T) {
	assert := assert.New(t)

	path := ""
	_, _, err := parseBucket(path)

	assert.Error(err)
	assert.Equal(errorNoBucket, err)

	path = "/"
	_, _, err = parseBucket(path)

	assert.Error(err)
	assert.Equal(errorNoBucket, err)
}
