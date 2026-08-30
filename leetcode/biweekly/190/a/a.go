package main

// https://space.bilibili.com/206214
func minBishopMoves(source, target []int) int {
	sx, sy := source[0], source[1]
	tx, ty := target[0], target[1]

	// 颜色不同
	if (sx+sy)%2 != (tx+ty)%2 {
		return -1
	}

	// 两点之间的直线，如果斜率是 -1 或者 1，那么两点可以直接到达，否则要走两步
	if sx+sy == tx+ty || sx-sy == tx-ty {
		return 1
	}
	return 2
}
