package main

import (
	"slices"
	"sort"
)

// https://space.bilibili.com/206214
func findDisappearedNumbers(nums []int, lower, upper int) (ans [][]int) {
	nums = append(nums, lower-1, upper+1)
	slices.Sort(nums)

	l := sort.SearchInts(nums, lower)
	r := sort.SearchInts(nums, upper+1)

	for i := l; i <= r; i++ {
		if nums[i]-nums[i-1] > 1 {
			ans = append(ans, []int{nums[i-1] + 1, nums[i] - 1})
		}
	}
	return
}
