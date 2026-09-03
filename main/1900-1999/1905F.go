package main

import (
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf1905F(in io.Reader, out io.Writer) {
	var T, n, v int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n)
		pos := make([]int, n+1)
		for i := 1; i <= n; i++ {
			Fscan(in, &v)
			pos[v] = i
		}

		sufMin := make([]int, n+2)
		sufMin[n+1] = n + 1
		for i := n; i > 0; i-- {
			sufMin[i] = min(sufMin[i+1], pos[i])
		}

		cnt := map[[2]int]int{}
		var a, b, sum int
		for i := 1; i <= n; i++ {
			if pos[i] > a {
				b, a = a, pos[i]
			} else {
				b = max(b, pos[i])
			}
			if pos[i] == i && a == i {
				sum++
			}
			if pos[i] == i && a > i && b == i {
				cnt[[2]int{sufMin[i], a}]++
			}
			if pos[i] < i && a == i {
				cnt[[2]int{pos[i], i}]++
			}
			if pos[i] > i && b < i {
				cnt[[2]int{i, pos[i]}]++
			}
		}

		mx := -2
		for _, x := range cnt {
			mx = max(mx, x)
		}
		Fprintln(out, sum+mx)
	}
}

//func main() { cf1905F(bufio.NewReader(os.Stdin), os.Stdout) }
