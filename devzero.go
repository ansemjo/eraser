// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"io"
)

// zero is a simple struct that only returns zeroes
// keeping a buffer to avoid permanent reallocation
type zero struct {
	buf []byte
}

func devZero() (r io.Reader) {
	return &zero{}
}

func (z *zero) Read(dst []byte) (n int, err error) {
	dlen := len(dst)
	if dlen > len(z.buf) {
		z.buf = make([]byte, dlen)
	}
	dst = z.buf[:dlen]
	for i := range dst {
		dst[i] = 0
	}
	return dlen, nil
}
