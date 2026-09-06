package main

import "math"

// https://space.bilibili.com/206214
func minScoreTriangulation1(v []int) int {
	n := len(v)
	memo := make([][]int, n)
	for i := range memo {
		memo[i] = make([]int, n)
		for j := range memo[i] {
			memo[i][j] = -1 // -1 表示还没有计算过
		}
	}
	var dfs func(int, int) int
	dfs = func(i, j int) int {
		if i+1 == j { // 只有两个点，无法组成三角形
			return 0
		}
		p := &memo[i][j]
		if *p != -1 { // 之前计算过
			return *p
		}
		res := math.MaxInt
		for k := i + 1; k < j; k++ { // 枚举顶点 k
			res = min(res, dfs(i, k)+dfs(k, j)+v[i]*v[j]*v[k])
		}
		*p = res // 记忆化
		return res
	}
	return dfs(0, n-1)
}

func minScoreTriangulation(v []int) int {
	n := len(v)
	f := make([][]int, n)
	for i := range f {
		f[i] = make([]int, n)
	}
	for i := n - 3; i >= 0; i-- {
		for j := i + 2; j < n; j++ {
			f[i][j] = math.MaxInt
			for k := i + 1; k < j; k++ {
				f[i][j] = min(f[i][j], f[i][k]+f[k][j]+v[i]*v[j]*v[k])
			}
		}
	}
	return f[0][n-1]
}
