package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf582C(in io.Reader, out io.Writer) {
	gcd := func(a, b int) int {
		for a != 0 {
			a, b = b%a, a
		}
		return b
	}
	var n, ans int
	Fscan(in, &n)
	a := make([]int, n)
	for i := range a {
		Fscan(in, &a[i])
	}

	g := make([]int, n+1)
	for i := 1; i <= n; i++ {
		g[i] = gcd(i, n)
	}

	b := make([]int, n*2)
	c := make([]int, n+1)
	mx := make([]int, n)
	for i := 1; i < n; i++ {
		if n%i != 0 {
			continue
		}
		for j := range i {
			mx[j] = 0
		}
		for j := 1; j <= n; j++ {
			c[j] = c[j-1]
			if g[j] == i {
				c[j]++
			}
		}
		for j := range n {
			mx[j%i] = max(mx[j%i], a[j])
		}
		for j := range n {
			if a[j] == mx[j%i] {
				b[j] = 1
				b[j+n] = 1
			} else {
				b[j] = 0
				b[j+n] = 0
			}
		}
		for j := n*2 - 2; j >= 0; j-- {
			if b[j] != 0 {
				b[j] = b[j+1] + 1
			}
		}
		for j := range n {
			ans += c[min(n-1, b[j])]
		}
	}
	Fprint(out, ans)
}

//func main() { cf582C(bufio.NewReader(os.Stdin), os.Stdout) }
