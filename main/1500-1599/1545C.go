package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf1545C(in io.Reader, out io.Writer) {
	const mod = 998244353
	var T, n int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n)
		a := make([][]int, n*2+1)
		for i := range a {
			a[i] = make([]int, n+1)
		}
		c := make([][]int, n+1)
		for i := range c {
			c[i] = make([]int, n+1)
		}

		ans := 1
		for i := 1; i <= n*2; i++ {
			for j := 1; j <= n; j++ {
				Fscan(in, &a[i][j])
				c[a[i][j]][j]++
			}
		}

		v := make([]int, n*2+1)
		res := make([]int, 0, n)

		for i := 1; i <= n; i++ {
			pos := 0
			for j := 1; j <= n*2; j++ {
				if v[j] == 0 {
					for k := 1; k <= n; k++ {
						if c[a[j][k]][k] == 1 {
							pos = j
						}
					}
				}
			}

			if pos == 0 {
				for pos = 1; v[pos] != 0; pos++ {
				}
				ans = ans * 2 % mod
			}

			for j := 1; j <= n*2; j++ {
				if v[j] == 0 {
					for k := 1; k <= n; k++ {
						if a[j][k] == a[pos][k] {
							v[j] = 1
						}
					}
					if v[j] != 0 {
						for k := 1; k <= n; k++ {
							c[a[j][k]][k]--
						}
					}
				}
			}
			res = append(res, pos)
		}

		Fprintln(out, ans)
		for _, x := range res {
			Fprint(out, x, " ")
		}
		Fprintln(out)
	}
}

//func main() { cf1545C(bufio.NewReader(os.Stdin), os.Stdout) }
