package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf913F(in io.Reader, out io.Writer) {
	const mod = 998244353
	pow := func(x, n int) int {
		res := 1
		for ; n > 0; n /= 2 {
			if n%2 > 0 {
				res = res * x % mod
			}
			x = x * x % mod
		}
		return res
	}

	var n, a, b int
	Fscan(in, &n, &a, &b)
	P := a * pow(b, mod-2) % mod
	g := make([]int, n+1)
	c := make([]int, n+1)
	f := make([][]int, n+1)
	for i := range f {
		f[i] = make([]int, n+1)
	}
	for i := 0; i <= n; i++ {
		f[i][0] = 1
		for j := 1; j <= i; j++ {
			f[i][j] = (f[i-1][j-1]*pow(mod+1-P, i-j) + f[i-1][j]*pow(P, j)) % mod
		}
	}
	for i := 1; i <= n; i++ {
		c[i] = 1
		for j := 1; j < i; j++ {
			c[i] = (c[i] + (mod-c[j])*f[i][j]) % mod
		}
	}
	for i := 2; i <= n; i++ {
		for j := 1; j < i; j++ {
			g[i] = (g[i] + c[j]*f[i][j]%mod*(j*(i-j)+j*(j-1)/2+g[j]+g[i-j])) % mod
		}
		g[i] = (g[i] + i*(i-1)/2*c[i]) % mod * pow(mod+1-c[i], mod-2) % mod
	}
	Fprint(out, g[n])
}

//func main() { cf913F(bufio.NewReader(os.Stdin), os.Stdout) }
