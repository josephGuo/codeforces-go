package main

import (
	"bufio"
	. "fmt"
	"io"
	"math/rand/v2"
)

// https://github.com/EndlessCheng
func cf2014H(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var T, n, q, v, l, r int
	mp := map[int]uint64{}
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n, &q)
		s := make([]uint64, n+1)
		for i := range n {
			Fscan(in, &v)
			if mp[v] == 0 {
				mp[v] = rand.Uint64()
			}
			s[i+1] = s[i] ^ mp[v]
		}
		for range q {
			Fscan(in, &l, &r)
			l--
			if (r-l)&1 == 0 && s[l] == s[r] {
				Fprintln(out, "YES")
			} else {
				Fprintln(out, "NO")
			}
		}
	}
}

//func main() { cf2014H(bufio.NewReader(os.Stdin), os.Stdout) }
