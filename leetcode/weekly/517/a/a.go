package main

// https://space.bilibili.com/206214
func countSpecialIntegers(nums []int) (ans int) {
	pos := map[int][]int{}
	for i, x := range nums {
		pos[x] = append(pos[x], i)
	}

	for _, p := range pos {
		if p[len(p)-1]-p[0]+1 == len(p) {
			ans++
		}
	}
	return
}
