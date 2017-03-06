# barf

Barf is a tool for producing data streams from files on S3.

Any files in your S3 bucket that end with `.gz` will be automatically gzip decompressed.

Each line from every file is sent through the stream.

## cli usage

`barf` can be used a command line tool that will write all lines of data to stdout.

### basic

```sh
$ barf s3://myawsbucket/prefix/to/my_files
```

This will print the contents of all the files found in your `myawsbucket`
with the prefix `prefix/to/my_files` to stdout.

### advanced

```sh
$ barf s3://myawsbucket/prefix/to/my_files -flow="1.0" -duration="3s" > output
```

In this example we set a value for `flow` which controls the rate of http calls / sec.
we also set a `duration` to get a small amount of data.


## golang library

It is possible to use the underlying golang library in your own projects.

### example

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/seedboxtech/barf/gobarf"
)


func main() {
	// setup new stream of data
	b := barf.New(barf.Stream{
		Src:      "s3://myawsbucket/prefix/to/my_files",
	})
	defer b.Close()

	// read stream and print stdout
  	for line := range b.Barf() {
		// do something with the line of content
	}
}
```
