## 方法一：Dijkstra 算法

[Dijkstra 算法介绍](https://leetcode.cn/problems/network-delay-time/solution/liang-chong-dijkstra-xie-fa-fu-ti-dan-py-ooe8/)

在网格图最短路的基础上，额外增加两个参数 $k$ 和 $\textit{idx}$，定义 $\textit{dis}[k][i][j][\textit{idx}]$ 表示从起点到 $(i,j)$ 的最小路径代价，此时还剩下 $k$ 次转向机会，且方向为 $\textit{idx}$。这里 $\textit{idx}$ 是一个 $[0,3]$ 中的整数（分别对应左右上下四个方向）。

枚举左右上下四个方向：

- 如果前进方向与 $\textit{idx}$ 相同，那么 $k$ 不变。
- 否则，必须满足 $k>0$，然后把 $k$ 减少一。

下午两点 [B站@灵茶山艾府](https://space.bilibili.com/206214) 直播讲题，欢迎关注~

其他语言稍后添加。

```py [sol-Python3]
class Solution:
    def minCost(self, grid: list[list[int]], k0: int) -> int:
        dirs = ((0, -1), (0, 1), (-1, 0), (1, 0))  # 左右上下
        m, n = len(grid), len(grid[0])
        dis = [[[[inf] * 4 for _ in range(n)] for _ in range(m)] for _ in range(k0 + 1)]

        # 初始方向可以向右（1）或向下（3）
        h = [(grid[0][0], k0, 0, 0, 1), (grid[0][0], k0, 0, 0, 3)]
        while h:
            d, k, i, j, idx = heappop(h)
            if i == m - 1 and j == n - 1:
                return d
            if d > dis[k][i][j][idx]:
                continue
            for new_idx, (dx, dy) in enumerate(dirs):
                x, y = i + dx, j + dy
                if 0 <= x < m and 0 <= y < n:
                    new_k = k
                    if new_idx != idx:
                        if k == 0:
                            continue
                        new_k -= 1
                    new_d = d + grid[x][y]
                    if new_d < dis[new_k][x][y][new_idx]:
                        dis[new_k][x][y][new_idx] = new_d
                        heappush(h, (new_d, new_k, x, y, new_idx))
        return -1
```

```java [sol-Java]

```

```cpp [sol-C++]

```

```go [sol-Go]
var dirs = []struct{ x, y int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} // 左右上下

func minCost(grid [][]int, k0 int) int {
	m, n := len(grid), len(grid[0])
	dis := make([][][][4]int, k0+1)
	for k := range dis {
		dis[k] = make([][][4]int, m)
		for i := range dis[k] {
			dis[k][i] = make([][4]int, n)
			for j := range dis[k][i] {
				for idx := range dis[k][i][j] {
					dis[k][i][j][idx] = math.MaxInt / 2
				}
			}
		}
	}

	// 初始方向可以向右（1）或向下（3）
	h := hp{{grid[0][0], k0, 0, 0, 1}, {grid[0][0], k0, 0, 0, 3}}
	for len(h) > 0 {
		top := heap.Pop(&h).(tuple)
		d, k, i, j, idx := top.dis, top.k, top.i, top.j, top.idx
		if i == m-1 && j == n-1 {
			return d
		}
		if d > dis[k][i][j][idx] {
			continue
		}
		for newIdx, dir := range dirs {
			x, y := i+dir.x, j+dir.y
			if 0 <= x && x < m && 0 <= y && y < n {
				newK := k
				if newIdx != idx {
					if k == 0 {
						continue
					}
					newK--
				}
				newD := d + grid[x][y]
				if newD < dis[newK][x][y][newIdx] {
					dis[newK][x][y][newIdx] = newD
					heap.Push(&h, tuple{newD, newK, x, y, newIdx})
				}
			}
		}
	}
	return -1
}

type tuple struct{ dis, k, i, j, idx int }
type hp []tuple

func (h hp) Len() int           { return len(h) }
func (h hp) Less(i, j int) bool { return h[i].dis < h[j].dis }
func (h hp) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *hp) Push(v any)        { *h = append(*h, v.(tuple)) }
func (h *hp) Pop() (v any)      { a := *h; *h, v = a[:len(a)-1], a[len(a)-1]; return }
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(kmn\log(kmn))$，其中 $m$ 和 $n$ 分别是 $\textit{grid}$ 的行数和列数。
- 空间复杂度：$\mathcal{O}(kmn)$。

## 方法二：动态规划

由于直走不消耗 $k$，走回头路必定消耗 $k$，所以我们不会回到相同的状态，所以添加参数 $k$ 后的分层图是一个**有向无环图**（DAG）。

所以可以直接跑记忆化搜索，无需 Dijkstra。

```py [sol-Python3]
class Solution:
    def minCost(self, grid: list[list[int]], k0: int) -> int:
        dirs = ((0, -1), (0, 1), (-1, 0), (1, 0))  # 左右上下
        m, n = len(grid), len(grid[0])

        @cache
        def dfs(k: int, i: int, j: int, idx: int) -> int:
            if i == j == 0:
                return grid[0][0]

            res = inf
            for new_idx, (dx, dy) in enumerate(dirs):
                x, y = i + dx, j + dy
                if 0 <= x < m and 0 <= y < n:
                    new_k = k
                    if new_idx != idx:
                        if k == 0:
                            continue
                        new_k -= 1
                    res = min(res, dfs(new_k, x, y, new_idx))
            return res + grid[i][j]

        ans = min(dfs(k0, m - 1, n - 1, 0), dfs(k0, m - 1, n - 1, 2))
        return -1 if ans == inf else ans
```

```go [sol-Go]
var dirs = []struct{ x, y int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} // 左右上下

func minCost(grid [][]int, k0 int) int {
	m, n := len(grid), len(grid[0])
	memo := make([][][][4]int, k0+1)
	for i := range memo {
		memo[i] = make([][][4]int, m)
		for j := range memo[i] {
			memo[i][j] = make([][4]int, n)
			for p := range memo[i][j] {
				for q := range memo[i][j][p] {
					memo[i][j][p][q] = -1
				}
			}
		}
	}

	var dfs func(int, int, int, int) int
	dfs = func(k, i, j, idx int) int {
		if i == 0 && j == 0 {
			return grid[0][0]
		}

		p := &memo[k][i][j][idx]
		if *p != -1 {
			return *p
		}

		res := math.MaxInt / 2
		for newIdx, dir := range dirs {
			x, y := i+dir.x, j+dir.y
			if 0 <= x && x < m && 0 <= y && y < n {
				newK := k
				if newIdx != idx {
					if k == 0 {
						continue
					}
					newK--
				}
				res = min(res, dfs(newK, x, y, newIdx))
			}
		}
		res += grid[i][j]

		*p = res
		return res
	}

	ans := min(dfs(k0, m-1, n-1, 0), dfs(k0, m-1, n-1, 2))
	if ans < math.MaxInt/2 {
		return ans
	}
	return -1
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(kmn)$，其中 $m$ 和 $n$ 分别是 $\textit{grid}$ 的行数和列数。
- 空间复杂度：$\mathcal{O}(kmn)$。

## 专题训练

1. 图论题单的「**§3.1 单源最短路：Dijkstra 算法**」。
2. 动态规划题单的「**二、网格图 DP**」。

## 分类题单

[如何科学刷题？](https://leetcode.cn/discuss/post/3141566/ru-he-ke-xue-shua-ti-by-endlesscheng-q3yd/)

1. [滑动窗口与双指针（定长/不定长/单序列/双序列/三指针/分组循环）](https://leetcode.cn/discuss/post/3578981/ti-dan-hua-dong-chuang-kou-ding-chang-bu-rzz7/)
2. [二分算法（二分答案/最小化最大值/最大化最小值/第K小）](https://leetcode.cn/discuss/post/3579164/ti-dan-er-fen-suan-fa-er-fen-da-an-zui-x-3rqn/)
3. [单调栈（基础/矩形面积/贡献法/最小字典序）](https://leetcode.cn/discuss/post/3579480/ti-dan-dan-diao-zhan-ju-xing-xi-lie-zi-d-u4hk/)
4. [网格图（DFS/BFS/综合应用）](https://leetcode.cn/discuss/post/3580195/fen-xiang-gun-ti-dan-wang-ge-tu-dfsbfszo-l3pa/)
5. [位运算（基础/性质/拆位/试填/恒等式/思维）](https://leetcode.cn/discuss/post/3580371/fen-xiang-gun-ti-dan-wei-yun-suan-ji-chu-nth4/)
6. [图论算法（DFS/BFS/拓扑排序/基环树/最短路/最小生成树/网络流）](https://leetcode.cn/discuss/post/3581143/fen-xiang-gun-ti-dan-tu-lun-suan-fa-dfsb-qyux/)
7. [动态规划（入门/背包/划分/状态机/区间/状压/数位/数据结构优化/树形/博弈/概率期望）](https://leetcode.cn/discuss/post/3581838/fen-xiang-gun-ti-dan-dong-tai-gui-hua-ru-007o/)
8. [常用数据结构（前缀和/差分/栈/队列/堆/字典树/并查集/树状数组/线段树）](https://leetcode.cn/discuss/post/3583665/fen-xiang-gun-ti-dan-chang-yong-shu-ju-j-bvmv/)
9. [数学算法（数论/组合/概率期望/博弈/计算几何/随机算法）](https://leetcode.cn/discuss/post/3584388/fen-xiang-gun-ti-dan-shu-xue-suan-fa-shu-gcai/)
10. [贪心与思维（基本贪心策略/反悔/区间/字典序/数学/思维/脑筋急转弯/构造）](https://leetcode.cn/discuss/post/3091107/fen-xiang-gun-ti-dan-tan-xin-ji-ben-tan-k58yb/)
11. [链表、树与回溯（前后指针/快慢指针/DFS/BFS/直径/LCA）](https://leetcode.cn/discuss/post/3142882/fen-xiang-gun-ti-dan-lian-biao-er-cha-sh-6srp/)
12. [字符串（KMP/Z函数/Manacher/字符串哈希/AC自动机/后缀数组/子序列自动机）](https://leetcode.cn/discuss/post/3144832/fen-xiang-gun-ti-dan-zi-fu-chuan-kmpzhan-ugt4/)

[我的题解精选（已分类）](https://github.com/EndlessCheng/codeforces-go/blob/master/leetcode/SOLUTIONS.md)

欢迎关注 [B站@灵茶山艾府](https://space.bilibili.com/206214)
