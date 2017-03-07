package barf

import (
	"runtime"
	"sync"

	"github.com/smaxwellstewart/go-resiliency/limiter"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

var (
	// DefaultFlow is the default rate limit at which we can make AWS S3 API calls.
	// It is set to a conservative value for most modern hardware.
	DefaultFlow = 1.0
	// DefaultS3Region is the default AWS region for making AWS S3 API calls.
	// For more info, http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
	DefaultS3Region = aws.USEast
	// DefaultWokers is the default number of concurrent wokers for reading s3 files.
	// By default it is set to the num of available cpus - 1
	DefaultWorkers = runtime.NumCPU() - 1
)

type Barfer interface {
	Barf() <-chan string
	Close()
}

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

		// apply defaults...

		// set default region to
		if strm.Region.Name == "" {
			strm.Region = DefaultS3Region
		}

		// trust me you want a rate limit!
		if strm.Flow == 0.0 {
			strm.Flow = DefaultFlow
		}
		rl := limiter.New(1, strm.Flow)

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
			Readers:  DefaultWorkers,
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

func (b *S3Barfer) Barf() <-chan string {
	// start listing gopher!
	listings := b.Lister.list(b.doneChan)

	// start reading gophers
	chans := make([]<-chan string, b.Readers)
	for i := 0; i < b.Readers; i++ {
		chans[i] = b.Reader.read(b.doneChan, listings)
	}

	return merge(b.doneChan, chans...)
}

func (b *S3Barfer) Close() {
	close(b.doneChan)
}

func merge(done chan struct{}, chs ...<-chan string) chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
		defer wg.Done()
		for n := range c {
			// if bytes.Contains(n, []byte("7f190fa3-b3b6-40eb-8696-099c400f55c6")) {
			// 	fmt.Printf("[merger] found imp: %v %s\n", len(n), n)
			// }
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
