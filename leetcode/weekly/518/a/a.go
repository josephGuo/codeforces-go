package main

// https://space.bilibili.com/206214
func countRotations(s string, k int) int {
	n := len(s)
	c := 0
	if s[0] == s[n-1] {
		c = 1
	}
	for i := 1; i < n; i++ {
		if s[i-1] == s[i] {
			c++
		}
	}

	if c == k {
		return n - c
	}
	if c == k+1 {
		return c
	}
	return 0
}
