package main

import "math"

// https://space.bilibili.com/206214
func minOperations(nums []int, sum int) int {
	f := make([]int, sum+1)
	for i := 1; i <= sum; i++ {
		f[i] = math.MaxInt / 2
	}

	for _, x := range nums {
		for i := sum; i > 0; i-- {
			// 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
			// 本题是分组背包，要枚举选哪个物品（枚举除法操作次数为 a，乘法操作次数为 b）
			for a := 0; x>>a > 0; a++ {
				for b := 0; x>>a<<b <= i; b++ {
					// 物品体积为 x>>a<<b，价值为 a+b
					f[i] = min(f[i], f[i-x>>a<<b]+a+b)
				}
			}
		}
	}

	if f[sum] == math.MaxInt/2 {
		return -1
	}
	return f[sum]
}

func minOperationsOld(nums []int, sum int) int {
	f := make([]int, sum+1)
	for i := 1; i <= sum; i++ {
		f[i] = math.MaxInt / 2
	}

	for _, x := range nums {
		// todo map 优化
		for i := sum; i > 0; i-- {
			for a := 0; x>>a > 0; a++ {
				// 如果 a > 0 且 x>>a&1 是 0 且 x>>(a-1)&1 也是 0
				// 此时有 x>>a<<b = x>>(a-1)<<(b-1)，所以只需考虑 b = 0，更大的 b 在 x>>(a-1) 中枚举过了
				if a > 0 && x>>a&1 == 0 && x>>(a-1)&1 == 0 && x>>a <= i {
					f[i] = min(f[i], f[i-x>>a]+a)
					continue
				}
				for b := 0; x>>a<<b <= i; b++ {
					// 物品体积为 x>>a<<b，价值为 a+b
					f[i] = min(f[i], f[i-x>>a<<b]+a+b)
				}
			}
		}
	}

	if f[sum] == math.MaxInt/2 {
		return -1
	}
	return f[sum]
}
