package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf464D(in io.Reader, out io.Writer) {
	var n int
	var k float64
	Fscan(in, &n, &k)
	f := [1002]float64{}
	for range n {
		for x := 1; x <= 1000; x++ {
			f[x] *= (float64(x)/float64(x+1) + k - 1) / k
			f[x] += (float64(x)/2 + (f[x+1]+float64(x))/float64(x+1)) / k
		}
	}
	Fprintf(out, "%.9f", k*f[1])
}

//func main() { cf464D(bufio.NewReader(os.Stdin), os.Stdout) }
