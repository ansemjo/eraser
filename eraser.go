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

	// open filename in first argument
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
	fmt.Printf("%s is %d bytes\n", file.Name(), size)

	var reader io.Reader
	if *flagZero {
		reader = devZero()
	} else if *flagRand {
		reader = devAES()
	}

	_, err = io.CopyN(file, reader, size)
	if err != nil {
		fatal(err.Error())
	}
	file.Sync()

}
