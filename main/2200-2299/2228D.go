package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2228D(in io.Reader, out io.Writer) {
	buf := make([]byte, 4096)
	_i, _n := 0, 0
	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return 0
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}
	rd := func() (x int) {
		b := rc()
		for ; '0' > b; b = rc() {
		}
		for ; '0' <= b; b = rc() {
			x = x*10 + int(b&15)
		}
		return
	}

	for range rd() {
		n := rd()
		a := make([]struct{ l, r int }, n)
		for i := range a {
			a[i].l = 1e9
		}
		has := make([]bool, n+1)
		for range n {
			x, y := rd()-1, rd()
			a[x].l = min(a[x].l, y)
			a[x].r = max(a[x].r, y)
			has[y] = true
		}

		sum := make([]int, n+2)
		for i, b := range has {
			sum[i+1] = sum[i]
			if b {
				sum[i+1]++
			}
		}

		sufMax := make([]int, n+1)
		sufMin := make([]int, n+1)
		sufMin[n] = 1e9
		for i := n - 1; i >= 0; i-- {
			sufMin[i] = min(sufMin[i+1], a[i].l)
			sufMax[i] = max(sufMax[i+1], a[i].r)
		}

		ans := 0
		preMin, preMax := n+1, 0
		for i, p := range a {
			if p.l > n {
				continue
			}
			preMin = min(preMin, p.l)
			preMax = max(preMax, p.r)
			l := max(preMin, sufMin[i+1])
			r := min(preMax, sufMax[i+1])
			if l < r {
				ans += sum[r] - sum[l]
			}
		}
		Fprintln(out, ans)
	}
}

//func main() { cf2228D(os.Stdin, os.Stdout) }
