package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2241E(in io.Reader, out io.Writer) {
	isSQ := [1e6 + 1]bool{}
	for i := 1; i*i < len(isSQ); i++ {
		isSQ[i*i] = true
	}

	var T, n int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n)
		a := make([]int, n)
		for i := range a {
			Fscan(in, &a[i])
		}
		g := make([][]int, n)
		for range n - 1 {
			var v, w int
			Fscan(in, &v, &w)
			v--
			w--
			g[v] = append(g[v], w)
			g[w] = append(g[w], v)
		}

		ans := 0
		var dfs func(int, int) int
		dfs = func(v, fa int) (size int) {
			sum2 := 0
			for _, w := range g[v] {
				if w == fa {
					continue
				}
				sz := dfs(w, v)
				if isSQ[a[v]] {
					ans += (sum2 + size) * sz
					sum2 += size * sz
				}
				size += sz
			}
			if isSQ[a[v]] {
				ans += (sum2 + size) * (n - 1 - size)
			}
			return size + 1
		}
		dfs(0, -1)
		Fprintln(out, ans)
	}
}

//func main() { cf2241E(bufio.NewReader(os.Stdin), os.Stdout) }
