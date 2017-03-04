package barf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEmptyPath(t *testing.T) {
	assert := assert.New(t)

	url := ""
	_, _, err := parsePath(url)

	assert.Error(err)
	assert.Equal(errorBadPath, err)
}

func TestParseFSPath(t *testing.T) {
	assert := assert.New(t)

	url := "fs://yoyoyo/asd"
	src, path, err := parsePath(url)

	assert.NoError(err)
	assert.Equal(src, "fs")
	assert.Equal(path, "yoyoyo/asd")
}

func TestParseS3Path(t *testing.T) {
	assert := assert.New(t)

	url := "s3://yoyoyo/asd"
	src, path, err := parsePath(url)

	assert.NoError(err)
	assert.Equal(src, "s3")
	assert.Equal(path, "yoyoyo/asd")
}

func TestParseDefaultPath(t *testing.T) {
	assert := assert.New(t)

	url := "yoyoyo/asd"
	src, path, err := parsePath(url)

	assert.NoError(err)
	assert.Equal(src, "fs")
	assert.Equal(path, "yoyoyo/asd")
}

func TestParseUnssuportedSource(t *testing.T) {
	assert := assert.New(t)

	url := "hdfs://yoyoyo/asd"
	_, _, err := parsePath(url)

	assert.Error(err)
	assert.Equal(errorUnsupportedProtocol, err)

}

func TestParseCrazyPath(t *testing.T) {
	assert := assert.New(t)

	url := "s3://yoyoyo/asds3://foo"
	_, _, err := parsePath(url)

	assert.Error(err)
	assert.Equal(errorBadPath, err)

}
