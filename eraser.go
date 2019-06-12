// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var version = "unknown"

func main() {

	// data source flags
	flagZero := flag.Bool("zero", false, "use zeroes")
	flagRand := flag.Bool("rand", false, "use pseudorandom noise")

	// add erasure note flag
	flagNote := flag.Bool("note", false, "add timestamped erasure note")

	// print version
	printVersion := func() {
		fmt.Printf("eraser version %s\n", version)
	}

	// custom usage message
	flag.Usage = func() {
		printVersion()
		fmt.Fprintf(os.Stderr, "Usage: $ %s [options] filename\n", os.Args[0])
		flag.PrintDefaults()
	}

	// parse argv
	flag.Parse()

	// zero and rand cannot both be false
	if !*flagRand && !*flagZero {
		fmt.Fprintln(os.Stderr, "err: one of -rand or -zero required")
		flag.Usage()
		os.Exit(1)
	}

	// warn on unused flags
	if len(flag.Args()) > 1 {
		check(fmt.Errorf("err: unused arguments: %s", flag.Args()[1:]))
	}

	// get filename
	filename := flag.Arg(0)
	if filename == "" {
		fmt.Fprintln(os.Stderr, "err: target filename required")
		flag.Usage()
		os.Exit(1)
	}

	// open file for writing
	file, err := os.OpenFile(filename, os.O_WRONLY, 0)
	check(err)

	// get size by seeking
	size, err := file.Seek(0, os.SEEK_END)
	check(err)
	_, err = file.Seek(0, 0)
	check(err)

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

	// optionally add erasure note
	if *flagNote {
		_, err = file.WriteAt([]byte(fmt.Sprintf("ERASURE ON %s\n", time.Now().UTC().Format(time.RFC3339))), 0)
		check(err)
	}

	// finish progress meter
	meter.done()

	// flush to disk
	file.Sync()

	// close the descriptor
	err = file.Close()
	check(err)

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
