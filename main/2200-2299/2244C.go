package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2244C(in io.Reader, out io.Writer) {
	gcd := func(a, b int) int {
		for a != 0 {
			a, b = b%a, a
		}
		return b
	}
	var T, n, x, y, v int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n, &x, &y)
		g := gcd(x, y)
		ok := true
		for i := 1; i <= n; i++ {
			Fscan(in, &v)
			if v%g != i%g {
				ok = false
			}
		}
		if ok {
			Fprintln(out, "YES")
		} else {
			Fprintln(out, "NO")
		}
	}
}

//func main() { cf2244C(bufio.NewReader(os.Stdin), os.Stdout) }
