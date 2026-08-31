package main

// https://space.bilibili.com/206214
func countSpecialIntegers1(nums []int) (ans int) {
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

func countSpecialIntegers(nums []int) (ans int) {
	cnt := map[int]int{}
	for i, x := range nums {
		if i == 0 || x != nums[i-1] {
			cnt[x]++
			if cnt[x] == 1 { // 首次遇到 x，暂时认为 x 是特殊整数
				ans++
			} else if cnt[x] == 2 { // x 不是特殊整数，撤销
				ans--
			}
		}
	}
	return
}
