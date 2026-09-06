package main

// https://space.bilibili.com/206214
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
