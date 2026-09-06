package main

import "slices"

// https://space.bilibili.com/206214
func carFleet1(target int, position, speed []int) int {
	type pair struct{ p, v int }
	a := make([]pair, len(position))
	for i, p := range position {
		a[i] = pair{p, speed[i]}
	}
	slices.SortFunc(a, func(a, b pair) int { return a.p - b.p })

	st := []pair{}
	for _, p := range a {
		for len(st) > 0 {
			q := st[len(st)-1]
			if (target-q.p)*p.v > (target-p.p)*q.v {
				break
			}
			st = st[:len(st)-1]
		}
		st = append(st, p)
	}
	return len(st)
}

func carFleet(target int, position, speed []int) int {
	n := len(position)
	type pair struct{ p, v int }
	a := make([]pair, n)
	for i, p := range position {
		a[i] = pair{p, speed[i]}
	}
	slices.SortFunc(a, func(a, b pair) int { return a.p - b.p })

	mx := a[n-1]
	ans := 1
	for i := n - 2; i >= 0; i-- {
		p := a[i]
		if (target-p.p)*mx.v > (target-mx.p)*p.v {
			mx = p
			ans++
		}
	}
	return ans
}
