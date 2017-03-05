package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/seedboxtech/barf/gobarf"
)

var (
	src      = flag.String("src", "", "required source url of data. only supported is s3 url.")
	flow     = flag.Float64("flow", barf.FlowDefault, "rate limit of data gathering process")
	duration = flag.Duration("duration", 0, "rate limit of data gathering process")
)

func main() {

	flag.Parse()

	// check if src provided!
	if *src == "" {
		fmt.Println("Usage: ")
		flag.PrintDefaults()
		log.Fatal("Must provide a src flag!")
	}

	// make new barf stream
	b := barf.New(barf.Stream{
		Src:  *src,
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

	// give small amount of time for println to complete
	// <-time.After(1000 * time.Millisecond)
}
