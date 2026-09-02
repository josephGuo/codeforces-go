package main

import (
	"slices"
	"sort"
)

// github.com/EndlessCheng/codeforces-go
func maxNumOfSubstrings(s string) (ans []string) {
	// 记录每种字母的出现位置
	pos := [26][]int{}
	for i, b := range s {
		b -= 'a'
		pos[b] = append(pos[b], i)
	}

	// 构建有向图
	g := [26][]int{}
	for i, p := range pos {
		if p == nil {
			continue
		}
		l, r := p[0], p[len(p)-1]
		for j, q := range pos {
			if j == i {
				continue
			}
			k := sort.SearchInts(q, l)
			// [l,r] 包含第 j 个小写字母
			if k < len(q) && q[k] <= r {
				g[i] = append(g[i], j)
			}
		}
	}

	// 遍历有向图
	vis := [26]bool{}
	var l, r int
	var dfs func(int)
	dfs = func(x int) {
		vis[x] = true
		p := pos[x]
		l = min(l, p[0]) // 合并区间
		r = max(r, p[len(p)-1])
		for _, y := range g[x] {
			if !vis[y] {
				dfs(y)
			}
		}
	}

	type pair struct{ l, r int }
	intervals := []pair{}
	for i, p := range pos {
		if p == nil {
			continue
		}
		// 如果要包含第 i 个小写字母，最终得到的区间是什么？
		vis = [26]bool{}
		l, r = len(s), 0
		dfs(i)
		intervals = append(intervals, pair{l, r})
	}

	// 435. 无重叠区间
	// 直接计算最多能选多少个区间
	slices.SortFunc(intervals, func(a, b pair) int { return a.r - b.r })
	preR := -1
	for _, p := range intervals {
		if p.l > preR {
			ans = append(ans, s[p.l:p.r+1])
			preR = p.r
		}
	}
	return
}
