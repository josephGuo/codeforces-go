package main

import "math/rand/v2"

// https://space.bilibili.com/206214
func validSubarrays(nums []int, k int, l0 int, r0 int, q int) []bool {
	n := len(nums)
	sum := make([]uint64, n+1)
	hash := map[int]uint64{}
	for i, x := range nums {
		// 把 nums[i] 映射成一个随机的 uint64
		if _, ok := hash[x]; !ok {
			hash[x] = rand.Uint64()
		}
		sum[i+1] = sum[i] ^ hash[x]
	}

	calcLeft := func(k int) []int {
		lefts := make([]int, n)
		cnt := map[int]int{}
		l := 0
		for i, x := range nums {
			cnt[x]++
			for len(cnt) >= k {
				v := nums[l]
				if cnt[v] > 1 {
					cnt[v]--
				} else {
					delete(cnt, v) // 保证 len(cnt) 是窗口内的不同元素个数
				}
				l++
			}
			lefts[i] = l
		}
		return lefts
	}

	l1 := calcLeft(k + 1)
	l2 := calcLeft(k)

	ans := make([]bool, q)
	l, r := l0, r0
	for i := range ans {
		if i > 0 {
			g := r - l
			if ans[i-1] {
				g = l + r
			}
			l = (l ^ g) % n
			r = (r ^ g) % n
			if l > r {
				l, r = r, l
			}
		}
		ans[i] = sum[r+1] == sum[l] && l1[r] <= l && l < l2[r]
	}
	return ans
}
