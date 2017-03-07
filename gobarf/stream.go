package barf

import "launchpad.net/goamz/aws"

type Stream struct {
	Src    string     // Src is the source of your file on S3. Example, s3://bucket/prefix/
	Flow   float64    // Flow is the rate limit of http calls per second.
	Region aws.Region // Region of AWS bucket.
}
