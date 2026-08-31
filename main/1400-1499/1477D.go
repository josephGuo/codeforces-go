package main

import (
	"bufio"
	. "fmt"
	"io"
	"slices"
)

// https://github.com/EndlessCheng
func cf1477D(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var T, n, m int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n, &m)
		g := make([][]int, n+2)
		for i := 1; i <= n; i++ {
			g[i] = append(g[i], i, n+1)
		}
		for i := 1; i <= m; i++ {
			var x, y int
			Fscan(in, &x, &y)
			g[x] = append(g[x], y)
			g[y] = append(g[y], x)
		}

		s1 := make([]int, n+2)
		s2 := make([]int, n+2)
		pos := make([]int, n+2)
		sz := make([]int, n+2)
		g2 := make([][]int, n+2)
		cnt := 0
		for i := 1; i <= n; i++ {
			if len(g[i]) == n+1 {
				cnt++
				s1[i] = cnt
				s2[i] = cnt
				continue
			}
			if pos[i] != 0 {
				continue
			}

			x := 0
			slices.Sort(g[i])
			for j := 0; j < len(g[i]); j++ {
				if g[i][j] != j+1 {
					x = j + 1
					break
				}
			}

			fa := pos[x]
			if fa == 0 {
				pos[x] = i
				pos[i] = i
				sz[i] = 2
			} else if fa == x {
				pos[i] = x
				sz[x]++
			} else if sz[fa] == 2 {
				sz[fa] = 0
				pos[fa] = x
				pos[x] = x
				pos[i] = x
				sz[x] = 3
			} else {
				sz[fa]--
				pos[x] = i
				pos[i] = i
				sz[i] = 2
			}
		}

		for i := 1; i <= n; i++ {
			if pos[i] != i {
				g2[pos[i]] = append(g2[pos[i]], i)
			}
		}

		for i := 1; i <= n; i++ {
			if len(g2[i]) > 0 {
				cnt++
				s1[i] = cnt
				for _, j := range g2[i] {
					s2[j] = cnt
					cnt++
					s1[j] = cnt
				}
				s2[i] = cnt
			}
		}

		for i := 1; i <= n; i++ {
			Fprint(out, s1[i], " ")
		}
		Fprintln(out)
		for i := 1; i <= n; i++ {
			Fprint(out, s2[i], " ")
		}
		Fprintln(out)
	}
}

//func main() { cf1477D(bufio.NewReader(os.Stdin), os.Stdout) }
