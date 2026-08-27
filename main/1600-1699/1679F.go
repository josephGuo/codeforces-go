package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf1679F(in io.Reader, out io.Writer) {
	const mod = 998244353
	var n, m, v, w int
	Fscan(in, &n, &m)
	g := [10]int{}
	for range m {
		Fscan(in, &v, &w)
		g[v] |= 1 << w
		g[w] |= 1 << v
	}

	f := make([]int, 1<<9)
	nf := make([]int, 1<<9)
	f[0] = 1
	for i := 1; i <= n; i++ {
		clear(nf)
		for j := range 1 << 9 {
			if f[j] == 0 {
				continue
			}
			for k := range 10 {
				if j>>k&1 == 0 {
					p := (j | (1<<k - 1)) & g[k]
					nf[p] = (nf[p] + f[j]) % mod
				}
			}
		}
		f, nf = nf, f
	}

	ans := 0
	for _, v := range f {
		ans += v
	}
	Fprint(out, ans%mod)
}

//func main() { cf1679F(bufio.NewReader(os.Stdin), os.Stdout) }
