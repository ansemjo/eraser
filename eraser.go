package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {

	// data source flags
	flagZero := flag.Bool("zero", false, "use zeroes")
	flagRand := flag.Bool("rand", true, "use random noise")

	// custom usage message
	flag.Usage = func() {
		fmt.Println("Usage: $ eraser [options] filename")
		flag.PrintDefaults()
	}

	// parse argv
	flag.Parse()

	// zero and rand cannot both be false
	if !*flagRand && !*flagZero {
		fmt.Fprintln(os.Stderr, "eraser: one of -rand or -zero required")
		flag.Usage()
		os.Exit(1)
	}

	// warn on unused flags
	if len(flag.Args()) > 1 {
		fmt.Fprintf(os.Stderr, "Warning: unused arguments: %s\n", flag.Args()[1:])
	}

	// get filename
	filename := flag.Arg(0)
	if filename == "" {
		fmt.Fprintln(os.Stderr, "eraser: target filename required")
		flag.Usage()
		os.Exit(1)
	}

	// open file for writing
	file, err := os.OpenFile(filename, os.O_WRONLY, 0)
	check(err)
	defer file.Close()

	// get size
	stat, err := file.Stat()
	check(err)
	size := stat.Size()

	// initialize desired reader
	var reader io.Reader
	if *flagZero {
		reader = devZero()
	} else if *flagRand {
		reader = devAES()
	}

	// draw progress meter and copy chunks of data
	meter := newProgress(size)
	meter.draw()
	for meter.current < size {
		chunk := min(size-meter.current, 256*1024)
		chunk, err = io.CopyN(file, reader, chunk)
		check(err)
		meter.add(chunk)
	}
	meter.done()

	// flush to disk
	file.Sync()

}

// check for fatal errors
func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// return smaller number
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
