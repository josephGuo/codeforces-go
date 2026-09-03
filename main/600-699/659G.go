package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf659G(in io.Reader, out io.Writer) {
	const mod = 1_000_000_007
	var n int
	Fscan(in, &n)
	a := make([]int, n+2)
	for i := 1; i <= n; i++ {
		Fscan(in, &a[i])
		a[i]--
	}

	ans := 0
	f := 0
	for i := 1; i <= n; i++ {
		ans += f*min(a[i-1], a[i])%mod + a[i]
		f = (f*min(a[i-1], a[i], a[i+1]) + min(a[i], a[i+1])) % mod
	}
	Fprint(out, ans%mod)
}

//func main() { cf659G(bufio.NewReader(os.Stdin), os.Stdout) }
