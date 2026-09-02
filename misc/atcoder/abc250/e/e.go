package main

import (
	"bufio"
	. "fmt"
	"io"
	"math/rand/v2"
	"os"
)

// https://github.com/EndlessCheng
func run(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var n, q, x, y int
	Fscan(in, &n)
	mp := map[int]uint{}

	f := func() []uint {
		s := make([]uint, n+1)
		vis := map[int]bool{}
		for i := range n {
			Fscan(in, &x)
			s[i+1] = s[i]
			if !vis[x] {
				vis[x] = true
				if mp[x] == 0 {
					mp[x] = rand.Uint()
				}
				s[i+1] ^= mp[x]
			}
		}
		return s
	}

	sa := f()
	sb := f()

	Fscan(in, &q)
	for range q {
		Fscan(in, &x, &y)
		if sa[x] == sb[y] {
			Fprintln(out, "Yes")
		} else {
			Fprintln(out, "No")
		}
	}
}

func main() { run(bufio.NewReader(os.Stdin), os.Stdout) }
