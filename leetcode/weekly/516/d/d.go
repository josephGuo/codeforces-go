package main

import "math/rand/v2"

// https://space.bilibili.com/206214
// 模板来源 https://leetcode.cn/discuss/post/3583665/
type fenwick []int

func newFenwickTree(n int) fenwick {
	return make(fenwick, n+1) // 使用下标 1 到 n
}

// a[i] 增加 val
// 时间复杂度 O(log n)
func (f fenwick) update(i, val int) {
	for i++; i < len(f); i += i & -i {
		f[i] += val
	}
}

// 求前缀和 a[1] + ... + a[i]
// 时间复杂度 O(log n)
func (f fenwick) pre(i int) (res int) {
	for i++; i > 0; i &= i - 1 {
		res += f[i]
	}
	return
}

// 求区间和 a[l] + ... + a[r]
// 时间复杂度 O(log n)
func (f fenwick) query(l, r int) int {
	return f.pre(r) - f.pre(l-1)
}

func validSubarrays1(nums []int, k int, queries [][]int) []bool {
	n := len(nums)
	sum := make([]uint64, n+1)
	hash := map[int]uint64{}
	for i, x := range nums {
		// 把不同的 nums[i] 映射成一个随机的 uint64
		if _, ok := hash[x]; !ok {
			hash[x] = rand.Uint64()
		}
		sum[i+1] = sum[i] ^ hash[x]
	}

	// 离线询问：按照右端点分组
	type pair struct{ l, qid int }
	groups := make([][]pair, n)
	for i, q := range queries {
		groups[q[1]] = append(groups[q[1]], pair{q[0], i})
	}

	t := newFenwickTree(n)
	last := make(map[int]int, len(hash)) // 预分配空间
	ans := make([]bool, len(queries))
	for r, x := range nums {
		if i, ok := last[x]; ok {
			t.update(i, -1)
		}
		last[x] = r
		t.update(r, 1)
		for _, p := range groups[r] {
			ans[p.qid] = sum[r+1] == sum[p.l] && t.query(p.l, r) == k
		}
	}
	return ans
}

func validSubarrays(nums []int, k int, queries [][]int) []bool {
	n := len(nums)
	sum := make([]uint64, n+1)
	hash := map[int]uint64{}
	for i, x := range nums {
		// 把 nums[i] 映射成一个随机的 uint64
		if _, ok := hash[x]; !ok {
			hash[x] = rand.Uint64()
		}
		sum[i+1] = sum[i] ^ hash[x]
	}

	calcLeft := func(k int) []int {
		lefts := make([]int, n)
		cnt := map[int]int{}
		l := 0
		for i, x := range nums {
			cnt[x]++
			for len(cnt) >= k {
				v := nums[l]
				if cnt[v] > 1 {
					cnt[v]--
				} else {
					delete(cnt, v) // 保证 len(cnt) 是窗口内的不同元素个数
				}
				l++
			}
			lefts[i] = l
		}
		return lefts
	}

	l1 := calcLeft(k + 1)
	l2 := calcLeft(k)

	ans := make([]bool, len(queries))
	for i, p := range queries {
		l, r := p[0], p[1]
		ans[i] = sum[r+1] == sum[l] && l1[r] <= l && l < l2[r]
	}
	return ans
}
