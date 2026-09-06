package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2208C(in io.Reader, out io.Writer) {
	var T, n int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n)
		a := make([]struct{ c, p int }, n)
		for i := range a {
			Fscan(in, &a[i].c, &a[i].p)
		}

		f := 0.
		for i := n - 1; i >= 0; i-- {
			p := a[i]
			f = max(f, float64(p.c)+float64(100-p.p)/100*f)
		}
		Fprintf(out, "%.6f\n", f)
	}
}

//func main() { cf2208C(bufio.NewReader(os.Stdin), os.Stdout) }
