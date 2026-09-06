package main

// https://space.bilibili.com/206214
func maxProfit1(prices []int, strategy []int, k int) int64 {
	n := len(prices)
	sum := make([]int, n+1)
	sumSell := make([]int, n+1)
	for i, p := range prices {
		sum[i+1] = sum[i] + p*strategy[i]
		sumSell[i+1] = sumSell[i] + p
	}

	ans := sum[n] // 不修改
	for i := k; i <= n; i++ {
		res := sum[i-k] + sum[n] - sum[i] + sumSell[i] - sumSell[i-k/2]
		ans = max(ans, res)
	}
	return int64(ans)
}

func maxProfit(prices, strategy []int, k int) int64 {
	var total, extra, maxExtra int
	for i, p := range prices {
		total += strategy[i] * p
		if i < k/2 {
			continue
		}

		// 1. 入
		extra -= strategy[i-k/2] * prices[i-k/2] // 左边的窗口
		extra += (1 - strategy[i]) * p           // 右边的窗口

		left := i - k + 1
		if left < 0 { // 尚未形成第一个窗口
			continue
		}

		// 2. 更新
		maxExtra = max(maxExtra, extra)

		// 3. 出
		extra += strategy[left] * prices[left]               // 左边的窗口
		extra -= (1 - strategy[left+k/2]) * prices[left+k/2] // 右边的窗口
	}
	return int64(total + maxExtra)
}
