package main

// zero is a simple struct that only returns zeroes
// keeping a buffer to avoid permanent reallocation
type zero struct {
	buf []byte
}

func devZero() (z *zero) {
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
