package main

import (
	. "fmt"
	"io"
	"strings"
)

// https://github.com/EndlessCheng
func cf1906L(in io.Reader, out io.Writer) {
	var n, k int
	Fscan(in, &n, &k)
	if k*2 < n || k == n {
		Fprint(out, -1)
		return
	}

	m := n - k
	Fprint(out, strings.Repeat("()", (k-n/2+1)/2),
		strings.Repeat("(", m),
		strings.Repeat(")", m),
		strings.Repeat("()", (k-n/2)/2))
}

//func main() { cf1906L(os.Stdin, os.Stdout) }
