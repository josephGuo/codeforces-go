package main

import "math"

// https://space.bilibili.com/206214
const mod = 1_000_000_007

func pow(x, n int) int {
	res := 1
	for ; n > 0; n /= 2 {
		if n%2 > 0 {
			res = res * x % mod
		}
		x = x * x % mod
	}
	return res
}

func sumDecoded(nums []int64) (ans int) {
	for _, X := range nums {
		x := int(X)
		d := x / 10
		// 计算 d 的十进制长度
		lengthD := 0
		for v := d; v > 0; v /= 10 {
			lengthD++
		}
		pow10 := int(math.Pow10(lengthD - x%10))
		// 根据 pow10 求出 x = d/pow10 和 y = d%pow10
		ans += pow(d/pow10, d%pow10)
	}
	return ans % mod
}
