package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func fatal(err string) {
	if err != "" {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {

	filename := flag.String("f", "", "target filename")
	flagZero := flag.Bool("zero", false, "use zeroes")
	flagRand := flag.Bool("rand", true, "use random noise")

	flag.Parse()

	// zero and rand cannot both be false
	if !*flagRand && !*flagZero {
		fmt.Fprintln(os.Stderr, "eraser: one of -rand or -zero required")
		flag.Usage()
		fatal("")
	}

	// open filename
	if *filename == "" {
		fmt.Fprintln(os.Stderr, "eraser: target filename required")
		flag.Usage()
		fatal("")
	}

	file, err := os.OpenFile(*filename, os.O_WRONLY, 0)
	if err != nil {
		fatal(err.Error())
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fatal(err.Error())
	}
	size := stat.Size()

	var reader io.Reader
	if *flagZero {
		reader = devZero()
	} else if *flagRand {
		reader = devAES()
	}

	spinner := newProgressSpinner(size)
	spinner.draw()
	for spinner.current < size {
		chunk := min(size-spinner.current, 256*1024)
		chunk, err = io.CopyN(file, reader, chunk)
		if err != nil {
			fatal(err.Error())
		}
		spinner.add(chunk)
	}
	spinner.done()

	file.Sync()

}
