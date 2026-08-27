package main

import (
	"bufio"
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf1620F(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var T, n int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n)
		a := make([]int, n+1)
		f := make([][2]int, n+1)
		g := make([][2]bool, n+1)
		for i := 1; i <= n; i++ {
			f[i][0], f[i][1] = n+1, n+1
			Fscan(in, &a[i])
		}

		f[1][0], f[1][1] = -n-1, -n-1
		for i := 2; i <= n; i++ {
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					v := a[i]
					if x == 0 {
						v = -a[i]
					}
					pre := a[i-1]
					if y == 0 {
						pre = -a[i-1]
					}
					if pre < v && f[i][x] > f[i-1][y] {
						f[i][x] = f[i-1][y]
						g[i][x] = y != 0
					}
					if f[i-1][y] < v && f[i][x] > pre {
						f[i][x] = pre
						g[i][x] = y != 0
					}
				}
			}
		}

		if f[n][0] > n && f[n][1] > n {
			Fprintln(out, "NO")
			continue
		}
		Fprintln(out, "YES")
		t := 0
		if f[n][0] > n {
			t = 1
		}
		for i := n; i >= 1; i-- {
			if t == 0 {
				a[i] *= -1
			}
			if !g[i][t] {
				t = 0
			} else {
				t = 1
			}
		}
		for i := 1; i <= n; i++ {
			Fprint(out, a[i], " ")
		}
		Fprintln(out)
	}
}

//func main() { cf1620F(bufio.NewReader(os.Stdin), os.Stdout) }
