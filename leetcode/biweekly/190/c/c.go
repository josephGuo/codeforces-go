package main

import (
	"bytes"
	"math/bits"
)

// https://space.bilibili.com/206214
func largestString(nums []int) []string {
	ans := make([]string, len(nums))
	for idx, x := range nums {
		// 单独处理 'z'
		s := bytes.Repeat([]byte{'z'}, x>>25)
		// 然后从 'y' 到 'a'
		for i := min(24, bits.Len(uint(x))-1); i >= 0; i-- {
			if x>>i&1 > 0 { // x 的 i 位是 1，所以有字母 i
				s = append(s, 'a'+byte(i))
			}
		}
		ans[idx] = string(s)
	}
	return ans
}
