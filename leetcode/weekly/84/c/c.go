package main

// github.com/EndlessCheng/codeforces-go
func largestOverlap(img1, img2 [][]int) (ans int) {
	n := len(img1)
	for dx := 1 - n; dx < n; dx++ {
		for dy := 1 - n; dy < n; dy++ {
			cnt1 := 0
			for i := max(-dx, 0); i < min(n-dx, n); i++ {
				for j := max(-dy, 0); j < min(n-dy, n); j++ {
					// 两个数都是 1，才能让 cnt1 增加 1
					cnt1 += img1[i][j] * img2[i+dx][j+dy]
				}
			}
			ans = max(ans, cnt1)
		}
	}
	return
}
