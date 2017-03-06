package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/smaxwellstewart/barf/gobarf"
)

var (
	flow     = flag.Float64("flow", barf.FlowDefault, "rate limit of data gathering process")
	duration = flag.Duration("duration", 0, "rate limit of data gathering process")
)

func main() {

	flag.Parse()

	// check we have supplied a src argument
	if len(flag.Args()) < 1 {
		printUsage()
		log.Fatal("Fatal error! Must supply an s3 url argument.")
	}

	// make new barf stream
	b := barf.New(barf.Stream{
		Src:  flag.Args()[0],
		Flow: *flow,
	})
	defer b.Close()
	strm := b.Barf()

	// start reading stream to stdout
	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for l := range strm {
			select {
			case <-done:
				return
			default:
				fmt.Println(string(l))
			}
		}
	}()

	// did we set duration for stream?
	if duration.Nanoseconds() != 0 {
		quit := time.After(*duration)
		<-quit
		close(done) // quit reading early
	}

	// wait until we either forced to stop or end of datastream
	wg.Wait()
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("$ barf s3://bucket/prefix/to/files")
	fmt.Println("Config flags:")
	flag.PrintDefaults()
}
