## 方法一：正序遍历 + 栈

首先，找出在 $t=0$ 时刻就合并的机器人：如果机器人 $i$ 满足 $i=n-1$ 或者 $\textit{position}[i+1] - \textit{position}[i] > \textit{distance}$，那么机器人 $i$ 是「最右侧机器人」。

对于这些「最右侧机器人」，速度慢的机器人会在一定时间后与其左侧更快的机器人合并。我们可以把遍历过的机器人的速度保存在一个栈中，如果当前机器人的速度比栈顶小，那么栈顶机器人就会在一定时间后与当前机器人合并，于是弹出栈顶。反复直到栈为空，或者栈顶机器人的速度 $\le$ 当前机器人的速度。

最后，从栈底到栈顶，机器人的速度是递增的（允许相邻相等），每个机器人都追不上其右边的机器人，所以相邻机器人的距离始终大于 $\textit{distance}$，无法合并。最终答案为栈的大小。

[本题视频讲解](https://www.bilibili.com/video/BV18sbp6xEeE/?t=13m20s)，欢迎点赞关注~

```py [sol-Python3]
class Solution:
    def countGroups(self, position: list[int], speed: list[int], distance: int) -> int:
        n = len(position)
        st = [-1]  # 哨兵
        for i, p in enumerate(position):
            # 找到这一组的最右侧机器人
            if i == n - 1 or position[i + 1] - p > distance:
                v = speed[i]
                while st[-1] > v:
                    st.pop()  # 栈顶机器人速度更快，一段时间后与机器人 i 合并
                st.append(v)
        return len(st) - 1  # 减去哨兵
```

```java [sol-Java]
class Solution {
    public int countGroups(int[] position, int[] speed, int distance) {
        // 更快的写法见【Java 写法二】
        Deque<Integer> st = new ArrayDeque<>();
        st.push(-1); // 哨兵
        for (int i = 0; i < position.length; i++) {
            // 找到这一组的最右侧机器人
            if (i == position.length - 1 || position[i + 1] - position[i] > distance) {
                int v = speed[i];
                while (!st.isEmpty() && st.peek() > v) {
                    st.pop(); // 栈顶机器人速度更快，一段时间后与机器人 i 合并
                }
                st.push(v);
            }
        }
        return st.size() - 1; // 减去哨兵
    }
}
```

```java [sol-Java 写法二]
class Solution {
    public int countGroups(int[] position, int[] speed, int distance) {
        int[] st = speed; // 原地栈
        int top = -1;
        for (int i = 0; i < position.length; i++) {
            // 找到这一组的最右侧机器人
            if (i == position.length - 1 || position[i + 1] - position[i] > distance) {
                int v = speed[i];
                while (top >= 0 && st[top] > v) {
                    top--; // 栈顶机器人速度更快，一段时间后与机器人 i 合并
                }
                st[++top] = v; // 入栈
            }
        }
        return top + 1;
    }
}
```

```cpp [sol-C++]
class Solution {
public:
    int countGroups(vector<int>& position, vector<int>& speed, int distance) {
        int n = position.size();
        stack<int> st;
        st.push(-1); // 哨兵
        for (int i = 0; i < n; i++) {
            // 找到这一组的最右侧机器人
            if (i == n - 1 || position[i + 1] - position[i] > distance) {
                int v = speed[i];
                while (st.top() > v) {
                    st.pop(); // 栈顶机器人速度更快，一段时间后与机器人 i 合并
                }
                st.push(v);
            }
        }
        return st.size() - 1; // 减去哨兵
    }
};
```

```go [sol-Go]
func countGroups(position, speed []int, distance int) int {
	st := speed[:0] // 原地栈
	for i, p := range position {
		// 找到这一组的最右侧机器人
		if i == len(position)-1 || position[i+1]-p > distance {
			v := speed[i]
			for len(st) > 0 && st[len(st)-1] > v {
				st = st[:len(st)-1] // 栈顶机器人速度更快，一段时间后与机器人 i 合并
			}
			st = append(st, v)
		}
	}
	return len(st)
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n)$，其中 $n$ 是 $\textit{position}$ 的长度。虽然我们写了个二重循环，但每个元素至多入栈出栈各一次，所以二重循环的**总**循环次数是 $\mathcal{O}(n)$ 的，所以时间复杂度是 $\mathcal{O}(n)$。
- 空间复杂度：$\mathcal{O}(n)$ 或 $\mathcal{O}(1)$。如果把输入的数组当作栈用，可以做到 $\mathcal{O}(1)$ 空间。

## 方法二：倒序遍历

速度小的机器人会把左边速度更大的机器人删除（合并）。

于是倒序遍历，维护遍历过的「最右侧机器人」速度的最小值 $\textit{mn}$。

- 如果当前机器人的速度比 $\textit{mn}$ 大，那它可以追上 $\textit{mn}$，不计入答案。
- 否则答案加一，更新 $\textit{mn}$ 为当前机器人的速度。

```py [sol-Python3]
class Solution:
    def countGroups(self, position: list[int], speed: list[int], distance: int) -> int:
        mn = speed[-1]
        ans = 1
        for i in range(len(speed) - 2, -1, -1):
            if speed[i] <= mn and position[i + 1] - position[i] > distance:
                mn = speed[i]
                ans += 1
        return ans
```

```java [sol-Java]
class Solution {
    public int countGroups(int[] position, int[] speed, int distance) {
        int n = speed.length;
        int mn = speed[n - 1];
        int ans = 1;
        for (int i = n - 2; i >= 0; i--) {
            if (speed[i] <= mn && position[i + 1] - position[i] > distance) {
                mn = speed[i];
                ans++;
            }
        }
        return ans;
    }
}
```

```cpp [sol-C++]
class Solution {
public:
    int countGroups(vector<int>& position, vector<int>& speed, int distance) {
        int n = speed.size();
        int mn = speed[n - 1];
        int ans = 1;
        for (int i = n - 2; i >= 0; i--) {
            if (speed[i] <= mn && position[i + 1] - position[i] > distance) {
                mn = speed[i];
                ans++;
            }
        }
        return ans;
    }
};
```

```go [sol-Go]
func countGroups(position, speed []int, distance int) int {
	n := len(speed)
	mn := speed[n-1]
	ans := 1
	for i := n - 2; i >= 0; i-- {
		if speed[i] <= mn && position[i+1]-position[i] > distance {
			mn = speed[i]
			ans++
		}
	}
	return ans
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n)$，其中 $n$ 是 $\textit{position}$ 的长度。
- 空间复杂度：$\mathcal{O}(1)$。

## 相似题目

- [853. 车队](https://leetcode.cn/problems/car-fleet/)
- [1776. 车队 II](https://leetcode.cn/problems/car-fleet-ii/)

## 专题训练

1. 双指针题单的「**六、分组循环**」。
2. 数据结构题单的「**§3.3 邻项消除**」。

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
