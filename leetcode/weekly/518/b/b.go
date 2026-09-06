package main

// https://space.bilibili.com/206214
func countGoodRotations1(nums []int) (ans int) {
	n := len(nums)
	sum := make([]int, n*2+1)
	for i := range n * 2 {
		sum[i+1] = sum[i] + nums[i%n]
	}

	for i := n; i < n*2; i++ {
		// sum[i-n/2]-sum[i-n] > sum[i]-sum[i-n/2]
		if sum[i-n/2]*2 > sum[i]+sum[i-n] {
			ans++
		}
	}
	return
}

func countGoodRotations(nums []int) (ans int) {
	n := len(nums)
	sum1, sum2 := 0, 0
	for i := n / 2; i < n*2-1; i++ {
		// 1. 入
		sum1 += nums[(i-n/2)%n]
		sum2 += nums[i%n]

		left := i - n + 1
		if left < 0 { // 尚未形成第一个窗口
			continue
		}

		// 2. 更新答案
		if sum1 > sum2 {
			ans++
		}

		// 3. 出
		sum1 -= nums[left]
		sum2 -= nums[(left+n/2)%n]
	}
	return
}
