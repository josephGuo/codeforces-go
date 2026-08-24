package main

import (
	"bufio"
	. "fmt"
	"io"
	"slices"
)

// https://github.com/EndlessCheng
const mod42 = 998244353

func pow42(x, n int) (res int) {
	res = 1
	for ; n > 0; n /= 2 {
		if n%2 > 0 {
			res = res * x % mod42
		}
		x = x * x % mod42
	}
	return
}

type comb42 struct{ _f, _invF []int }

func newComb42(mx int) *comb42 {
	c := &comb42{[]int{1}, []int{1}}
	c._grow(mx)
	return c
}

func (c *comb42) _grow(mx int) {
	n := len(c._f)
	c._f = slices.Grow(c._f, mx+1)[:mx+1]
	for i := n; i <= mx; i++ {
		c._f[i] = c._f[i-1] * i % mod42
	}
	c._invF = slices.Grow(c._invF, mx+1)[:mx+1]
	c._invF[mx] = pow42(c._f[mx], mod42-2)
	for i := mx; i > n; i-- {
		c._invF[i-1] = c._invF[i] * i % mod42
	}
}

func (c *comb42) f(n int) int {
	if n >= len(c._f) {
		c._grow(n * 2)
	}
	return c._f[n]
}

func (c *comb42) invF(n int) int {
	if n >= len(c._f) {
		c._grow(n * 2)
	}
	return c._invF[n]
}

func (c *comb42) c(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	return c.f(n) * c.invF(k) % mod42 * c.invF(n-k) % mod42
}

func cf1942G(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var cm = newComb42(0)
	var t int
	Fscan(in, &t)

	type query struct{ a, b, c int }
	qs := make([]query, t)
	mx := 0
	for i := 0; i < t; i++ {
		Fscan(in, &qs[i].a, &qs[i].b, &qs[i].c)
		mx = max(mx, qs[i].a+qs[i].c+5)
	}

	for _, q := range qs {
		a, c := q.a, q.c
		s := 0
		for i := 0; i <= min(a, c); i++ {
			f := (cm.c(2*i+4, i) - cm.c(2*i+4, i-1) + mod42) % mod42
			s = (s + f*cm.c(i+5, 5)%mod42*cm.c(a+c-2*i, a-i)) % mod42
		}
		if a < c {
			f := (cm.c(a+c+5, c) - cm.c(a+c+5, a) + mod42) % mod42
			s = (s + f*cm.c(a+5, 5)) % mod42
		}
		den := cm.c(a+c+5, c) * cm.c(a+5, 5) % mod42
		Fprintln(out, s*pow42(den, mod42-2)%mod42)
	}
}

//func main() { cf1942G(bufio.NewReader(os.Stdin), os.Stdout) }
