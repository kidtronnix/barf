package barf

import (
	"runtime"
	"sync"

	"github.com/smaxwellstewart/go-resiliency/limiter"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

var (
	FlowDefault     = 1.0
	DefaultS3Region = aws.USEast
)

type Barfer interface {
	Barf() <-chan []byte
	Close()
}

//

func New(strm Stream) Barfer {
	src, path, err := parseSrc(strm.Src)
	if err != nil {
		panic(err)
	}
	var barfer Barfer
	switch src {
	case "s3":
		auth, err := aws.EnvAuth()
		if err != nil {
			panic(err.Error())
		}

		bucketname, prefix, err := parseBucket(path)
		if err != nil {
			panic(err)
		}
		// trust me you want a limit!
		if strm.Flow == 0.0 {
			strm.Flow = FlowDefault
		}

		// create a rate limiter for ourselves
		rl := limiter.New(0, strm.Flow)

		// set default region to
		if strm.Region.Name == "" {
			strm.Region = DefaultS3Region
		}

		// connect to our s3 bucket
		s := s3.New(auth, strm.Region)
		bucket := s.Bucket(bucketname)

		// configure our barfer
		barfer = &S3Barfer{
			Lister: &S3Lister{
				Bucket:  bucket,
				Limiter: rl.Limiter(),
				Path:    prefix,
			},
			Reader: &S3Reader{
				Bucket:  bucket,
				Limiter: rl.Limiter(),
			},
			Readers:  runtime.NumCPU() - 1,
			doneChan: make(chan struct{}),
		}
	default:
		panic("Unsupported storage protocol!")
	}
	return barfer
}

type S3Barfer struct {
	Lister   *S3Lister
	Reader   *S3Reader
	Readers  int
	doneChan chan struct{}
}

func (b *S3Barfer) Barf() <-chan []byte {
	// start listing gopher!
	listings := b.Lister.list(b.doneChan)

	// start reading gophers
	chans := make([]<-chan []byte, b.Readers)
	for i := 0; i < b.Readers; i++ {
		chans[i] = b.Reader.read(b.doneChan, listings)
	}

	return merge(b.doneChan, chans...)
}

func (b *S3Barfer) Close() {
	close(b.doneChan)
}

func merge(done chan struct{}, chs ...<-chan []byte) chan []byte {
	var wg sync.WaitGroup
	out := make(chan []byte)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan []byte) {
		defer wg.Done()

		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(chs))
	for _, c := range chs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// TODO: refactor s3 specific barfer creation logic
// func newS3Barfer(bucket string, prefix string) {
//
// }
