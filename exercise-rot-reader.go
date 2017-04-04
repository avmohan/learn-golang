package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (x rot13Reader) Read(buf []byte) (int, error) {
	n, e := x.r.Read(buf)
	for i, v := range buf {
		switch {
		case v >= 'A' && v <= 'M' || v >= 'a' && v <= 'm':
			buf[i] = v + 13
		case v >= 'N' && v <= 'Z' || v >= 'n' && v <= 'z':
			buf[i] = v - 13
		}
	}
	return n, e
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
