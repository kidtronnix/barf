<img src="https://raw.githubusercontent.com/smaxwellstewart/barf/master/logo.png" width="300" style="float:right;" />
# barf

Barf is a data wrangling tool for producing data streams from an AWS S3 bucket.

Every line of every file is sent through the stream until it closes.

As a bonus, any files in your S3 bucket that end with `.gz` or `.bz2` will be
automatically decompressed using gzip or bzip2 algorithms respectively.

## cli

### Install

The binaries for mac, windows and linux are hosted on the [releases](https://github.com/smaxwellstewart/barf/releases)
page of the github project.

#### OSX

```sh
curl https://github.com/smaxwellstewart/barf/releases/download/v0.1/barf_osx -o barf
chmod +x barf
# move to some directory on your `$PATH`, example:
sudo mv barf /usr/local/bin
```

#### Linux

```sh
wget https://github.com/smaxwellstewart/barf/releases/download/v0.1/barf_linux
mv barf_linux barf
chmod +x barf
# move to some directory on your `$PATH`, example:
sudo mv barf /usr/local/bin/barf
```

#### Windows

```bash
# Binary available: https://github.com/smaxwellstewart/barf/releases/download/v0.1/barf.exe
# Method not implemented yet.
```

### Usage

Think of this as a recursive `cat` command for all files
in your S3 bucket under a certain prefix.

```sh
$ barf s3://bucket/prefix/
# prepare for the data vom!
```

In the more advanced example below, we set a value for `flow` which controls the
rate of http calls / sec. We also set a `duration` to get a small amount of data
by opening up our stream for a limited time.
Finally we pipe everything to some output file for safe keeping.

```sh
$ barf -flow="10.0" -duration="3s" s3://bucket/prefix/ > output
```


## go library

The other way of using barf is as golang library in your own projects.

### Example

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
