package main

import (
	"bufio"
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf949E(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var n int
	Fscan(in, &n)
	a := make([]int, n+2)
	b := make([]int, n+2)
	for i := 1; i <= n; i++ {
		Fscan(in, &a[i])
	}
	n++

	ans := []int{}
	for ok, tot := true, 0; ok; {
		for s := 0; ; s++ {
			mn, mx := int(1e9), -int(1e9)
			for i := 1; i <= n; i++ {
				b[i] = (a[i] + 1<<30) & (2<<s - 1)
				if b[i] >= 1<<s {
					b[i] -= 1 << (s + 1)
				}
				mn = min(mn, b[i])
				mx = max(mx, b[i])
			}
			if mx-mn+1 <= (1 << s) {
				for i := 0; i < s; i++ {
					if mx>>i&1 != 0 {
						ans = append(ans, 1<<(i+tot))
					} else {
						ans = append(ans, -(1 << (i + tot)))
					}
				}
				ok = false
				for i := 1; i <= n; i++ {
					a[i] = (a[i] - b[i]) >> (s + 1)
					if a[i] != 0 {
						ok = true
					}
				}
				tot += s + 1
				break
			}
		}
	}

	Fprintln(out, len(ans))
	for _, x := range ans {
		Fprint(out, x, " ")
	}
}

//func main() { cf949E(bufio.NewReader(os.Stdin), os.Stdout) }
