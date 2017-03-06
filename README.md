# barf

Barf is a tool for producing a data streams from S3.

Any files in your s3 bucket that end with `.gz` will be automatically gzip decompressed.

## cli usage

barf can be used a command line tool that will write all data to stdout.

### basic

```sh
$ barf s3://myawsbucket/prefix/to/my_files
```

this will print the contents of all the files found in your `myawsbucket`
with the prefix `prefix/to/my_files` to stdout.

### advanced

```sh
$ barf s3://myawsbucket/prefix/to/my_files -flow="1.0" -duration="3s" > output
```

in this example we set a value for `flow` which controls the rate of http calls / sec.
we also set a `duration` to get a small amount of data.


## golang library

It is possible to use the underlying golang library in your own projects.

### example:

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
  	for l := range b.Barf() {
		fmt.Println(string(l))
	}
}
```
