package main

import (
	"bufio"
	. "fmt"
	"io"
	"slices"
)

// https://github.com/EndlessCheng
func cf2252D(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var T, n int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n)
		a := make([]int, n)
		d := make([]int, n-1)
		for i := range a {
			Fscan(in, &a[i])
			if i > 0 {
				d[i-1] = a[i] - a[i-1]
			}
		}

		st := 0
		for i := range n - 1 {
			if i == n-2 || a[i]&1 != a[i+2]&1 {
				slices.Sort(d[st : i+1])
				st = i + 1
			}
		}

		s := a[0]
		Fprint(out, s)
		for _, v := range d {
			s += v
			Fprint(out, " ", s)
		}
		Fprintln(out)
	}
}

//func main() { cf2252D(bufio.NewReader(os.Stdin), os.Stdout) }
