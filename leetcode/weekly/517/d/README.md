推荐先完成更简单的 [4040. 构造子集和的最少操作次数 I](https://leetcode.cn/problems/minimum-operations-to-form-subset-sum-i/)，[我的题解](https://leetcode.cn/problems/minimum-operations-to-form-subset-sum-i/solution/fen-zu-bei-bao-pythonjavacgo-by-endlessc-o7ao/)。

由于 $x$ 先乘一次再除一次，$x$ 不变，所以乘法操作后面跟着除法操作是没有意义的，白白浪费操作次数。

所以对于同一个数，最优操作顺序为：

- 先执行若干次（可能零次）除法操作，再执行若干次（可能零次）乘法操作。

所以可以把 $x$ 变成

$$
x' = \left\lfloor\dfrac{x}{2^a}\right\rfloor 2^b
$$

其中 $a$ 是除法操作次数，$b$ 是乘法操作次数。相当于把 $x$ 先右移 $a$ 次，再左移 $b$ 次。

由于 $\left\lfloor\dfrac{x}{2^a}\right\rfloor$ 的不同结果只有 $\mathcal{O}(\log x)$ 个，所以只需考虑 $\mathcal{O}(\log x)$ 个不同的 $a$。

由于 $x'\le \textit{sum}$（超过 $\textit{sum}$ 就不能选了），所以有 $\mathcal{O}(\log \textit{sum})$ 个不同的 $b$。

所以只有 $\mathcal{O}(\log x\log \textit{sum})$ 个不同的 $(a,b)$，只有 $\mathcal{O}(\log x\log \textit{sum})$ 个不同的 $x'$。

设 $U=\max(\textit{nums})$，现在问题变成：

- 给你 $n$ 组物品，每组物品有 $\mathcal{O}(\log U\log \textit{sum})$ 个，每个物品的体积为 $\left\lfloor\dfrac{x}{2^a}\right\rfloor 2^b$，价值为 $a+b$。每组物品至多选一个。计算满足体积和恰好为 $\textit{sum}$ 的最小价值和。

这是标准的**分组背包问题**。

请先学习 [0-1 背包](https://www.bilibili.com/video/BV16Y411v7Y6/)（包括空间优化写法），写一些 0-1 背包题目，例如 [2915. 和为目标值的最长子序列的长度](https://leetcode.cn/problems/length-of-the-longest-subsequence-that-sums-to-target/)，[我的题解](https://leetcode.cn/problems/length-of-the-longest-subsequence-that-sums-to-target/solutions/2502839/mo-ban-qia-hao-zhuang-man-xing-0-1-bei-b-0nca/)。

分组背包只需在 0-1 背包的基础上，把「选或不选」改成「**枚举选这一组内的哪个物品，或者一个都不选**」。

[本题视频讲解](https://www.bilibili.com/video/BV1NV4X6bEhr/?t=6m5s)，欢迎点赞关注~

## 优化前

```py [sol-Python3]
class Solution:
    def minOperations(self, nums: list[int], sum: int) -> int:
        f = [0] + [inf] * sum

        for x in nums:
            for i in range(sum, 0, -1):
                # 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
                # 本题是分组背包，要枚举选哪个物品（枚举除法操作次数为 a，乘法操作次数为 b）
                a = 0
                while x >> a:
                    b = 0
                    while x >> a << b <= i:
                        # 物品体积为 x>>a<<b，价值为 a+b
                        f[i] = min(f[i], f[i - (x >> a << b)] + a + b)
                        b += 1
                    a += 1
 
            if f[sum] == 0:  # 优化：提前返回
                return 0

        return -1 if f[sum] == inf else f[sum]
```

```java [sol-Java]
class Solution {
    public int minOperations(int[] nums, int sum) {
        int[] f = new int[sum + 1];
        Arrays.fill(f, Integer.MAX_VALUE / 2); // 避免加法溢出
        f[0] = 0;

        for (int x : nums) {
            for (int i = sum; i > 0; i--) {
                // 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
                // 本题是分组背包，要枚举选哪个物品（枚举除法操作次数为 a，乘法操作次数为 b）
                for (int a = 0; (x >> a) > 0; a++) {
                    for (int b = 0; (x >> a << b) <= i; b++) {
                        // 物品体积为 x>>a<<b，价值为 a+b
                        f[i] = Math.min(f[i], f[i - (x >> a << b)] + a + b);
                    }
                }
            }
            if (f[sum] == 0) { // 优化：提前返回
                return 0;
            }
        }

        return f[sum] == Integer.MAX_VALUE / 2 ? -1 : f[sum];
    }
}
```

```cpp [sol-C++]
class Solution {
public:
    int minOperations(vector<int>& nums, int sum) {
        vector<int> f(sum + 1, INT_MAX / 2); // 避免加法溢出
        f[0] = 0;

        for (int x : nums) {
            for (int i = sum; i > 0; i--) {
                // 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
                // 本题是分组背包，要枚举选哪个物品（枚举除法操作次数为 a，乘法操作次数为 b）
                for (int a = 0; x >> a; a++) {
                    for (int b = 0; (x >> a << b) <= i; b++) {
                        // 物品体积为 x>>a<<b，价值为 a+b
                        f[i] = min(f[i], f[i - (x >> a << b)] + a + b);
                    }
                }
            }
            if (f[sum] == 0) { // 优化：提前返回
                return 0;
            }
        }

        return f[sum] == INT_MAX / 2 ? -1 : f[sum];
    }
};
```

```go [sol-Go]
func minOperations(nums []int, sum int) int {
	f := make([]int, sum+1)
	for i := 1; i <= sum; i++ {
		f[i] = math.MaxInt / 2 // 避免加法溢出
	}

	for _, x := range nums {
		for i := sum; i > 0; i-- {
			// 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
			// 本题是分组背包，要枚举选哪个物品（枚举除法操作次数为 a，乘法操作次数为 b）
			for a := 0; x>>a > 0; a++ {
				for b := 0; x>>a<<b <= i; b++ {
					// 物品体积为 x>>a<<b，价值为 a+b
					f[i] = min(f[i], f[i-x>>a<<b]+a+b)
				}
			}
		}
		if f[sum] == 0 { // 优化：提前返回
			return 0
		}
	}

	if f[sum] == math.MaxInt/2 {
		return -1
	}
	return f[sum]
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n\cdot \textit{sum}\cdot \log U\log \textit{sum})$，其中 $n$ 是 $\textit{nums}$ 的长度，$U=\max(\textit{nums})$。
- 空间复杂度：$\mathcal{O}(\textit{sum})$。

## 优化

对于 $x = \textit{nums}[i]$，$\left\lfloor\dfrac{x}{2^a}\right\rfloor 2^b$ 可能有重复值，例如 $x=6$ 先除以 $2$ 再乘以 $2$ 还是 $6$，所以 $a=b=1$ 与 $a=b=0$ 这两个物品体积相同。我们可以先生成所有不同的 $\left\lfloor\dfrac{x}{2^a}\right\rfloor 2^b$，如果有相同的 $\left\lfloor\dfrac{x}{2^a}\right\rfloor 2^b$，只保留 $a+b$ 小的。然后再枚举这一组物品，计算转移。

> **注**：实际上，按照代码的循环顺序，首次把 key 插入哈希表时的 value 就是最小的。

```py [sol-Python3]
class Solution:
    def minOperations(self, nums: list[int], sum: int) -> int:
        f = [0] + [inf] * sum

        for x in nums:
            # 生成这一组的所有物品，相同体积的物品，只保留价值最小的物品
            costs = defaultdict(lambda: inf)
            a = 0
            while x >> a:
                b = 0
                while x >> a << b <= sum:
                    v = x >> a << b
                    costs[v] = min(costs[v], a + b)
                    b += 1
                a += 1

            # 按照体积从小到大排序，方便跳出循环
            items = sorted(costs.items())

            for i in range(sum, 0, -1):
                for v, c in items:
                    if v > i:
                        break
                    f[i] = min(f[i], f[i - v] + c)  # 手写 min 更快，见【Python3 更快的写法】

            if f[sum] == 0:
                return 0

        return -1 if f[sum] == inf else f[sum]
```

```py [sol-Python3 更快的写法]
class Solution:
    def minOperations(self, nums: list[int], sum: int) -> int:
        f = [0] + [inf] * sum

        for x in nums:
            costs = {}
            a = 0
            while x >> a:
                b = 0
                while x >> a << b <= sum:
                    v = x >> a << b
                    if v not in costs:
                        costs[v] = a + b
                    b += 1
                a += 1

            # 按照体积从小到大排序，方便跳出循环
            items = sorted(costs.items())

            for i in range(sum, 0, -1):
                mn = f[i]  # 避免在循环中反复访问 list
                for v, c in items:
                    if v > i:
                        break
                    tmp = f[i - v] + c
                    if tmp < mn:  # 手写 min
                        mn = tmp
                f[i] = mn

            if f[sum] == 0:
                return 0

        return -1 if f[sum] == inf else f[sum]
```

## 专题训练

见下面动态规划题单的「**§3.4 分组背包**」。

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
