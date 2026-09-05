package main

// https://space.bilibili.com/206214
func maxPalindromesDP(s string, k int) int {
	n := len(s)
	// f[i] 表示从 s[:i] 中选出的回文子串的最大数目
	f := make([]int, n+1)
	for i := k; i <= n; i++ {
		f[i] = f[i-1] // 不考虑 s[i-1]
		if isPalindrome(s[i-k : i]) {
			f[i] = max(f[i], f[i-k]+1)
		}
		if i > k && isPalindrome(s[i-k-1:i]) {
			f[i] = max(f[i], f[i-k-1]+1)
		}
	}
	return f[n]
}

func isPalindrome(s string) bool {
	n := len(s)
	for i := range n / 2 {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

func maxPalindromes1(s string, k int) (ans int) {
	n := len(s)
	for i := 0; i <= n-k; {
		if isPalindrome(s[i : i+k]) {
			ans++
			i += k // 计算 s[i+k:] 中的最优方案
		} else if i < n-k && isPalindrome(s[i:i+k+1]) {
			// 如果跳过不选，即使 s[i+1:i+1+k] 是回文串，剩余内容仍然是 s[i+k+1:]，并不会更优
			// 所以不需要考虑跳过 s[i:i+k+1] 的情况
			ans++
			i += k + 1 // 计算 s[i+k+1:] 中的最优方案
		} else {
			i++ // 计算 s[i+1:] 中的最优方案
		}
	}
	return
}

func maxPalindromes(s string, k int) (ans int) {
	// Manacher 模板
	// 将 s 改造为 t，这样就不需要讨论 len(s) 的奇偶性，因为新串 t 的每个回文子串都是奇回文串（都有回文中心）
	// s 和 t 的下标转换关系：
	// (si+1)*2 = ti
	// ti/2-1 = si
	// ti 为偶数，对应奇回文串（从 2 开始）
	// ti 为奇数，对应偶回文串（从 3 开始）
	n := len(s)
	t := append(make([]byte, 0, n*2+3), '^')
	for _, c := range s {
		t = append(t, '#', byte(c))
	}
	t = append(t, '#', '$')

	// 定义一个奇回文串的回文半径=(长度+1)/2，即保留回文中心，去掉一侧后的剩余字符串的长度
	// halfLen[i] 表示在 t 上的以 t[i] 为回文中心的最长回文子串的回文半径
	// 即 [i-halfLen[i]+1,i+halfLen[i]-1] 是 t 上的一个回文子串
	halfLen := make([]int, len(t)-2)
	halfLen[1] = 1

	// boxR 表示当前右边界下标最大的回文子串的右边界下标+1
	// boxM 为该回文子串的中心位置
	// 二者的关系为 boxR = boxM + halfLen[boxM]
	boxM, boxR := 0, 0
	for i := 2; i < len(halfLen); i++ {
		hl := 1
		if i < boxR {
			// 记 i 关于 boxM 的对称位置 i'=boxM*2-i
			// 若以 i' 为中心的最长回文子串范围超出了以 boxM 为中心的回文串的范围
			// 则 halfLen[i] 应先初始化为已知的回文半径 boxR-i，然后再继续暴力匹配
			// 否则 halfLen[i] 与 halfLen[i'] 相等
			hl = min(boxR-i, halfLen[boxM*2-i])
		}

		// 暴力扩展
		for t[i-hl] == t[i+hl] {
			hl++
			boxM, boxR = i, i+hl
		}

		halfLen[i] = hl
	}

	// 判断子串 s[l:r]（左闭右开）是否为回文串
	// 根据下标转换关系得到子串 s[l:r] 在 t 中对应的回文中心下标为 l+r+1
	// t 中回文子串的长度为 hl*2-1
	// 由于其中 # 的数量总是比字母的数量多 1
	// 因此其在 s 中对应的回文子串的长度为 hl-1
	isPalindrome := func(l, r int) bool {
		return halfLen[l+r+1] > r-l // halfLen[l+r+1]-1 >= r-l
	}

	for i := 0; i <= n-k; {
		if isPalindrome(i, i+k) {
			ans++
			i = i + k // 计算 s[i+k:] 中的最优方案
		} else if i < n-k && isPalindrome(i, i+k+1) {
			// 如果跳过不选，即使 s[i+1:i+1+k] 是回文串，剩余内容仍然是 s[i+k+1:]，并不会更优
			// 所以不需要考虑跳过 s[i:i+k+1] 的情况
			ans++
			i = i + k + 1 // 计算 s[i+k+1:] 中的最优方案
		} else {
			i++ // 计算 s[i+1:] 中的最优方案
		}
	}
	return
}
