package main

// github.com/EndlessCheng/codeforces-go
func minSumOfLengths1(arr []int, target int) int {
	n := len(arr)
	// sufMin[i] 表示左端点 >= i 的和为 target 的最短子数组长度
	// 不存在子数组时，长度设为 n+1
	sufMin := make([]int, n)
	minLen := n + 1
	sum := 0
	r := n - 1
	for l := n - 1; l > 0; l-- {
		sum += arr[l]
		for sum > target {
			sum -= arr[r]
			r--
		}
		if sum == target {
			minLen = min(minLen, r-l+1)
		}
		sufMin[l] = minLen
	}

	ans := n + 1
	sum = 0
	l := 0
	for r, x := range arr[:n-1] {
		sum += x
		for sum > target {
			sum -= arr[l]
			l++
		}
		if sum == target {
			ans = min(ans, r-l+1+sufMin[r+1])
		}
	}

	if ans > n {
		return -1
	}
	return ans
}

func minSumOfLengths(arr []int, target int) int {
	n := len(arr)
	ans := n + 1
	// preMin[i] 表示右端点 < i 的和为 target 的最短子数组长度
	// 不存在子数组时，长度设为 n+1
	preMin := make([]int, n+1)
	preMin[0] = n + 1
	minLen := n + 1
	sum := 0
	l := 0
	for r, x := range arr {
		sum += x
		for sum > target {
			sum -= arr[l]
			l++
		}
		if sum == target {
			// 枚举第二个子数组的右端点为 r，用滑动窗口算出此时第二个子数组的左端点为 l
			// 那么第一个子数组的右端点必须 < l
			ans = min(ans, preMin[l]+r-l+1)
			minLen = min(minLen, r-l+1)
		}
		preMin[r+1] = minLen
	}

	if ans > n {
		return -1
	}
	return ans
}
