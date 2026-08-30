package main

import "sort"

// https://space.bilibili.com/206214
func countValidSplits(nums []int, skip int) (cnt int) {
	n := len(nums)
	// suf[i] 是后缀 [i,n-1]（除去 skip）的 GCD
	suf := make([]int, n+1)
	for j := n - 1; j >= 0; j-- {
		if j != skip {
			suf[j] = gcd(suf[j+1], nums[j])
		} else {
			suf[j] = suf[j+1]
		}
	}

	pre := 0
	for j, x := range nums {
		if j != skip {
			pre = gcd(pre, x)
			// 现在 pre 是前缀 [0,j]（除去 skip）的 GCD
			if pre == suf[j+1] {
				cnt++
			}
		}
	}
	return
}

func maxValidSplits(nums []int) int {
	ans := countValidSplits(nums, -1) // 不删除元素

	// countValidSplits 只会调用 O(log max(nums)) 次
	g := 0
	for i, x := range nums {
		if g > 0 && x%g == 0 { // x 不改变前缀 GCD
			continue // 把 x 删了 ans 也不会变大
		}
		g = gcd(g, x)
		ans = max(ans, countValidSplits(nums, i)) // 删 x
	}

	return ans
}

func gcd(a, b int) int {
	for a != 0 {
		a, b = b%a, a
	}
	return b
}

//

func maxValidSplits2(nums []int) int {
	n := len(nums)
	preGcd := make([]int, n)
	g := 0
	for i, x := range nums {
		g = gcd(g, x)
		preGcd[i] = g
	}

	sufGcd := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		sufGcd[i] = gcd(sufGcd[i+1], nums[i])
	}

	// 不删任何数
	allGcd := sufGcd[0]
	p := sort.Search(n, func(i int) bool { return preGcd[i] == allGcd }) // [p,n-1] 都是 allGcd
	q := sort.SearchInts(sufGcd[:n], allGcd+1) - 1                       // [0,q] 都是 allGcd
	ans := max(q-p, 0)                                                   // 满足 i >= p 且 i+1 <= q 的 i 的个数

	for i := range n {
		if i > 0 && preGcd[i] == preGcd[i-1] {
			continue
		}

		// 删除 nums[i]
		newG := sufGcd[i+1]
		if i > 0 {
			newG = gcd(preGcd[i-1], sufGcd[i+1])
		}

		g = 0
		for j, x := range nums {
			if j == i {
				continue
			}
			g = gcd(g, x)
			if g == newG {
				p = j
				break
			}
		}

		g = 0
		for j := n - 1; j >= 0; j-- {
			if j == i {
				continue
			}
			g = gcd(g, nums[j])
			if g == newG {
				q = j
				break
			}
		}

		res := q - p
		if p <= i && i < q {
			res--
		}
		ans = max(ans, res)
	}

	return ans
}
