package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2241C(in io.Reader, out io.Writer) {
	var T, n int
	var s string
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n, &s)
		cnt := 1
		for i := 1; i < n; i++ {
			cnt += int(s[i] ^ s[i-1])
		}
		if cnt == 2 {
			Fprintln(out, 2)
		} else {
			Fprintln(out, 1)
		}
	}
}

//func main() { cf2241C(bufio.NewReader(os.Stdin), os.Stdout) }
