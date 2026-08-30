package main

import (
	"math"
	"math/bits"
)

// https://space.bilibili.com/206214
func minOperations(nums []int, sum int) int {
	f := make([]int, sum+1)
	for i := 1; i <= sum; i++ {
		f[i] = math.MaxInt / 2
	}

	for _, x := range nums {
		w := bits.Len(uint(x))
		for i := sum; i > 0; i-- {
			// 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
			// 本题是分组背包，要枚举选哪个物品（枚举乘了 a 次或者除了 a 次）
			for a := 0; x<<a <= i; a++ {
				// 物品体积为 x<<a，价值为 a
				f[i] = min(f[i], f[i-x<<a]+a)
			}
			// 从小到大枚举 x>>a，方便在 x>>a > i 时跳出循环
			for a := w - 1; a > 0 && x>>a <= i; a-- {
				// 物品体积为 x>>a，价值为 a
				f[i] = min(f[i], f[i-x>>a]+a)
			}
		}
	}

	if f[sum] == math.MaxInt/2 {
		return -1
	}
	return f[sum]
}
