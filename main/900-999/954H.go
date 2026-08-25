package main

import (
	"bufio"
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf954H(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	const mod = 1_000_000_007
	pow := func(x, n int) int {
		res := 1
		for ; n > 0; n /= 2 {
			if n%2 > 0 {
				res = res * x % mod
			}
			x = x * x % mod
		}
		return res
	}

	var n int
	Fscan(in, &n)
	preMul := make([]int, n+1)
	sum := make([]int, 2*n+1)
	ans := make([]int, 2*n+1)
	a := make([]int, n+1)
	preMul[0] = 1
	for i := 1; i < n; i++ {
		Fscan(in, &a[i])
		preMul[i] = a[i] * preMul[i-1] % mod
	}

	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			ans[j-i] = (ans[j-i] + preMul[j]) % mod
		}
		inv := (mod + 1) / 2
		inv = inv * pow(preMul[i], mod-2) % mod
		inv = inv * pow(a[i+1], mod-2) % mod
		inv = inv * (a[i+1] - 1) % mod
		for j := 2 * i; j < 2*n-1; j++ {
			ans[j-2*i] = (ans[j-2*i] + inv*sum[j]) % mod
		}
		for j := i; j < n; j++ {
			sum[i+j] = (sum[i+j] + preMul[i]*preMul[j]) % mod
		}
		for j := i + 1; j < n; j++ {
			sum[i+j] = (sum[i+j] + preMul[i]*preMul[j]) % mod
		}
	}

	for i := 1; i < 2*n-1; i++ {
		Fprint(out, ans[i], " ")
	}
}

//func main() { cf954H(bufio.NewReader(os.Stdin), os.Stdout) }
