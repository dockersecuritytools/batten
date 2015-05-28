# Procfs

Procfs is a parser for the /proc virtual filesystem on Linux written in the Go programming lanugage.

Not all /proc interfaces are currently supported; but pull requests are welcome!

## Installation

go get github.com/jandre/procfs 

## Examples
See the `*_test` files for usage, but in short:

```go
import (
  "github.com/jandre/procfs"
)

// fetch all processes from /proc
// returns a map of pid -> *Process 
processes, err := procfs.Processes();

```

## Documentation

Documentation can be found at: http://godoc.org/github.com/jandre/procfs
