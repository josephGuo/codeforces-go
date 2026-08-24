package main

import (
	"math/bits"
)

// https://space.bilibili.com/206214
func isPalindromic(s string) bool {
	n := len(s)
	for i := range n/2 + 1 {
		if bits.Reverse8(s[i]) != s[n-1-i] {
			return false
		}
	}
	return true
}
