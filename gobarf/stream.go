package barf

import "launchpad.net/goamz/aws"

type Stream struct {
	Src    string
	Flow   float64    // leaky bucket drip rate!
	Region aws.Region // only needed for s3 streams
}
