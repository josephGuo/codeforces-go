package main

// https://space.bilibili.com/206214
const mx = 100_001

var primeFactors = [mx][]int{}

func init() {
	for i := 2; i < mx; i++ {
		if primeFactors[i] == nil { // i 是质数
			for j := i; j < mx; j += i { // i 的倍数 j 有质因子 i
				primeFactors[j] = append(primeFactors[j], i)
			}
		}
	}
}

func longestSubarray(nums []int, k int) (ans int) {
	cnt := map[int]int{}
	left := 0
	for i, x := range nums {
		for _, p := range primeFactors[x] {
			cnt[p]++
		}
		for len(cnt) > k {
			for _, p := range primeFactors[nums[left]] {
				if cnt[p] > 1 {
					cnt[p]--
				} else {
					delete(cnt, p) // 保证 len(cnt) 是窗口内的不同质因子个数
				}
			}
			left++
		}
		ans = max(ans, i-left+1)
	}
	return
}
