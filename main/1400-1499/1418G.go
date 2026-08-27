package main

import (
	. "fmt"
	"io"
	"math/rand/v2"
)

// https://github.com/EndlessCheng
func cf1418G(in io.Reader, out io.Writer) {
	var n, v, ans int
	Fscan(in, &n)
	mp := make([][2]uint, n+1)
	for i := 1; i <= n; i++ {
		mp[i] = [2]uint{uint(rand.Uint64()), uint(rand.Uint64())}
	}

	a := make([]int, n)
	tot := make([]int, n+1)
	cnt := make([]int8, n+1)
	cntS := map[uint]int{0: 1}
	s := make([]uint, n+1)
	l := 0
	for i := range a {
		Fscan(in, &v)
		a[i] = v

		s[i+1] = s[i]
		tot[v]++
		if tot[v]%3 == 1 {
			s[i+1] ^= mp[v][0]
		} else if tot[v]%3 == 2 {
			s[i+1] ^= mp[v][1]
		} else {
			s[i+1] ^= mp[v][0] ^ mp[v][1]
		}

		cnt[v]++
		for cnt[v] > 3 {
			cntS[s[l]]--
			cnt[a[l]]--
			l++
		}

		ans += cntS[s[i+1]]
		cntS[s[i+1]]++
	}
	Fprint(out, ans)
}

//func main() { cf1418G(bufio.NewReader(os.Stdin), os.Stdout) }
