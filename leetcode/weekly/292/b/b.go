package main

import . "github.com/EndlessCheng/codeforces-go/leetcode/testutil"

// github.com/EndlessCheng/codeforces-go
func averageOfSubtree(root *TreeNode) (ans int) {
	// 返回 node 子树的节点值之和、节点个数
	var dfs func(*TreeNode) (int, int)
	dfs = func(node *TreeNode) (int, int) {
		if node == nil {
			return 0, 0
		}
		leftSum, leftSize := dfs(node.Left)    // 拆解问题：递归计算左子树的信息
		rightSum, rightSize := dfs(node.Right) // 拆解问题：递归计算右子树的信息
		sum := leftSum + rightSum + node.Val   // node 子树的节点值之和
		size := leftSize + rightSize + 1       // node 子树的节点个数
		if node.Val == sum/size {              // 题目要求下取整
			ans++
		}
		return sum, size
	}

	dfs(root)
	return
}
