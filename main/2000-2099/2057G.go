package main

import (
	"bufio"
	. "fmt"
	"io"
)

// https://github.com/EndlessCheng
func cf2057G(in io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	var T, n, m int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n, &m)
		a := make([][]byte, n)
		for i := range a {
			Fscan(in, &a[i])
		}

		pos := [5][]int{}
		for i, row := range a {
			for j, b := range row {
				if b != '#' {
					continue
				}
				d := (i + j*2) % 5
				p := i*m + j
				pos[d] = append(pos[d], p)
				if i == 0 || a[i-1][j] == '.' {
					t := (d + 4) % 5
					pos[t] = append(pos[t], p)
				}
				if i == n-1 || a[i+1][j] == '.' {
					t := (d + 1) % 5
					pos[t] = append(pos[t], p)
				}
				if j == 0 || a[i][j-1] == '.' {
					t := (d + 3) % 5
					pos[t] = append(pos[t], p)
				}
				if j == m-1 || a[i][j+1] == '.' {
					t := (d + 2) % 5
					pos[t] = append(pos[t], p)
				}
			}
		}

		best := 0
		for i := 1; i < 5; i++ {
			if len(pos[i]) < len(pos[best]) {
				best = i
			}
		}

		for _, p := range pos[best] {
			a[p/m][p%m] = 'S'
		}
		for _, row := range a {
			Fprintf(out, "%s\n", row)
		}
	}
}

//func main() { cf2057G(bufio.NewReader(os.Stdin), os.Stdout) }
