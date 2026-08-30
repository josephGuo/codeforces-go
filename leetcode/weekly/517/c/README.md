由于先乘再除没有意义，所以只需考虑只有乘法操作，或者只有除法操作。

我们可以把 $x$ 变成 $x' = x\cdot 2^a$ 或者 $x' = \left\lfloor\dfrac{x}{2^a}\right\rfloor$，其中 $a$ 是操作次数。

由于 $x\cdot 2^a\le \textit{sum}$（超过 $\textit{sum}$ 就不能选了），所以有 $\mathcal{O}(\log \textit{sum})$ 个不同的乘法操作次数。

由于 $\left\lfloor\dfrac{x}{2^a}\right\rfloor$ 的不同结果只有 $\mathcal{O}(\log x)$ 个，所以有 $\mathcal{O}(\log x)$ 个不同的除法操作次数。

所以只有 $\mathcal{O}(\log \textit{sum} + \log x)$ 个不同的操作，只有 $\mathcal{O}(\log \textit{sum} + \log x)$ 个不同的 $x'$。

设 $U=\max(\textit{nums})$，现在问题变成：

- 给你 $n$ 组物品，每组物品有 $\mathcal{O}(\log \textit{sum} + \log U)$ 个，每个物品的体积为 $x\cdot 2^a\le \textit{sum}$ 或 $\left\lfloor\dfrac{x}{2^a}\right\rfloor$，价值为 $a$。每组物品至多选一个。计算满足体积和恰好为 $\textit{sum}$ 的最小价值和。

这是标准的**分组背包问题**。

请先学习 [0-1 背包](https://www.bilibili.com/video/BV16Y411v7Y6/)（包括空间优化写法），写一些 0-1 背包题目，例如 [2915. 和为目标值的最长子序列的长度](https://leetcode.cn/problems/length-of-the-longest-subsequence-that-sums-to-target/)，[我的题解](https://leetcode.cn/problems/length-of-the-longest-subsequence-that-sums-to-target/solutions/2502839/mo-ban-qia-hao-zhuang-man-xing-0-1-bei-b-0nca/)。

分组背包只需在 0-1 背包的基础上，把「选或不选」改成「**枚举选这一组内的哪个物品，或者一个都不选**」。

[本题视频讲解](https://www.bilibili.com/video/BV1NV4X6bEhr/?t=6m5s)，欢迎点赞关注~

```py [sol-Python3]
class Solution:
    def minOperations(self, nums: list[int], sum: int) -> int:
        f = [0] + [inf] * sum

        for x in nums:
            w = x.bit_length()  # x 的二进制长度
            for i in range(sum, 0, -1):
                # 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
                # 本题是分组背包，要枚举选哪个物品（枚举乘了 a 次或者除了 a 次）
                a = 0
                while x << a <= i:
                    # 物品体积为 x<<a，价值为 a
                    f[i] = min(f[i], f[i - (x << a)] + a)
                    a += 1

                # 从小到大枚举 x>>a，方便在 x>>a > i 时跳出循环
                a = w - 1
                while a > 0 and x >> a <= i:
                    # 物品体积为 x>>a，价值为 a
                    f[i] = min(f[i], f[i - (x >> a)] + a)
                    a -= 1

        return -1 if f[sum] == inf else f[sum]
```

```py [sol-Python3 记忆化搜索]
class Solution:
    def minOperations(self, nums: list[int], sum: int) -> int:
        @cache
        def dfs(i: int, s: int) -> int:
            if s == 0:  # 无需操作
                return 0
            if i < 0:  # 不足 sum，不合法
                return inf

            # 这一组一个也不选
            res = dfs(i - 1, s)

            x = nums[i]
            a = 0
            while x << a <= s:
                # 物品体积为 x<<a，价值为 a
                res = min(res, dfs(i - 1, s - (x << a)) + a)
                a += 1

            # 从小到大枚举 x>>a，方便在 x>>a > i 时跳出循环
            a = x.bit_length() - 1
            while a > 0 and x >> a <= s:
                # 物品体积为 x>>a，价值为 a
                res = min(res, dfs(i - 1, s - (x >> a)) + a)
                a -= 1

            return res

        ans = dfs(len(nums) - 1, sum)
        return -1 if ans == inf else ans
```

```java [sol-Java]
class Solution {
    public int minOperations(int[] nums, int sum) {
        int[] f = new int[sum + 1];
        Arrays.fill(f, Integer.MAX_VALUE / 2); // 避免加法溢出
        f[0] = 0;

        for (int x : nums) {
            int w = 32 - Integer.numberOfLeadingZeros(x); // x 的二进制长度
            for (int i = sum; i > 0; i--) {
                // 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
                // 本题是分组背包，要枚举选哪个物品（枚举乘了 a 次或者除了 a 次）
                for (int a = 0; (x << a) <= i; a++) {
                    // 物品体积为 x<<a，价值为 a
                    f[i] = Math.min(f[i], f[i - (x << a)] + a);
                }
                // 从小到大枚举 x>>a，方便在 x>>a > i 时跳出循环
                for (int a = w - 1; a > 0 && (x >> a) <= i; a--) {
                    // 物品体积为 x>>a，价值为 a
                    f[i] = Math.min(f[i], f[i - (x >> a)] + a);
                }
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
        vector f(sum + 1, INT_MAX / 2); // 避免加法溢出
        f[0] = 0;

        for (int x : nums) {
            int w = bit_width(1u * x); // x 的二进制长度
            for (int i = sum; i > 0; i--) {
                // 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
                // 本题是分组背包，要枚举选哪个物品（枚举乘了 a 次或者除了 a 次）
                for (int a = 0; (x << a) <= i; a++) {
                    // 物品体积为 x<<a，价值为 a
                    f[i] = min(f[i], f[i - (x << a)] + a);
                }
                // 从小到大枚举 x>>a，方便在 x>>a > i 时跳出循环
                for (int a = w - 1; a > 0 && (x >> a) <= i; a--) {
                    // 物品体积为 x>>a，价值为 a
                    f[i] = min(f[i], f[i - (x >> a)] + a);
                }
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
		w := bits.Len(uint(x)) // x 的二进制长度
		for i := sum; i > 0; i-- {
			// 回想一下，0-1 背包是选或不选，状态转移方程为 f[i] = min(f[i], f[i-物品体积] + 物品价值)
			// 本题是分组背包，要枚举选哪个物品（枚举乘了 a 次或者除了 a 次）
			for a := 0; x<<a <= i; a++ {
				// 物品体积为 x<<a，价值为 a
				f[i] = min(f[i], f[i-x<<a]+a)
			}
			// 从小到大枚举 x>>a，方便在 x>>a > i 时跳出循环
			for a := w - 1; a > 0 && x>>a <= i; a-- {
				// 物品体积为 x>>a，价值为 a
				f[i] = min(f[i], f[i-x>>a]+a)
			}
		}
	}

	if f[sum] == math.MaxInt/2 {
		return -1
	}
	return f[sum]
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n\cdot \textit{sum}\cdot (\log \textit{sum}+\log U))$，其中 $n$ 是 $\textit{nums}$ 的长度，$U=\max(\textit{nums})$。
- 空间复杂度：$\mathcal{O}(\textit{sum})$。

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
