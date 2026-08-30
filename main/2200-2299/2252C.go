package main

import (
	"container/heap"
	. "fmt"
	"io"
	"sort"
)

// https://github.com/EndlessCheng
func cf2252C(in io.Reader, out io.Writer) {
	var T, n, m int
	for Fscan(in, &T); T > 0; T-- {
		Fscan(in, &n, &m)
		v := make([]int, n)
		for i := range v {
			Fscan(in, &v[i])
		}
		a := make([][]int, n)
		for i := range a {
			a[i] = make([]int, m)
			for j := range a[i] {
				Fscan(in, &a[i][j])
			}
		}

		ans := m
		h := &hp52{}
		s := 0
		for i := n - 1; i >= 0; i-- {
			for _, x := range a[i] {
				heap.Push(h, x)
				s += x
			}
			for s >= v[i] {
				ans = min(ans, h.Len())
				s -= heap.Pop(h).(int)
			}
		}
		Fprintln(out, ans)
	}
}

//func main() { cf2252C(bufio.NewReader(os.Stdin), os.Stdout) }
type hp52 struct{ sort.IntSlice }
func (h *hp52) Push(v any) { h.IntSlice = append(h.IntSlice, v.(int)) }
func (h *hp52) Pop() any   { a := h.IntSlice; v := a[len(a)-1]; h.IntSlice = a[:len(a)-1]; return v }
