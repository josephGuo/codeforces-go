本题强制在线（不允许使用离线算法）。

做法同 [4033. 有效 K 个不同元素子数组 I](https://leetcode.cn/problems/valid-k-unique-subarrays-i/)，请看 [我的题解](https://leetcode.cn/problems/valid-k-unique-subarrays-i/solutions/4016355/yi-huo-ha-xi-chi-xian-shu-zhuang-shu-zu-pkxwi/) 的方法一。

```py [sol-Python3]
class Solution:
    def validSubarrays(self, nums: list[int], k: int, l0: int, r0: int, q: int) -> list[bool]:
        n = len(nums)
        s = [0] * (n + 1)
        # 把 nums[i] 映射成随机数
        # 不要在 [0, 2**L) 内随机，见 https://codeforces.com/blog/entry/153335
        mp = defaultdict(lambda: randrange(10 ** 18))
        for i, x in enumerate(nums):
            s[i + 1] = s[i] ^ mp[x]

        def calc_left(k: int) -> list[int]:
            lefts = [0] * n
            cnt = defaultdict(int)
            l = 0
            for i, x in enumerate(nums):
                cnt[x] += 1
                while len(cnt) >= k:
                    v = nums[l]
                    if cnt[v] > 1:
                        cnt[v] -= 1
                    else:
                        del cnt[v]  # 保证 len(cnt) 是窗口内的不同元素个数
                    l += 1
                lefts[i] = l
            return lefts

        l1 = calc_left(k + 1)
        l2 = calc_left(k)

        ans = [False] * q
        l, r = l0, r0
        for i in range(q):
            if i > 0:
                g = l + r if ans[i - 1] else r - l
                l = (l ^ g) % n
                r = (r ^ g) % n
                if l > r:
                    l, r = r, l
            ans[i] = s[r + 1] == s[l] and l1[r] <= l < l2[r]
        return ans
```

```java [sol-Java]
class Solution {
    private static final Random random = new Random();

    public boolean[] validSubarrays(int[] nums, int k, int l0, int r0, int q) {
        int n = nums.length;
        long[] sum = new long[n + 1];
        Map<Integer, Long> hash = new HashMap<>();
        for (int i = 0; i < n; i++) {
            // 把 nums[i] 映射成一个随机的 long
            long randVal = hash.computeIfAbsent(nums[i], _ -> random.nextLong());
            sum[i + 1] = sum[i] ^ randVal;
        }

        int[] l1 = calcLeft(nums, k + 1);
        int[] l2 = calcLeft(nums, k);

        boolean[] ans = new boolean[q];
        int l = l0;
        int r = r0;
        for (int i = 0; i < q; i++) {
            if (i > 0) {
                int g = ans[i - 1] ? l + r : r - l;
                l = (l ^ g) % n;
                r = (r ^ g) % n;
                if (l > r) {
                    int tmp = l;
                    l = r;
                    r = tmp;
                }
            }
            ans[i] = sum[r + 1] == sum[l] && l1[r] <= l && l < l2[r];
        }
        return ans;
    }

    private int[] calcLeft(int[] nums, int k) {
        int n = nums.length;
        int[] lefts = new int[n];
        Map<Integer, Integer> cnt = new HashMap<>();
        int l = 0;
        for (int i = 0; i < n; i++) {
            cnt.merge(nums[i], 1, Integer::sum); // ++cnt[nums[i]]
            while (cnt.size() >= k) {
                int c = cnt.merge(nums[l], -1, Integer::sum); // c = --cnt[nums[l]]
                if (c == 0) {
                    cnt.remove(nums[l]); // 保证 cnt.size() 是窗口内的不同元素个数
                }
                l++;
            }
            lefts[i] = l;
        }
        return lefts;
    }
}
```

```cpp [sol-C++]
class Solution {
    static inline mt19937_64 rng = mt19937_64(chrono::steady_clock::now().time_since_epoch().count());

public:
    vector<bool> validSubarrays(vector<int>& nums, int k, int l0, int r0, int q) {
        int n = nums.size();
        vector<uint64_t> sum(n + 1);
        unordered_map<int, uint64_t> hash;
        for (int i = 0; i < n; i++) {
            int x = nums[i];
            // 把 nums[i] 映射成一个随机的 uint64_t
            if (!hash.contains(x)) {
                hash[x] = rng() / 3; // https://codeforces.com/blog/entry/153335
            }
            sum[i + 1] = sum[i] ^ hash[x];
        }

        auto calc_left = [&](int k) -> vector<int> {
            vector<int> lefts(n);
            unordered_map<int, int> cnt;
            int l = 0;
            for (int i = 0; i < n; i++) {
                cnt[nums[i]]++;
                while (cnt.size() >= k) {
                    auto it = cnt.find(nums[l]);
                    if (--it->second == 0) {
                        cnt.erase(it); // 保证 cnt.size() 是窗口内的不同元素个数
                    }
                    l++;
                }
                lefts[i] = l;
            }
            return lefts;
        };

        auto l1 = calc_left(k + 1);
        auto l2 = calc_left(k);

        vector<bool> ans(q);
        int l = l0, r = r0;
        for (int i = 0; i < q; i++) {
            if (i > 0) {
                int g = ans[i - 1] ? l + r : r - l;
                l = (l ^ g) % n;
                r = (r ^ g) % n;
                if (l > r) {
                    swap(l, r);
                }
            }
            ans[i] = sum[r + 1] == sum[l] && l1[r] <= l && l < l2[r];
        }
        return ans;
    }
};
```

```go [sol-Go]
func validSubarrays(nums []int, k int, l0 int, r0 int, q int) []bool {
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

	ans := make([]bool, q)
	l, r := l0, r0
	for i := range ans {
		if i > 0 {
			g := r - l
			if ans[i-1] {
				g = l + r
			}
			l = (l ^ g) % n
			r = (r ^ g) % n
			if l > r {
				l, r = r, l
			}
		}
		ans[i] = sum[r+1] == sum[l] && l1[r] <= l && l < l2[r]
	}
	return ans
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n+q)$，其中 $n$ 是 $\textit{nums}$ 的长度。对于滑动窗口，虽然写了个二重循环，但是内层循环中对 $\ell$ 加一的**总**执行次数不会超过 $n$ 次，所以二重循环的循环次数为 $\mathcal{O}(n)$。
- 空间复杂度：$\mathcal{O}(n)$。返回值不计入。

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
