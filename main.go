package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/smaxwellstewart/barf/gobarf"
)

var (
	flow     = flag.Float64("flow", barf.DefaultFlow, "rate limit of data gathering process")
	duration = flag.Duration("duration", 0, "duration of stream. useful for taking samples. if zero or not set, stream will be endure until end of file content.")
)

func main() {
	// parse flags from cli
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

	// make chan for communicating when the timer is done
	done := make(chan struct{})

	// start the timer in the background
	go func() {
		// did we set duration for stream?
		if duration.Nanoseconds() != 0 {
			quit := time.After(*duration)
			<-quit
			close(done) // close reading from strm!
		}
	}()

	// loop through results
	for l := range strm {
		select {
		case <-done:
			return
		default:
			fmt.Println(l)
		}
	}

	// fmt.Println(i)

}

func printUsage() {
	fmt.Println("Usage: barf [options] <s3://url>")
	fmt.Println("\nOptions: ")
	flag.PrintDefaults()
}
