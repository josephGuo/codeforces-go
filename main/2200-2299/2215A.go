package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2215A(in io.Reader, out io.Writer) {
	var T, n, k, p, q int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n, &k, &p, &q)
		a := make([]int, n)
		var tot, s, sp, sqp int
		ans := int(1e18)
		for i := range a {
			Fscan(in, &a[i])
			v := a[i] % p
			w := a[i] % q % p
			tot += min(v, w)
			s += min(v, w)
			sp += v
			sqp += w
			l := i - k + 1
			if l < 0 {
				continue
			}
			ans = min(ans, min(sp, sqp)-s)
			v = a[l] % p
			w = a[l] % q % p
			s -= min(v, w)
			sp -= v
			sqp -= w
		}
		Fprintln(out, tot+ans)
	}
}

//func main() { cf2215A(bufio.NewReader(os.Stdin), os.Stdout) }
