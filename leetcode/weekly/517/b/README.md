按题意要求模拟即可。

其中 $x^y$ 可以用**快速幂**计算，原理见[【图解】一张图秒懂快速幂](https://leetcode.cn/problems/powx-n/solution/tu-jie-yi-zhang-tu-miao-dong-kuai-su-mi-ykp3i/)。

注意取模。为什么可以在**中途取模**？原理见 [模运算的世界：当加减乘除遇上取模](https://leetcode.cn/circle/discuss/mDfnkW/)。

```py [sol-Python3]
class Solution:
    def sumDecoded(self, nums: list[int]) -> int:
        MOD = 1_000_000_007
        ans = 0
        for x in nums:
            d = x // 10

            # 计算 d 的十进制长度
            length_d = 0
            v = d
            while v > 0:
                length_d += 1
                v //= 10

            pow10 = 10 ** (length_d - x % 10)
            # 根据 pow10 求出 x = d//pow10 和 y = d%pow10
            ans += pow(d // pow10, d % pow10, MOD)
        return ans % MOD
```

```java [sol-Java]
class Solution {
    public static final int MOD = 1_000_000_007;

    public int sumDecoded(long[] nums) {
        long ans = 0;
        for (long x : nums) {
            long d = x / 10;
            // 计算 d 的十进制长度
            int lengthD = 0;
            for (long v = d; v > 0; v /= 10) {
                lengthD++;
            }
            long pow10 = (long) Math.pow(10, lengthD - x % 10);
            // 根据 pow10 求出 x = d/pow10 和 y = d%pow10
            ans += pow(d / pow10, d % pow10);
        }
        return (int) (ans % MOD);
    }

    private long pow(long x, long n) {
        long res = 1;
        for (; n > 0; n /= 2) {
            if (n % 2 > 0) {
                res = res * x % MOD;
            }
            x = x * x % MOD;
        }
        return res;
    }
}
```

```cpp [sol-C++]
class Solution {
    const int MOD = 1'000'000'007;

    long long qpow(long long x, long long n) {
        long long res = 1;
        for (; n; n /= 2) {
            if (n % 2) {
                res = res * x % MOD;
            }
            x = x * x % MOD;
        }
        return res;
    }

public:
    int sumDecoded(vector<long long>& nums) {
        long long ans = 0;
        for (long long x : nums) {
            long long d = x / 10;
            // 计算 d 的十进制长度
            int length_d = 0;
            for (long long v = d; v > 0; v /= 10) {
                length_d++;
            }
            long long pow10 = pow(10, length_d - x % 10);
            // 根据 pow10 求出 x = d/pow10 和 y = d%pow10
            ans += qpow(d / pow10, d % pow10);
        }
        return ans % MOD;
    }
};
```

```go [sol-Go]
const mod = 1_000_000_007

func pow(x, n int) int {
	res := 1
	for ; n > 0; n /= 2 {
		if n%2 > 0 {
			res = res * x % mod
		}
		x = x * x % mod
	}
	return res
}

func sumDecoded(nums []int64) (ans int) {
	for _, X := range nums {
		x := int(X)
		d := x / 10
		// 计算 d 的十进制长度
		lengthD := 0
		for v := d; v > 0; v /= 10 {
			lengthD++
		}
		pow10 := int(math.Pow10(lengthD - x%10))
		// 根据 pow10 求出 x = d/pow10 和 y = d%pow10
		ans += pow(d/pow10, d%pow10)
	}
	return ans % mod
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n\log U)$，其中 $n$ 是 $\textit{nums}$ 的长度，$U=\max(\textit{nums})$。
- 空间复杂度：$\mathcal{O}(1)$。

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
