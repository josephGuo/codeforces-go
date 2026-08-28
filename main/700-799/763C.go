package main

import (
	. "fmt"
	"io"
	"slices"
)

// https://github.com/EndlessCheng
func cf763C(in io.Reader, out io.Writer) {
	var M, n int
	pow := func(x, n int) int {
		res := 1
		for ; n > 0; n /= 2 {
			if n%2 > 0 {
				res = res * x % M
			}
			x = x * x % M
		}
		return res
	}

	Fscan(in, &M, &n)
	a := make([]int, n)
	s := [2]int{}
	for i := 0; i < n; i++ {
		Fscan(in, &a[i])
		s[0] = (s[0] + a[i]) % M
		s[1] = (s[1] + a[i]*a[i]) % M
	}

	if n == 1 {
		Fprintln(out, a[0], 0)
		return
	}
	if n == M {
		Fprintln(out, 0, 1)
		return
	}

	slices.Sort(a)
	b := make([]int, n)
o:
	for i := 1; i < n; i++ {
		d := a[i] - a[0]
		x := (s[0] - n*(n-1)/2%M*d%M + M) % M
		x = x * pow(n, M-2) % M
		if s[1] == (n*x%M*x+n*(n-1)%M*d%M*x+n*(n-1)*(2*n-1)/6%M*d%M*d)%M {
			b[0] = x
			for j := 1; j < n; j++ {
				b[j] = (b[j-1] + d) % M
			}
			slices.Sort(b)
			for j := 0; j < n; j++ {
				if a[j] != b[j] {
					continue o
				}
			}
			Fprintln(out, x, d)
			return
		}
	}
	Fprint(out, -1)
}

//func main() { cf763C(bufio.NewReader(os.Stdin), os.Stdout) }
