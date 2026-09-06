package main

import (
	"bufio"
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf773D(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var n int
	Fscan(in, &n)
	a := make([][]int, n+1)
	for i := range a {
		a[i] = make([]int, n+1)
	}
	mn := int(1e9)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			Fscan(in, &a[i][j])
			a[j][i] = a[i][j]
			mn = min(mn, a[i][j])
		}
	}

	dis := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dis[i] = 1e18
		for j := 1; j <= n; j++ {
			if i != j {
				a[i][j] -= mn
				dis[i] = min(dis[i], a[i][j]*2)
			}
		}
	}

	vis := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		k := 0
		for j := 1; j <= n; j++ {
			if !vis[j] && (k == 0 || dis[j] < dis[k]) {
				k = j
			}
		}
		vis[k] = true
		for j := 1; j <= n; j++ {
			dis[j] = min(dis[j], dis[k]+a[k][j])
		}
	}

	for i := 1; i <= n; i++ {
		Fprintln(out, dis[i]+(n-1)*mn)
	}
}

//func main() { cf773D(bufio.NewReader(os.Stdin), os.Stdout) }
