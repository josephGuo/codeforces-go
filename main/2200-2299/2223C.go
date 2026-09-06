package main

import (
	"bufio"
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2223C(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
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
		n, q := rd(), rd()
		pa := make([]int, n)
		g := make([][]int, n)
		for w := 1; w < n; w++ {
			v := rd() - 1
			pa[w] = v
			g[v] = append(g[v], w)
		}
		dis := make([]int, n)
		for i := 1; i < n; i++ {
			dis[i] = dis[pa[i]] + rd()
		}
		m := make([]int, q)
		group := make([]int, q)
		for i := range m {
			m[i] = rd()
			group[i] = i
		}

		var dfs func(int, int, int, []int)
		dfs = func(v, rem, sonsLcm int, group []int) {
			if g[v] == nil {
				for _, i := range group {
					m[i] = v
				}
				return
			}

			sons := len(g[v])
			if sonsLcm > 1e18 || sonsLcm%sons == 0 {
				dfs(g[v][(rem+dis[v])%sons], rem, sonsLcm, group)
				return
			}

			newGroup := make([][]int, sons)
			for _, i := range group {
				r := (m[i] + dis[v]) % sons
				newGroup[r] = append(newGroup[r], i)
			}
			sonsLcm = lcm23(sonsLcm, sons)
			for i, ng := range newGroup {
				if ng != nil {
					dfs(g[v][i], m[ng[0]]%sonsLcm, sonsLcm, ng)
				}
			}
		}
		dfs(0, 0, 1, group)

		for _, v := range m {
			Fprint(out, v+1, " ")
		}
		Fprintln(out)
	}
}

//func main() { cf2223C(os.Stdin, os.Stdout) }

func gcd23(a, b int) int {
	for a != 0 {
		a, b = b%a, a
	}
	return b
}

func lcm23(a, b int) int {
	a /= gcd23(b, a)
	if a > 1e18/b {
		return 1e18 + 1
	}
	return a * b
}
