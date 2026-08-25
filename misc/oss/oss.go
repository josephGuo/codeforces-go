package oss

import (
	"fmt"
	"math"
	"math/bits"
	"slices"
	"strings"
	"unsafe"
)

/*
2026.6.13

《沉星之序》（Order of the Sinking Star）游戏原型：
《Heroes of Sokoban》https://www.puzzlescript.net/play.html?p=6860122
《Heroes of Sokoban II: Monsters》https://www.puzzlescript.net/play.html?p=6910207
《Heroes of Sokoban III: The Bard and The Druid》https://www.puzzlescript.net/play.html?p=7072276
《Mirror Isles》https://alan.draknek.org/games/puzzlescript/mirrors.php
《Skipping Stones To Lonely Homes》https://alan.draknek.org/games/puzzlescript/skipping-stones.php
《PROMESST》https://silverspaceship.com/promesst/
《PROMESST2》https://silverspaceship.com/promesst2/
《ENIGMASH》https://jacklance.github.io/PuzzleScript/play.html?p=cfdcc6e23f1fb3e9de2fd42fafaf4d4c

*/

type warriorArrType [warriorNumberInit]point
type thiefArrType [thiefNumberInit]point
type wizardArrType [wizardNumberInit]point
type priestArrType [priestNumberInit]point
type druidArrType [druidNumberInit]point
type bardArrType [bardNumberInit]point
type explorerArrType [explorerNumberInit]point
type sailorArrType [sailorNumberInit]point
type merchantArrType [merchantNumberInit]point

type stoneArrType [stoneNumberInit]point
type crystalArrType [crystalNumberInit + grassNumberInit]point
type grassArrType [(crystalNumberInit + grassNumberInit) * min(druidNumberInit, 1)]point
type skippingStoneArrType [skippingStoneNumberInit]point
type lilyArrType [lilyNumberInit]point
type goblinArrType [goblinNumberInit]pointWithDir
type dragonArrType [len(dragonDirInit)]pointWithDir
type beamArrType [len(beamDirInit)]pointWithDir
type mirrorArrType [len(mirrorDirInit) / 2]pointWithDir
type mirrorRefArrType [len(mirrorRefDirInit) / 2]pointWithDir
type mirrorAuxArrType [len(mirrorAuxDirInit) / 2]pointWithDir

type data struct {
	warrior  warriorArrType  // A 推多个对象
	thief    thiefArrType    // T 拉一个对象
	wizard   wizardArrType   // W 交换对象
	cleric   priestArrType   // C 自己以及上下左右无敌
	druid    druidArrType    // D 把对象变成石头
	bard     bardArrType     // B 同时移动切比雪夫距离 <= 2 的对象
	explorer explorerArrType // 7 普通角色，无法推对象
	sailor   sailorArrType   // 8 普通角色，推一个对象
	merchant merchantArrType // 9 普通角色，推一个对象

	// 石头   挡光，无法被反射
	stones stoneArrType // S

	// 水晶   透明，可以被反射
	crystals crystalArrType // s

	// 草
	grass grassArrType // w

	// 水漂石
	skippingStones skippingStoneArrType // k

	// 睡莲叶，z=-1
	lilies lilyArrType // l

	// 怪物
	goblins goblinArrType // g
	dragons dragonArrType // d

	// 镜子
	mirrors     mirrorArrType    // M
	mirrorRefs  mirrorRefArrType // R 主镜子 + 可以被反射
	mirrorAuxes mirrorAuxArrType // m 关卡名中称其为 mundane

	// 光束
	// 高 4 位是类型，低 4 位是方向
	beams beamArrType // b

	// 哪些角色变成了水晶
	isCrystalMask [min(druidNumberInit, 1)]uint8 // 不支持 9

	// todo 用于判断牧师是否漂浮
	//isPriestAttacked bool

	// 门的开闭，避免反复计算
	doorOpened        [doorKinds]bool
	monsterDoorOpened bool

	// 当前角色类型
	curCharTypeNum int8
}

const mapSizeH = int8(len(levelMap))

var mapSizeN, mapSizeM int8
var hasWater = false

func initMap() {
	fmt.Println("data 结构体大小:", unsafe.Sizeof(data{}))

	mapSizeN = int8(len(levelMap[0]))
	mapSizeM = int8(len(levelMap[0][0]))

	for i, p := range finals {
		finals[i] = changeNegPoint(p)
	}
	for i, p := range monsterDoors {
		monsterDoors[i] = changeNegPoint(p)
	}
	for _, ps := range doors {
		for i, p := range ps {
			ps[i].point = changeNegPoint(p.point)
		}
	}
	for _, ps := range weightSwitches {
		for i, p := range ps {
			ps[i] = changeNegPoint(p)
		}
	}

	var doorMask, switchMask int
	var warriorNum, thiefNum, wizardNum, priestNum, druidNum, bardNum, explorerNum, sailorNum, merchantNum int
	var stoneNum, crystalNum, skippingNum, grassesNum, lilyNum, beamNum, mirrorNum, mirrorRefNum, mirrorAuxNum int
	var goblinNum, dragonNum int

	for i, ps := range doors {
		if len(ps) > 0 {
			doorMask |= 1 << i
		}
	}
	for i, ps := range weightSwitches {
		if len(ps) > 0 {
			switchMask |= 1 << i
		}
	}

	if warriorPosInit != noPos {
		warriorNum++
	}
	if thiefPosInit != noPos {
		thiefNum++
	}
	if wizardPosInit != noPos {
		wizardNum++
	}
	if priestPosInit != noPos {
		priestNum++
	}
	if bardPosInit != noPos {
		bardNum++
	}
	if sailorPosInit != noPos {
		sailorNum++
	}

	checkGrid := func(grid []string) {
		for _, row := range grid {
			if len(row) != int(mapSizeM) {
				panic("行不等长")
			}
			for _, ch := range row {
				switch ch {
				case 'A':
					warriorNum++
				case 'T':
					thiefNum++
				case 'W':
					wizardNum++
				case 'C':
					priestNum++
				case 'D':
					druidNum++
				case 'B':
					bardNum++
				case '7':
					explorerNum++
				case '8':
					sailorNum++
				case '9':
					merchantNum++
				case 'S':
					stoneNum++
				case 's':
					crystalNum++
				case 'w':
					grassesNum++
				case 'k':
					skippingNum++
				case 'l':
					lilyNum++
				case 'g':
					goblinNum++
				case 'd':
					dragonNum++
				case 'b':
					beamNum++
				case 'M':
					mirrorNum++
				case 'R':
					mirrorRefNum++
				case 'm':
					mirrorAuxNum++
				case 'x', 'y', 'z', '{':
					switchMask |= 1 << (ch - 'x')
				case 'X', 'Y', 'Z', '[':
					doorMask |= 1 << (ch - 'X')
				case '~':
					hasWater = true
				}
			}
		}
	}
	if len(waterMap) > 0 {
		checkGrid(waterMap)
	}
	for _, grid := range levelMap {
		checkGrid(grid)
	}

	if bits.OnesCount(uint(doorMask)) != doorKinds {
		panic("没有修改 door kinds")
	}
	if bits.OnesCount(uint(switchMask)) != doorKinds {
		panic("压力开关错误，请检查地图")
	}

	if warriorNum != warriorNumberInit {
		panic("没有修改 warrior number")
	}
	if thiefNum != thiefNumberInit {
		panic("没有修改 thief number")
	}
	if wizardNum != wizardNumberInit {
		panic("没有修改 wizard number")
	}
	if priestNum != priestNumberInit {
		panic("没有修改 priest number")
	}
	if druidNum != druidNumberInit {
		panic("没有修改 druid number")
	}
	if bardNum != bardNumberInit {
		panic("没有修改 bard number")
	}
	if explorerNum != explorerNumberInit {
		panic("没有修改 explorer number")
	}
	if sailorNum != sailorNumberInit {
		panic("没有修改 sailor number")
	}
	if !allowCloneMan && merchantNum != merchantNumberInit {
		panic("没有修改 merchant number")
	}

	// 检查数组大小是否与 levelMap 匹配
	if stoneNum != stoneNumberInit {
		panic("没有修改 stone number")
	}
	if crystalNum != crystalNumberInit {
		panic("没有修改 crystal number")
	}
	if grassesNum != grassNumberInit {
		panic("没有修改 grass number")
	}
	if skippingNum != skippingStoneNumberInit {
		panic("没有修改 skipping stone number")
	}
	if lilyNum != lilyNumberInit {
		panic("没有修改 lily number")
	}
	if goblinNum != len(goblinArrType{}) {
		panic("没有修改 goblin number")
	}
	if dragonNum != len(dragonArrType{}) {
		panic("没有修改 dragon dir")
	}
	if beamNum != len(beamArrType{}) {
		panic("没有修改 beamDirInit")
	}
	if len(beamDirInit) != len(beamTypeInit) {
		panic("没有修改 beam type")
	}
	if mirrorNum != len(mirrorArrType{}) {
		panic("没有修改 mirror dir")
	}
	if mirrorRefNum != len(mirrorRefArrType{}) {
		panic("没有修改 mirror ref dir")
	}
	if mirrorAuxNum != len(mirrorAuxArrType{}) {
		panic("没有修改 mirror aux dir")
	}
}

func (d *data) areAllMonstersDied() bool {
	for _, p := range d.goblins {
		if p.point != noPos && p.dir&dirCrystalDelta == 0 { // 没有变成水晶
			return false
		}
	}
	for _, p := range d.dragons {
		if p.point != noPos && p.dir&dirCrystalDelta == 0 { // 没有变成水晶
			return false
		}
	}
	return true
}

// 可以用 bitset 优化
func (d *data) getAllCharPos(isBigMap bool) []point {
	if isBigMap {
		return nil
	}

	allChars := make([]point, 0,
		warriorNumberInit+
			thiefNumberInit+
			wizardNumberInit+
			priestNumberInit+
			druidNumberInit+
			bardNumberInit+
			explorerNumberInit+
			sailorNumberInit+
			merchantNumberInit,
	)
	if warriorNumberInit > 0 {
		for _, p := range d.warrior {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if thiefNumberInit > 0 {
		for _, p := range d.thief {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if wizardNumberInit > 0 {
		for _, p := range d.wizard {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if priestNumberInit > 0 {
		for _, p := range d.cleric {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if druidNumberInit > 0 {
		for _, p := range d.druid {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if bardNumberInit > 0 {
		for _, p := range d.bard {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if explorerNumberInit > 0 {
		for _, p := range d.explorer {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if sailorNumberInit > 0 {
		for _, p := range d.sailor {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	if merchantNumberInit > 0 {
		for _, p := range d.merchant {
			if p != noPos {
				allChars = append(allChars, p)
			}
		}
	}
	return allChars
}

func (d *data) getAllMovableObjPos(isBigMap bool) (all, chars, nonChars []point) {
	chars = d.getAllCharPos(isBigMap)
	all = chars
	if mirrorDirInit != "" {
		for _, p := range d.mirrors {
			if p.point != noPos {
				all = append(all, p.point)
			}
		}
	}
	if mirrorRefDirInit != "" {
		for _, p := range d.mirrorRefs {
			if p.point != noPos {
				all = append(all, p.point)
			}
		}
	}
	if mirrorAuxDirInit != "" {
		for _, p := range d.mirrorAuxes {
			if p.point != noPos {
				all = append(all, p.point)
			}
		}
	}
	for _, p := range d.stones {
		if p != noPos {
			all = append(all, p)
		}
	}
	for _, p := range d.crystals {
		if p != noPos {
			all = append(all, p)
		}
	}
	for _, p := range d.skippingStones {
		if p != noPos {
			all = append(all, p)
		}
	}
	if goblinNumberInit > 0 {
		for _, p := range d.goblins {
			if p.point != noPos {
				all = append(all, p.point)
			}
		}
	}
	if dragonDirInit != "" {
		for _, p := range d.dragons {
			if p.point != noPos {
				all = append(all, p.point)
			}
		}
	}
	if beamDirInit != "" {
		for _, p := range d.beams {
			if p.point != noPos {
				all = append(all, p.point)
			}
		}
	}
	return all, chars, all[len(chars):]
}

func inBound(p point) bool {
	return 0 <= p.x && p.x < mapSizeN &&
		0 <= p.y && p.y < mapSizeM &&
		p.z < mapSizeH
}

func (d *data) inAnyClosedDoors(p point) bool {
	if !d.monsterDoorOpened && slices.Contains(monsterDoors[:], p) { // 怪物门
		return true
	}
	for i, opened := range d.doorOpened {
		if !opened && pdContains(doors[i], p) { // 压力门
			return true
		}
	}
	return false
}

// p 不是固体（p 是空地，或者 p 是可移动对象）
func (d *data) isValidPos(p point) bool {
	x, y, z := p.x, p.y, p.z
	if !inBound(p) {
		return false
	}

	if z < 0 {
		if z != -1 {
			panic("invalid z")
		}
		// todo 水中的门
		if levelMap[0][x][y] != '~' {
			return false
		}
		return true
	}

	if levelMap[z][x][y] == '#' { // 墙
		return false
	}
	if slices.Contains(d.grass[:], p) { // 草
		return false
	}
	if d.inAnyClosedDoors(p) {
		return false
	}
	return true
}

// 返回 mask 表示在哪些类型的 beam 中
// endPoints（若 mask 有 endpoint）0 为发射器一端，1 为远端
// todo 多个 endPoints
type beamInfo struct {
	endPoints   [2]point
	endPointDir point
}

func (d *data) withinBeams(p point, allNonCharObjs []point) (typeMask uint16, beamNf beamInfo) {
	for _, beam := range d.beams {
		// beam.dir 高 4 位是类型，低 4 位是方向
		beamDir := directions6[beam.dir&0xf]

		const hasMirror = mirrorDirInit != "" || mirrorRefDirInit != "" || mirrorAuxDirInit != ""
		if !hasMirror {
			// 剪枝：先粗略判断是否在光束方向上（不考虑障碍）
			if beamDir.x != 0 {
				// 上下，必须同 y 同 z
				if beam.y != p.y || beam.z != p.z {
					continue
				}
				if beamDir.x > 0 != (beam.x < p.x) {
					continue
				}
			} else if beamDir.y != 0 {
				// 左右，必须同 x 同 z
				if beam.x != p.x || beam.z != p.z {
					continue
				}
				if beamDir.y > 0 != (beam.y < p.y) {
					continue
				}
			} else { // beamDir.z != 0
				// 高低，必须同 x 同 y
				if beam.x != p.x || beam.y != p.y {
					continue
				}
			}
		}

		cur := beam.point
		meetP := false
		for {
			old := cur
			cur = cur.add(beamDir)
			if cur == p {
				meetP = true
				if beam.dir>>4 != beamEndpoint {
					break
				}
				// 继续走，走到末端
			}
			// 特判：水晶可以穿透
			if slices.Contains(d.crystals[:], cur) {
				continue
			}
			// 出地图边界，或者遇到不可穿透对象（可穿透对象：墙、人、水晶）
			// 可以在地图边界加一圈 '#' 
			if !inBound(cur) || slices.Contains(allNonCharObjs, cur) || d.inAnyClosedDoors(cur) {
				cur = old
				break
			}
		}

		if meetP {
			beamType := beam.dir >> 4
			typeMask |= 1 << beamType
			if beamType == beamEndpoint {
				beamNf.endPoints[0] = beam.point.add(beamDir)
				beamNf.endPoints[1] = cur
				beamNf.endPointDir = beamDir
			}
		}
	}
	return
}

// 是否被牧师保护（或者自己是牧师）
func (d *data) isProtected(char point) bool {
	if priestNumberInit == 0 {
		return false
	}
	priest := d.cleric[:][0] // todo 多个牧师
	if char == priest {
		return true
	}
	if mapSizeH > 1 {
		return isNeighbor6(char, priest)
	}
	return isNeighbor4(char, priest)
}

func (d *data) isFallingIntoGround(p point) bool {
	return p.z > 0 && levelMap[p.z-1][p.x][p.y] == '.'
}

// 在水面上（z=0）且下面（z=-1）没有物品（或者门）的对象，落入水中
// todo 摧毁水中的镜子
func (d *data) isFallIntoWater(p point) bool {
	// todo z > 0 中途遇到障碍
	// todo 栏杆
	if !hasWater ||
		p.z != 0 ||
		p == noPos ||
		levelMap[0][p.x][p.y] != '~' {
		return false
	}
	downP := point{p.x, p.y, -1}
	// 水中的门
	for i, opened := range d.doorOpened {
		if !opened && pdContains(doors[i], downP) {
			return false
		}
	}
	// 水中的石头、水晶、水漂石或睡莲叶
	if len(d.stones) > 0 && slices.Contains(d.stones[:], downP) ||
		len(d.crystals) > 0 && slices.Contains(d.crystals[:], downP) ||
		len(d.dragons) > 0 && len(d.druid) > 0 && pdContains(d.dragons[:], downP) ||
		len(d.goblins) > 0 && len(d.druid) > 0 && pdContains(d.goblins[:], downP) ||
		skippingStoneNumberInit > 0 && slices.Contains(d.skippingStones[:], downP) ||
		lilyNumberInit > 0 && slices.Contains(d.lilies[:], downP) {
		return false
	}
	return true
}

func (d *data) isAttacked(p point, burnPos []point) bool {
	// 喷火龙
	if slices.Contains(burnPos, p) {
		return true
	}

	// 哥布林
	for _, g := range d.goblins {
		if g.point == noPos || g.dir&dirCrystalDelta > 0 { // 是石头
			continue
		}
		if mapSizeH > 1 {
			if isNeighbor6(g.point, p) { // todo 是这样吗？
				return true
			}
		} else {
			if isNeighbor4(g.point, p) {
				return true
			}
		}
	}

	return false
}

const (
	dieTypeNo = iota
	dieTypeCrushed
	dieTypeAttacked
	dieTypeDrown
)

func (d *data) getDieType(p point, burnPos []point, isChar bool) int {
	// 被门压死
	// todo 忽略向上的门（应该抬高角色）
	for i, opened := range d.doorOpened {
		if !opened && pdContains(doors[i], p) {
			return dieTypeCrushed
		}
	}

	// 被攻击的优先级更高
	if d.isAttacked(p, burnPos) {
		if isChar && (d.isProtected(p) || len(d.isCrystalMask) > 0 && d.isCrystalMask[:][0] > 0 &&
			d.isCrystalMask[:][0]>>(d.getCharType(p)-1)&1 > 0) {
			return dieTypeNo // 注：如果下面是空或者水，不会落下去
		}
		return dieTypeAttacked
	}

	// 淹死
	if d.isFallIntoWater(p) {
		// todo 如果自己是牧师且周围人被攻击，那么自己是悬浮的，不会淹死
		// todo 牧师不会落水？
		if priestNumberInit > 0 && isChar && p == d.cleric[:][0] {
			return dieTypeNo
		}
		return dieTypeDrown
	}

	return dieTypeNo
}

func (d *data) getCharType(p point) uint8 {
	switch {
	case len(d.warrior) > 0 && d.warrior[:][0] == p:
		return charWarrior
	case len(d.thief) > 0 && d.thief[:][0] == p:
		return charThief
	case len(d.wizard) > 0 && d.wizard[:][0] == p:
		return charWizard
	case len(d.cleric) > 0 && d.cleric[:][0] == p:
		return charCleric
	case len(d.druid) > 0 && d.druid[:][0] == p:
		return charDruid
	case len(d.bard) > 0 && d.bard[:][0] == p:
		return charBard
	case len(d.explorer) > 0 && d.explorer[:][0] == p:
		return charExplorer
	case len(d.sailor) > 0 && d.sailor[:][0] == p:
		return charSailor
	case len(d.merchant) > 0 && d.merchant[:][0] == p:
		return charMerchant
	default:
		return 0
	}
}

// 反射：从 mirror.point 出发，往 dir 方向走 step 步
func (d *data) reflectTo(mirror pointWithDir, dir point, step int, allMovableObjs []point) point {
	cur := mirror.point
	for k := range step {
		cur.x += dir.x
		cur.y += dir.y
		cur.z += dir.z
		// 遇到另一面主镜子
		if i := pdIndex(d.mirrors[:], cur); i >= 0 {
			if k == step-1 { // 按 X 反射
				return noPos // 最终反射到了镜子上，这不行
			}
			dir = d.mirrors[i].reflectToAnotherDir(dir)
			if dir == (point{}) {
				// 镜子背对我们
				if step == math.MaxInt { // 法师 todo 喷火龙
					return d.mirrors[i].point
				}
				return noPos
			}
			continue // 改变光路，继续反射
		}
		// 遇到另一面可以反射的镜子
		if i := pdIndex(d.mirrorRefs[:], cur); i >= 0 {
			if k == step-1 {
				return noPos // 最终反射到了镜子上
			}
			dir = d.mirrorRefs[i].reflectToAnotherDir(dir)
			if dir == (point{}) {
				// 镜子背对我们
				if step == math.MaxInt { // 法师 todo 喷火龙
					return d.mirrorRefs[i].point
				}
				return noPos
			}
			continue // 改变光路，继续反射
		}
		// 遇到另一面辅助镜子
		if i := pdIndex(d.mirrorAuxes[:], cur); i >= 0 {
			if k == step-1 {
				return noPos // 最终反射到了辅助镜子上
			}
			dir = d.mirrorAuxes[i].reflectToAnotherDir(dir)
			if dir == (point{}) {
				// 镜子背对我们
				if step == math.MaxInt { // 法师 todo 喷火龙
					return d.mirrorAuxes[i].point
				}
				return noPos
			}
			continue // 改变光路，继续反射
		}
		// 光路被（不可移动对象）挡住
		if !d.isValidPos(cur) {
			return noPos
		}
		// 光路被非镜子对象挡住
		if i := slices.Index(allMovableObjs, cur); i >= 0 {
			if step == math.MaxInt { // 法师
				return allMovableObjs[i]
			}
			return noPos
		}
	}
	// 按 X 反射
	return cur
}

func (d *data) changePos(oldP, newP point, newDir uint8) {
	// 人
	if warriorNumberInit > 0 {
		if i := slices.Index(d.warrior[:], oldP); i >= 0 {
			d.warrior[i] = newP
			return
		}
	}
	if thiefNumberInit > 0 {
		if i := slices.Index(d.thief[:], oldP); i >= 0 {
			d.thief[i] = newP
			return
		}
	}
	if wizardNumberInit > 0 {
		if i := slices.Index(d.wizard[:], oldP); i >= 0 {
			d.wizard[i] = newP
			return
		}
	}
	if priestNumberInit > 0 {
		if i := slices.Index(d.cleric[:], oldP); i >= 0 {
			d.cleric[i] = newP
			return
		}
	}
	if druidNumberInit > 0 {
		if i := slices.Index(d.druid[:], oldP); i >= 0 {
			d.druid[i] = newP
			return
		}
	}
	if bardNumberInit > 0 {
		if i := slices.Index(d.bard[:], oldP); i >= 0 {
			d.bard[i] = newP
			return
		}
	}
	if sailorNumberInit > 0 {
		if i := slices.Index(d.sailor[:], oldP); i >= 0 {
			d.sailor[i] = newP
			return
		}
	}
	if explorerNumberInit > 0 {
		if i := slices.Index(d.explorer[:], oldP); i >= 0 {
			d.explorer[i] = newP
			return
		}
	}
	if merchantNumberInit > 0 {
		if i := slices.Index(d.merchant[:], oldP); i >= 0 {
			d.merchant[i] = newP
			return
		}
	}

	// 物
	if mirrorDirInit != "" {
		if i := pdIndex(d.mirrors[:], oldP); i >= 0 {
			d.mirrors[i].point = newP
			if newDir != math.MaxUint8 {
				d.mirrors[i].dir = newDir
			}
			return
		}
	}
	if mirrorRefDirInit != "" {
		if i := pdIndex(d.mirrorRefs[:], oldP); i >= 0 {
			d.mirrorRefs[i].point = newP
			if newDir != math.MaxUint8 {
				d.mirrorRefs[i].dir = newDir
			}
			return
		}
	}
	if mirrorAuxDirInit != "" {
		if i := pdIndex(d.mirrorAuxes[:], oldP); i >= 0 {
			d.mirrorAuxes[i].point = newP
			if newDir != math.MaxUint8 {
				d.mirrorAuxes[i].dir = newDir
			}
			return
		}
	}

	if len(d.stones) > 0 {
		if i := slices.Index(d.stones[:], oldP); i >= 0 {
			d.stones[i] = newP
			return
		}
	}

	if len(d.crystals) > 0 {
		if i := slices.Index(d.crystals[:], oldP); i >= 0 {
			d.crystals[i] = newP
			return
		}
	}

	if skippingStoneNumberInit > 0 {
		if i := slices.Index(d.skippingStones[:], oldP); i >= 0 {
			d.skippingStones[i] = newP
			return
		}
	}

	if goblinNumberInit > 0 && canPushGoblin {
		if i := pdIndex(d.goblins[:], oldP); i >= 0 {
			d.goblins[i].point = newP
			return
		}
	}

	if dragonDirInit != "" && canPushDragon {
		if i := pdIndex(d.dragons[:], oldP); i >= 0 {
			d.dragons[i].point = newP
			if newDir != math.MaxUint8 {
				d.dragons[i].dir &^= 7
				d.dragons[i].dir |= newDir
			}
			return
		}
	}

	if beamDirInit != "" && canPushBeam {
		if i := pdIndex(d.beams[:], oldP); i >= 0 {
			d.beams[i].point = newP
			if newDir != math.MaxUint8 {
				d.beams[i].dir = newDir
			}
			return
		}
	}

	panic("没有发生修改，请检查代码")
}

func (d *data) getCurCharPos() (pos point) {
	switch d.curCharTypeNum {
	case charDefault:
		panic("代码有误，当前角色不能为 charDefault")
	case charWarrior:
		pos = d.warrior[:][0]
	case charThief:
		pos = d.thief[:][0]
	case charWizard:
		pos = d.wizard[:][0]
	case charCleric:
		pos = d.cleric[:][0]
	case charDruid:
		pos = d.druid[:][0]
	case charBard:
		pos = d.bard[:][0]
	case charExplorer:
		pos = d.explorer[:][0]
	case charSailor:
		pos = d.sailor[:][0]
	case charMerchant:
		pos = d.merchant[:][0]
	default:
		panic("未找到当前角色")
	}
	return
}

// 进入的时候切回 '8'，离开的时候才切换角色
func (newData *data) bigMapForceSwapChar(oldP, newP point) {
	if !isBigMap {
		return
	}

	isOldOutside := strings.ContainsRune("ATWCDB789", rune(levelMap[0][oldP.x][oldP.y]))
	isNewOutside := strings.ContainsRune("ATWCDB789", rune(levelMap[0][newP.x][newP.y]))

	if !isOldOutside && isNewOutside {
		// 从场景内部移到场景外部
		// 重置所有人的位置，除了 8
		if warriorNumberInit > 0 {
			newData.warrior[:][0] = noPos
		}
		if thiefNumberInit > 0 {
			newData.thief[:][0] = noPos
		}
		if wizardNumberInit > 0 {
			newData.wizard[:][0] = noPos
		}
		if priestNumberInit > 0 {
			newData.cleric[:][0] = noPos
		}
		if druidNumberInit > 0 {
			newData.druid[:][0] = noPos
		}
		if bardNumberInit > 0 {
			newData.bard[:][0] = noPos
		}
		if explorerNumberInit > 0 {
			newData.explorer[:][0] = noPos
		}
		//if sailorNumberInit > 0 {
		//	newData.sailor[:][0] = noPos
		//}
		if merchantNumberInit > 0 {
			newData.merchant[:][0] = noPos
		}
		newData.curCharTypeNum = charSailor
	} else if isOldOutside && !isNewOutside {
		// 从场景外部移到场景内部
		switch levelMap[0][oldP.x][oldP.y] {
		case 'A':
			newData.sailor[:][0] = noPos
			newData.warrior[:][0] = newP
			newData.curCharTypeNum = charWarrior
		case 'T':
			newData.sailor[:][0] = noPos
			newData.thief[:][0] = newP
			newData.curCharTypeNum = charThief
		case 'W':
			newData.sailor[:][0] = noPos
			newData.wizard[:][0] = newP
			newData.curCharTypeNum = charWizard
		case 'C':
			newData.sailor[:][0] = noPos
			newData.cleric[:][0] = newP
			newData.curCharTypeNum = charCleric
		case 'D':
			newData.sailor[:][0] = noPos
			newData.druid[:][0] = newP
			newData.curCharTypeNum = charDruid
		case 'B':
			newData.sailor[:][0] = noPos
			newData.bard[:][0] = newP
			newData.curCharTypeNum = charBard
		case '7':
			newData.sailor[:][0] = noPos
			newData.explorer[:][0] = newP
			newData.curCharTypeNum = charExplorer
		case '8':
			newData.sailor[:][0] = newP
			newData.curCharTypeNum = charSailor
		case '9':
			newData.sailor[:][0] = noPos
			newData.merchant[:][0] = newP
			newData.curCharTypeNum = charMerchant
		}
	}
}

func solveLevel() []string {
	initMap()

	warriorInitArr := warriorArrType{}
	for i := range warriorInitArr {
		warriorInitArr[i] = noPos
	}
	thiefInitArr := thiefArrType{}
	for i := range thiefInitArr {
		thiefInitArr[i] = noPos
	}
	wizardInitArr := wizardArrType{}
	for i := range wizardInitArr {
		wizardInitArr[i] = noPos
	}
	priestInitArr := priestArrType{}
	for i := range priestInitArr {
		priestInitArr[i] = noPos
	}
	druidInitArr := druidArrType{}
	for i := range druidInitArr {
		druidInitArr[i] = noPos
	}
	bardInitArr := bardArrType{}
	for i := range bardInitArr {
		bardInitArr[i] = noPos
	}
	explorerInitArr := explorerArrType{}
	for i := range explorerInitArr {
		explorerInitArr[i] = noPos
	}
	sailorInitArr := sailorArrType{}
	for i := range sailorInitArr {
		sailorInitArr[i] = noPos
	}
	merchantInitArr := merchantArrType{}
	for i := range merchantInitArr {
		merchantInitArr[i] = noPos
	}

	mirrorInitArr := mirrorArrType{}
	mirrorRefInitArr := mirrorRefArrType{}
	mirrorAuxInitArr := mirrorAuxArrType{}
	stoneInitArr := stoneArrType{}
	for i := range stoneInitArr {
		stoneInitArr[i] = noPos
	}
	crystalInitArr := crystalArrType{}
	for i := range crystalInitArr {
		crystalInitArr[i] = noPos
	}
	grassInitArr := grassArrType{}
	for i := range grassInitArr {
		grassInitArr[i] = noPos
	}
	skippingStoneInitArr := skippingStoneArrType{}
	for i := range skippingStoneInitArr {
		skippingStoneInitArr[i] = noPos
	}
	lilyInitArr := lilyArrType{}
	goblinInitArr := goblinArrType{}
	dragonInitArr := dragonArrType{}
	beamInitArr := beamArrType{}

	__curCharTypeNum := initCharTypeNum
	if isBigMap {
		__curCharTypeNum = charSailor
	}

	__warriors := warriorInitArr[:0]
	if warriorPosInit != noPos {
		__warriors = append(__warriors, warriorPosInit)
	}
	__thiefs := thiefInitArr[:0]
	if thiefPosInit != noPos {
		__thiefs = append(__thiefs, thiefPosInit)
	}
	__wizards := wizardInitArr[:0]
	if wizardPosInit != noPos {
		__wizards = append(__wizards, wizardPosInit)
	}
	__priests := priestInitArr[:0]
	if priestPosInit != noPos {
		__priests = append(__priests, priestPosInit)
	}
	__druids := druidInitArr[:0]
	__bards := bardInitArr[:0]
	if bardPosInit != noPos {
		__bards = append(__bards, bardPosInit)
	}
	__explorers := explorerInitArr[:0]
	__sailor := sailorPosInit
	__sailors := sailorInitArr[:0]
	if sailorPosInit != noPos {
		__sailors = append(__sailors, sailorPosInit)
	}
	__merchants := merchantInitArr[:0]

	__mirrors := mirrorInitArr[:0]
	__mirrorRefs := mirrorRefInitArr[:0]
	__mirrorAuxes := mirrorAuxInitArr[:0]
	__stones := stoneInitArr[:0]
	__crystals := crystalInitArr[:0]
	__grass := grassInitArr[:0]
	__skippingStones := skippingStoneInitArr[:0]
	__lilies := lilyInitArr[:0]
	__goblins := goblinInitArr[:0]
	__dragons := dragonInitArr[:0]
	__beams := beamInitArr[:0]

	parseGrid := func(z int, grid []string) {
		for x, row := range grid {
			for y, ch := range row {
				p := point{int8(x), int8(y), int8(z)}
				switch ch {
				case 'A':
					if isBigMap {
						if __sailor == noPos {
							__sailor = p
						}
					} else {
						if __curCharTypeNum < 0 {
							__curCharTypeNum = charWarrior
						}
						__warriors = append(__warriors, p)
					}
				case 'T':
					if isBigMap {
						if __sailor == noPos {
							__sailor = p
						}
					} else {
						if __curCharTypeNum < 0 {
							__curCharTypeNum = charThief
						}
						__thiefs = append(__thiefs, p)
					}
				case 'W':
					if isBigMap {
						if __sailor == noPos {
							__sailor = p
						}
					} else {
						if __curCharTypeNum < 0 {
							__curCharTypeNum = charWizard
						}
						__wizards = append(__wizards, p)
					}
				case 'C':
					if isBigMap {
						if __sailor == noPos {
							__sailor = p
						}
					} else {
						if __curCharTypeNum < 0 {
							__curCharTypeNum = charCleric
						}
						__priests = append(__priests, p)
					}
				case 'D':
					if isBigMap {
						if __sailor == noPos {
							__sailor = p
						}
					} else {
						if __curCharTypeNum < 0 {
							__curCharTypeNum = charDruid
						}
						__druids = append(__druids, p)
					}
				case 'B':
					if isBigMap {
						if __sailor == noPos {
							__sailor = p
						}
					} else {
						if __curCharTypeNum < 0 {
							__curCharTypeNum = charBard
						}
						__bards = append(__bards, p)
					}
				case '7':
					if isBigMap {
						if __sailor == noPos {
							__sailor = p
						}
					} else {
						if __curCharTypeNum < 0 {
							__curCharTypeNum = charExplorer
						}
						__explorers = append(__explorers, p)
					}
				case '8':
					if __curCharTypeNum < 0 {
						__curCharTypeNum = charSailor
					}
					if __sailor == noPos {
						__sailor = p
					}
					__sailors = append(__sailors, p)
				case '9':
					if __curCharTypeNum < 0 {
						__curCharTypeNum = charMerchant
					}
					__merchants = append(__merchants, p)
				case 'M':
					idx := len(__mirrors)
					dir0 := getDir(mirrorDirInit[idx*2])
					dir1 := getDir(mirrorDirInit[idx*2+1])
					if dir0 > dir1 {
						dir0, dir1 = dir1, dir0
					}
					// 低小，高大
					__mirrors = append(__mirrors, pointWithDir{p, dir1<<4 | dir0})
				case 'R':
					idx := len(__mirrorRefs)
					dir0 := getDir(mirrorRefDirInit[idx*2])
					dir1 := getDir(mirrorRefDirInit[idx*2+1])
					if dir0 > dir1 {
						dir0, dir1 = dir1, dir0
					}
					// 低小，高大
					__mirrorRefs = append(__mirrorRefs, pointWithDir{p, dir1<<4 | dir0})
				case 'm':
					idx := len(__mirrorAuxes)
					dir0 := getDir(mirrorAuxDirInit[idx*2])
					dir1 := getDir(mirrorAuxDirInit[idx*2+1])
					if dir0 > dir1 {
						dir0, dir1 = dir1, dir0
					}
					// 低小，高大
					__mirrorAuxes = append(__mirrorAuxes, pointWithDir{p, dir1<<4 | dir0})
				case 'S': // 大写，不透光
					__stones = append(__stones, p)
				case 's': // 小写，透光
					__crystals = append(__crystals, p)
				case 'w':
					__grass = append(__grass, p)
				case 'k':
					__skippingStones = append(__skippingStones, p)
				case 'l':
					p.z-- // todo 注意这里
					__lilies = append(__lilies, p)
				case 'g':
					__goblins = append(__goblins, pointWithDir{p, 0}) // todo 默认方向为 dirs[0]
				case 'd':
					idx := len(__dragons)
					__dragons = append(__dragons, pointWithDir{p, getDir(dragonDirInit[idx])})
				case 'b':
					idx := len(__beams)
					dir := getDir(beamDirInit[idx])
					tp := beamTypeInit[idx] - '0'
					__beams = append(__beams, pointWithDir{p, tp<<4 | dir})
				case 'x', 'y', 'z', '{':
					weightSwitches[ch-'x'] = append(weightSwitches[ch-'x'], p)
					if also := sameSwitches[ch]; also > 0 {
						weightSwitches[also-'x'] = append(weightSwitches[also-'x'], p)
					}
				case 'X', 'Y', 'Z', '[':
					dir := getDir('n') // 默认方向向下，除非手动设置 doorDirString
					if len(doors) < len(doorDirString) {
						dir = getDir(doorDirString[len(doors)])
					}
					doors[ch-'X'] = append(doors[ch-'X'], pointWithDir{p, dir})
				case 'N':
					monsterDoors = append(monsterDoors, p)
				case 'f':
					finals = append(finals, p)
				case '.', '#', '~':
					// ignore
				default:
					panic(fmt.Sprintf("不支持的符号 %c", ch))
				}
			}
		}
	}
	if len(waterMap) > 0 {
		parseGrid(-1, waterMap)
	}
	for z, grid := range levelMap {
		parseGrid(z, grid)
	}

	// 有时候会手动添加 finals 的初始值，总体不一定是有序的
	slices.SortFunc(finals, cmpPoint)

	validChars := []int8{}
	if warriorNumberInit > 0 {
		validChars = append(validChars, charWarrior)
	}
	if thiefNumberInit > 0 {
		validChars = append(validChars, charThief)
	}
	if wizardNumberInit > 0 {
		validChars = append(validChars, charWizard)
	}
	if priestNumberInit > 0 {
		validChars = append(validChars, charCleric)
	}
	if druidNumberInit > 0 {
		validChars = append(validChars, charDruid)
	}
	if bardNumberInit > 0 {
		validChars = append(validChars, charBard)
	}
	if explorerNumberInit > 0 {
		validChars = append(validChars, charExplorer)
	}
	if sailorNumberInit > 0 { // todo __sailor != noPos
		validChars = append(validChars, charSailor)
	}
	if merchantNumberInit > 0 {
		validChars = append(validChars, charMerchant)
	}

	if !slices.Contains(validChars, __curCharTypeNum) {
		panic(fmt.Sprint("请修改 initCharTypeNum"))
	}

	levelData := data{
		warrior:  warriorInitArr,
		thief:    thiefInitArr,
		wizard:   wizardInitArr,
		cleric:   priestInitArr,
		bard:     bardInitArr,
		druid:    druidInitArr,
		explorer: explorerInitArr,
		sailor:   sailorInitArr,
		merchant: merchantInitArr,

		stones:   stoneInitArr,
		crystals: crystalInitArr,
		grass:    grassInitArr,
		goblins:  goblinInitArr,
		dragons:  dragonInitArr,

		skippingStones: skippingStoneInitArr,
		lilies:         lilyInitArr,

		mirrors:     mirrorInitArr,
		mirrorRefs:  mirrorRefInitArr,
		mirrorAuxes: mirrorAuxInitArr,

		beams: beamInitArr,

		//doorOpened:        doorOpenedInit,
		monsterDoorOpened: monsterDoorOpenedInit,

		curCharTypeNum: __curCharTypeNum,
	}

	type pair struct {
		data
		info string
	}
	from := map[data]pair{} // 同时充当 vis 的功能
	queue := []data{}
	defer func() { fmt.Printf("// 搜索了 %d 个状态\n", len(from)) }()

	add := func(last, d data, info string) {
		//fmt.Println(d.warrior, d.thief, d.beams[0].point)

		allMovableObjs, _, _ := d.getAllMovableObjPos(isBigMap)

		// 先确定门的开闭
		for i, weightSwitch := range weightSwitches {
			opened := !doorOpenedInit[i]
			// 如果有一个开关没有被压住，那么 opened 为初始状态
			for _, p := range weightSwitch {
				pushed := slices.Contains(allMovableObjs, p) ||
					len(d.grass) > 0 && slices.Contains(d.grass[:], p) // 草也可以压住开关
				if !pushed {
					opened = !opened // 变回 doorOpenedInit[i]
					break
				}
			}
			d.doorOpened[i] = opened

			// 水晶被门压碎（水晶在门中，但门没有打开）
			if !opened {
				for j, p := range d.stones {
					if pdContains(doors[i], p) {
						if !canDestroyObj {
							return
						}
						d.stones[j] = noPos
					}
				}
				for j, p := range d.crystals {
					if pdContains(doors[i], p) {
						if !canDestroyObj {
							return
						}
						d.crystals[j] = noPos
					}
				}
			}
		}

		// 被喷火龙攻击到的位置
		var burnedPos []point
		if len(d.dragons) > 0 {
			for _, dra := range d.dragons {
				if dra.z < 0 || dra.dir&dirCrystalDelta > 0 { // 是石头
					continue
				}
				dir := directions4[dra.dir]
				cur := point{dra.x, dra.y, dra.z}
				for {
					cur.x += dir.x
					cur.y += dir.y
					cur.z += dir.z
					if !d.isValidPos(cur) {
						break
					}

					if len(d.mirrors) > 0 || len(d.mirrorRefs) > 0 || len(d.mirrorAuxes) > 0 {
						// 这里的逻辑和法师是一样的
						mir := noPosDir
						if i := pdIndex(d.mirrors[:], cur); i >= 0 && d.mirrors[i].canReflect(dir) {
							mir = d.mirrors[i]
						} else if i := pdIndex(d.mirrorRefs[:], cur); i >= 0 && d.mirrorRefs[i].canReflect(dir) {
							mir = d.mirrorRefs[i]
						} else if i := pdIndex(d.mirrorAuxes[:], cur); i >= 0 && d.mirrorAuxes[i].canReflect(dir) {
							mir = d.mirrorAuxes[i]
						}

						// 面对的是镜子的正面
						if mir.point != noPos {
							dir2 := mir.reflectToAnotherDir(dir)
							// 沿着光路搜索，找第一个可交换对象
							refP := d.reflectTo(mir, dir2, math.MaxInt, allMovableObjs)
							if refP != noPos {
								burnedPos = append(burnedPos, refP)
							}
							break
						}
					}

					if slices.Contains(allMovableObjs, cur) {
						burnedPos = append(burnedPos, cur)
						break
					}
				}
			}
		}

		// 对象下落到 z >= 0
		// todo 整合后面的落水逻辑
		if mapSizeH > 1 {
			for _, p := range allMovableObjs {
				if p.z <= 0 {
					continue
				}
				// 如果 p 是牧师或其邻居，且正被攻击，那么 p 不会下落
				if d.isProtected(p) && d.isAttacked(p, burnedPos) {
					continue
				}
				oldP := p
				fallCnt := 0
				p.z--
				for p.z >= 0 && d.isValidPos(p) && !slices.Contains(allMovableObjs, p) {
					fallCnt++
					p.z--
				}
				p.z++
				if oldP.z != p.z {
					if !allowFallIntoGround {
						return
					}
					info += strings.Repeat("W", fallCnt)
					d.changePos(oldP, p, math.MaxUint8)
				}
			}
		}

		// 先判断是否有角色死亡
		for _, char := range d.getAllCharPos(false) {
			if d.getDieType(char, burnedPos, true) != dieTypeNo {
				return
			}
		}

		dieType := dieTypeNo
		// 一开始，以及切换角色，都不结算怪物之间的攻击
		isSwitching := info[0] == 'c' || '1' <= info[0] && info[0] <= '9'
		if !isSwitching && !d.monsterDoorOpened {
			// 哥布林
			goblins := d.goblins
			if len(d.goblins) > 0 {
				for i, p := range d.goblins {
					if p.z < 0 {
						continue
					}
					if p.dir&dirCrystalDelta > 0 { // 是石头
						// 落水
						if d.isFallIntoWater(p.point) {
							if !allowFallIntoWater {
								return
							}
							info += "W"
							goblins[i].z = -1
						}
						continue
					}
					// todo 变成水晶的哥布林 + 水晶哥布林落水
					if tp := d.getDieType(p.point, burnedPos, false); tp != dieTypeNo {
						if !canDestroyObj {
							return
						}
						dieType = tp
						goblins[i] = noPosDir
					}
				}
				if len(goblins) > 1 {
					slices.SortFunc(goblins[:], cmpPointWithDir)
				}
			}

			// 喷火龙
			dragons := d.dragons
			if len(d.dragons) > 0 {
				for i, p := range d.dragons {
					if p.z < 0 {
						continue
					}
					if p.dir&dirCrystalDelta > 0 { // 是石头
						// 落水
						if d.isFallIntoWater(p.point) {
							if !allowFallIntoWater {
								return
							}
							info += "W"
							dragons[i].z = -1
						}
						continue
					}
					if tp := d.getDieType(p.point, burnedPos, false); tp != dieTypeNo {
						if !canDestroyObj {
							return
						}
						dieType = tp
						dragons[i] = noPosDir
					}
				}
				if len(dragons) > 1 {
					slices.SortFunc(dragons[:], cmpPointWithDir)
				}
			}

			if canDestroyObj {
				d.goblins = goblins
				d.dragons = dragons
				d.monsterDoorOpened = d.areAllMonstersDied()
			}
		}

		// todo 石头/镜子落入水中的镜子，水中的镜子会被摧毁

		// 镜子
		if len(d.mirrors) > 0 {
			mir := d.mirrors[:]
			for i, p := range mir {
				if d.isFallIntoWater(p.point) {
					if !allowFallIntoWater {
						return
					}
					info += "W"
					mir[i].z = -1
				}
			}
			if len(d.mirrors) > 1 {
				slices.SortFunc(mir, cmpPointWithDir)
			}
		}

		// 可以被反射的镜子
		if len(d.mirrorRefs) > 0 {
			mir := d.mirrorRefs[:]
			for i, p := range mir {
				if d.isFallIntoWater(p.point) {
					if !allowFallIntoWater {
						return
					}
					info += "W"
					mir[i].z = -1
				}
			}
			if len(d.mirrorRefs) > 1 {
				slices.SortFunc(mir, cmpPointWithDir)
			}
		}

		// 辅助镜子
		if len(d.mirrorAuxes) > 0 {
			mir := d.mirrorAuxes[:]
			for i, p := range mir {
				if d.isFallIntoWater(p.point) {
					if !allowFallIntoWater {
						return
					}
					info += "W"
					mir[i].z = -1
				}
			}
			if len(d.mirrorAuxes) > 1 {
				slices.SortFunc(mir, cmpPointWithDir)
			}
		}

		// 石头
		if len(d.stones) > 0 {
			sto := d.stones[:]
			for i, p := range sto {
				if d.isFallIntoWater(p) {
					if !allowFallIntoWater {
						return
					}
					info += "W"
					sto[i].z = -1
				}
			}
			if len(d.stones) > 1 {
				slices.SortFunc(sto, cmpPoint)
			}
		}

		// 水晶
		if len(d.crystals) > 0 {
			cry := d.crystals[:]
			for i, p := range cry {
				if d.isFallIntoWater(p) {
					if !allowFallIntoWater {
						return
					}
					info += "W"
					cry[i].z = -1
				}
			}
			if len(d.crystals) > 1 {
				slices.SortFunc(cry, cmpPoint)
			}
		}

		// 草
		if len(d.grass) > 1 {
			slices.SortFunc(d.grass[:], cmpPoint)
		}

		// 水漂石
		if len(d.skippingStones) > 0 {
			sto := d.skippingStones[:]
			for i, p := range sto {
				if d.isFallIntoWater(p) {
					if !allowFallIntoWater {
						return
					}
					info += "W"
					sto[i].z = -1
				}
			}
			if len(d.skippingStones) > 1 {
				slices.SortFunc(sto, cmpPoint)
			}
		}

		// 睡莲叶
		if len(d.lilies) > 1 {
			slices.SortFunc(d.lilies[:], cmpPoint)
		}

		// 光束
		if len(d.beams) > 1 {
			slices.SortFunc(d.beams[:], cmpPointWithDir)
		}

		// 人
		if len(d.warrior) > 1 {
			slices.SortFunc(d.warrior[:], cmpPoint)
		}
		if len(d.thief) > 1 {
			slices.SortFunc(d.thief[:], cmpPoint)
		}
		if len(d.wizard) > 1 {
			slices.SortFunc(d.wizard[:], cmpPoint)
		}
		if len(d.cleric) > 1 {
			slices.SortFunc(d.cleric[:], cmpPoint)
		}
		if len(d.druid) > 1 {
			slices.SortFunc(d.druid[:], cmpPoint)
		}
		if len(d.bard) > 1 {
			slices.SortFunc(d.bard[:], cmpPoint)
		}
		if len(d.explorer) > 1 {
			slices.SortFunc(d.explorer[:], cmpPoint)
		}
		if len(d.sailor) > 1 {
			slices.SortFunc(d.sailor[:], cmpPoint)
		}
		if len(d.merchant) > 1 {
			slices.SortFunc(d.merchant[:], cmpPoint)
		}

		if _, ok := from[d]; !ok {
			if dieType == dieTypeAttacked {
				info += "K" // 怪物攻击动画
			} else if dieType == dieTypeDrown {
				info += "W"
			}
			from[d] = pair{last, info}
			queue = append(queue, d)
		}
	}

	add(data{}, levelData, "c")

	for len(queue) > 0 {
		// 注意入队的时候修改了物品的位置（重力落下）
		d := queue[0]
		queue = queue[1:]

		allMovableObjs, allChars, allNonChars := d.getAllMovableObjPos(isBigMap)

		var pass bool
		if !targetIsClearAllMonsters {
			// 标准版：所有人都到达终点
			if isBigMap {
				p := d.getCurCharPos()
				pass = slices.Equal([]point{p}, finals)
			} else if len(d.isCrystalMask) == 0 || d.isCrystalMask[:][0] == 0 { // 不能有人是水晶
				if len(allChars) > 1 {
					slices.SortFunc(allChars, cmpPoint)
				}
				pass = slices.Equal(allChars, finals)
			}
		} else {
			// 简化版：怪物门开启（怪物都被杀）
			pass = d.monsterDoorOpened
		}
		if pass {
			// 生成操作序列
			path := []string{}
			for {
				var ok bool
				pre, ok := from[d]
				if !ok {
					panic("代码修改了 d，与存入的 d 不符")
				}
				d = pre.data
				if d == (data{}) { // 初始状态
					break
				}
				if pre.info != "IGNORE" {
					//fmt.Println(pre.thief[0], pre.cleric[0]) // DEBUG
					path = append(path, pre.info)
				}
			}
			slices.Reverse(path)
			return path
		}

		// todo 如果角色的头上有物品，物品会跟着移动（注意镜子的方向会变）    堆叠上限是多少？？
		// todo 即使人没有移动，切换方向也会改变头上物品（镜子、激光等）的方向
		// todo 多控时，如果下一个位置是没有石头的水，则一个角色无法移动（已在商人中实现）

		// todo 修改 changePos 的代码，添加一个参数 alsoMoveTop bool，
		//      使得当物品移动时，物品上方的物品（如果有）也跟着移动

		// 先考虑按 x 镜子反射对象，这样后面移动更流畅
		// 只要有一个镜子反射失败（红光），就直接 return
		doMirrors := func() {
			newData := d
			swapped := uint(0)
		nextMirror:
			for _, mirror := range append(d.mirrors[:], d.mirrorRefs[:]...) {
				// 找两个方向最近的可反射的对象
				cur0 := mirror.point
				cur1 := mirror.point
				dir0 := directions6[mirror.dir&0xf]
				dir1 := directions6[mirror.dir>>4]
				foundMirror := uint8(0)
				for step := 1; ; step++ {
					justFound := uint8(0) // 是否找到了非镜子对象
					// 检查方向 0
					if foundMirror&1 == 0 {
						cur0 = cur0.add(dir0)
						if !d.isValidPos(cur0) {
							continue nextMirror
						}
						// todo bitset
						if pdContains(d.mirrors[:], cur0) || pdContains(d.mirrorAuxes[:], cur0) {
							foundMirror |= 1
						} else if slices.Contains(allMovableObjs, cur0) {
							justFound |= 1
						}
					}
					// 检查方向 1
					if foundMirror>>1 == 0 {
						cur1 = cur1.add(dir1)
						if !d.isValidPos(cur1) {
							continue nextMirror
						}
						if pdContains(d.mirrors[:], cur1) || pdContains(d.mirrorAuxes[:], cur1) {
							foundMirror |= 2
						} else if slices.Contains(allMovableObjs, cur1) {
							justFound |= 2
						}
					}
					if foundMirror == 3 {
						return // 不能两方向最近都是镜子
					}
					if justFound == 3 {
						return // 不能反射位置都是对象
					}
					if justFound == 0 {
						continue // 都是空，继续找
					}

					oldP := cur0
					dir := dir1
					if justFound == 2 {
						oldP = cur1
						dir = dir0 // 往另一个方向反射
					}

					// 无法反射的石头，视作墙壁
					if slices.Contains(d.stones[:], oldP) {
						continue nextMirror
					}

					// 反射
					newP := d.reflectTo(mirror, dir, step, allMovableObjs)
					if newP == noPos {
						return // 反射失败
					}
					itemIdx := slices.Index(allMovableObjs, oldP)
					if swapped>>itemIdx&1 > 0 {
						// todo 所有对象的分身
						// 不能再分身了
						if slices.Contains(d.merchant[:], oldP) {
							if newData.merchant[:][0] != noPos {
								return
							}
							newData.merchant[:][0] = newP
						} else if slices.Contains(d.crystals[:], oldP) {
							//newData.crystals[0] = newP // todo
						} else {
							// todo 其他对象的分身
						}
					} else {
						swapped |= 1 << itemIdx
						if i := pdIndex(newData.dragons[:], oldP); i >= 0 {
							// 如果是 oldP 是喷火龙，则朝向会变
							newDir := mirror.reflectDragon(newData.dragons[i].dir)
							newData.dragons[i] = pointWithDir{newP, newDir}
						} else if i := pdIndex(newData.mirrorRefs[:], oldP); i >= 0 {
							// 如果是 oldP 是可被反射的镜子，则与 mir 垂直的镜子会前后翻转
							newDir := mirror.reflectMirrorRef(newData.mirrorRefs[i].dir)
							newData.mirrorRefs[i] = pointWithDir{newP, newDir}
						} else {
							newData.changePos(oldP, newP, math.MaxUint8)
						}
					}
					break
				}

				// 合二为一
				if allowMerge {
					// todo 这里恰有两人
					man := newData.merchant[:]
					if man[0] != noPos && man[0] == man[1] {
						man[0] = noPos
					}
				}
			}

			if swapped == 0 {
				return
			}

			add(d, newData, "x")
		}
		doMirrors()

		// 移动当前角色
		switch d.curCharTypeNum {
		case charWarrior:
			// 普通移动一步

			p0 := d.warrior[:][0] // todo 暂时支持一个人
			withinBeams, _ := d.withinBeams(p0, allNonChars)

			// 在墙里面，但不能穿透
			if withinBeams>>beamThrough&1 == 0 && inBound(p0) && levelMap[p0.z][p0.x][p0.y] == '#' {
				goto afterSwitch
			}

			for dIdx, dir := range directions4 {
				// 该方向有多少个连续的对象
				cnt := 0
				cur := p0.add(dir)
				for slices.Contains(allMovableObjs, cur) {
					cnt++
					cur = cur.add(dir)
				}
				// 前面是否有空地
				if !(withinBeams>>beamThrough&1 > 0 && inBound(cur) && levelMap[cur.z][cur.x][cur.y] == '#') &&
					!d.isValidPos(cur) {
					continue // 枚举另一个方向
				}

				newData := d
				for range cnt {
					// 倒着回来
					nxt := cur.sub(dir) // 这是个物品
					newData.changePos(nxt, cur, math.MaxUint8)

					// todo 多层
					oldTop := point{nxt.x, nxt.y, nxt.z + 1}
					// todo 喷火龙 / 镜子
					if slices.Contains(allMovableObjs, oldTop) {
						newTop := cur
						newTop.z++
						newData.changePos(oldTop, newTop, uint8(dIdx))
					}
					cur = nxt
				}

				newP := p0.add(dir)

				if mapSizeH > 1 {
					oldTop := point{p0.x, p0.y, p0.z + 1}
					// 如果原位置头上有喷火龙或者镜子，修改其位置和朝向
					if i := pdIndex(newData.dragons[:], oldTop); i >= 0 {
						newTop := newP
						newTop.z++
						if !d.isValidPos(newTop) || slices.Contains(allMovableObjs, newTop) {
							continue // todo 暂时禁止喷火龙落地 
						}
						// todo 如果喷火龙和人的方向不同呢？
						newData.dragons[i] = pointWithDir{newTop, uint8(dIdx)}
					} else if slices.Contains(allMovableObjs, oldTop) {
						newTop := newP
						newTop.z++
						newData.changePos(oldTop, newTop, uint8(dIdx))
					}
					// todo 镜子
				}

				newData.warrior[:][0] = newP // todo 暂时支持一个人
				newData.bigMapForceSwapChar(p0, newP)
				add(d, newData, dir4String[dIdx])
			}
		case charThief:
			// 普通移动一步
			p0 := d.thief[:][0] // todo
			withinBeams, _ := d.withinBeams(p0, allNonChars)

			// 在墙里面，但不能穿透
			if withinBeams>>beamThrough&1 == 0 && inBound(p0) && levelMap[p0.z][p0.x][p0.y] == '#' {
				goto afterSwitch
			}

			for dIdx, dir := range directions4 {
				newP := p0.add(dir)

				// 前面是否有空地
				if !(withinBeams>>beamThrough&1 > 0 && inBound(newP) && levelMap[newP.z][newP.x][newP.y] == '#') &&
					(!d.isValidPos(newP) || slices.Contains(allMovableObjs, newP)) {
					continue // 枚举另一个方向
				}

				newData := d
				back := point{p0.x - dir.x, p0.y - dir.y, p0.z - dir.z}
				if slices.Contains(allMovableObjs, back) {
					// 拉人/物 -> 当前位置
					newData.changePos(back, p0, math.MaxUint8)
				}
				newData.thief[:][0] = newP
				newData.bigMapForceSwapChar(p0, newP)
				add(d, newData, dir4String[dIdx])
			}
		case charWizard:
			p0 := d.wizard[:][0]
			withinBeams, beamNf := d.withinBeams(p0, allNonChars)

			// 在墙里面，但不能穿透
			if withinBeams>>beamThrough&1 == 0 && inBound(p0) && levelMap[p0.z][p0.x][p0.y] == '#' {
				goto afterSwitch
			}

		nextDir:
			for dIdx, dir := range directions4 {
				var newP point
				if withinBeams>>beamDouble&1 > 0 {
					// todo 先看走一步是不是物品，黄光的规则是这样的吗？
					//newP = point{p0.x + dir.x, p0.y + dir.y, p0.z + dir.z}
					//if slices.Contains(allMovableObjs, newP) {
					//	// 和对象交换位置
					//	newData := d
					//	newData.changePos(newP, p0, math.MaxUint8) // newP 换到 p0
					//	newData.wizard[:][0] = newP                // 法师换到 newP
					//	add(d, newData, dir4String[dIdx]+"P")      // swap
					//	continue
					//}

					// 如果在 double 光中，优先级更高，只能往该方向走两步
					// todo 多个 double 光的情况，要叠加
					// todo 绿光
					const multi = 2
					newP = point{p0.x + dir.x*multi, p0.y + dir.y*multi, p0.z + dir.z*multi}
					if !d.isValidPos(newP) {
						continue // 出界或者有障碍物（墙、草）
					}
					if slices.Contains(allMovableObjs, newP) {
						// 和对象交换位置
						newData := d
						newData.changePos(newP, p0, math.MaxUint8) // newP 换到 p0
						newData.wizard[:][0] = newP                // 法师换到 newP
						add(d, newData, dir4String[dIdx]+"P")      // swap
						continue
					}
					// 空地，直接走过去（逻辑在后面）
				} else if withinBeams>>beamEndpoint&1 > 0 && (dir == beamNf.endPointDir || dir == beamNf.endPointDir.rev()) {
					// 如果在 endpoint 光中，优先级更高，只能往该方向走到终点
					// 如果面朝发射器移动，移动到发射器前一格子
					// 如果背朝发射器移动，移动到末端格子
					// 如果移动到的格子不合法，则不能移动，否则移动过去 / 交换物品
					if dir != beamNf.endPointDir { // 方向相反
						newP = beamNf.endPoints[0]
					} else {
						newP = beamNf.endPoints[1]
					}
					if newP == p0 || !d.isValidPos(newP) {
						continue // 原地不动 or 出界或者有障碍物（墙、草）
					}
					if slices.Contains(allMovableObjs, newP) {
						// 和对象交换位置
						newData := d
						newData.changePos(newP, p0, math.MaxUint8) // newP 换到 p0
						newData.wizard[:][0] = newP                // 法师换到 newP
						add(d, newData, dir4String[dIdx]+"P")      // swap
						continue
					}
				} else {
					// dir 方向是否有可交换对象
					cur := p0
					for {
						cur = cur.add(dir)
						newP := cur
						if !d.isValidPos(newP) {
							break // 出界或者有障碍物
						}
						if !slices.Contains(allMovableObjs, newP) {
							continue // 空地
						}

						mir := noPosDir
						if i := pdIndex(d.mirrors[:], newP); i >= 0 && d.mirrors[i].canReflect(dir) {
							mir = d.mirrors[i]
						} else if i := pdIndex(d.mirrorRefs[:], newP); i >= 0 && d.mirrorRefs[i].canReflect(dir) {
							mir = d.mirrorRefs[i]
						} else if i := pdIndex(d.mirrorAuxes[:], newP); i >= 0 && d.mirrorAuxes[i].canReflect(dir) {
							mir = d.mirrorAuxes[i]
						}

						// 面对的是镜子的正面
						if mir.point != noPos {
							dir2 := mir.reflectToAnotherDir(dir)
							// 沿着光路搜索，找第一个可交换对象
							newP = d.reflectTo(mir, dir2, math.MaxInt, allMovableObjs)
							if newP == noPos {
								break // 镜子反射路径没有任何对象，只能普通移动一步
							}
						}

						// 和对象交换位置
						// 注：这里可能自己和自己交换
						newData := d
						newData.changePos(newP, p0, math.MaxUint8) // newP 换到 p0
						newData.wizard[:][0] = newP                // 法师换到 newP
						add(d, newData, dir4String[dIdx]+"P")      // swap
						continue nextDir
					}

					// 没有可交换对象，那就普通移动
					newP = p0.add(dir)
					if !d.isValidPos(newP) || slices.Contains(allMovableObjs, newP) {
						continue // 枚举另一个方向
					}
				}

				newData := d
				newData.wizard[:][0] = newP
				newData.bigMapForceSwapChar(p0, newP)
				add(d, newData, dir4String[dIdx]) // move
			}
		case charCleric:
			// 普通移动一步
			p0 := d.cleric[:][0]
			withinBeams, _ := d.withinBeams(p0, allNonChars)
			for dIdx, dir := range directions4 {
				newP := p0.add(dir)
				if !d.isValidPos(newP) || slices.Contains(allMovableObjs, newP) {
					continue // 枚举另一个方向
				}
				newData := d
				if allowAllPushItem && withinBeams>>beamStrong&1 > 0 {
					if i := slices.Index(allMovableObjs, newP); i >= 0 {
						// 可以推物品
						nxt2 := newP.add(dir)
						if !d.isValidPos(nxt2) || slices.Contains(allMovableObjs, nxt2) {
							continue // 枚举另一个方向
						}
						newData.changePos(newP, nxt2, math.MaxUint8)
					}
				}
				newData.cleric[:][0] = newP
				newData.bigMapForceSwapChar(p0, newP)
				add(d, newData, dir4String[dIdx])
			}
		case charBard:
			p0 := d.bard[:][0]
			items := []point{}
			if isBigMap {
				items = append(items, p0)
			}
			for _, p := range allMovableObjs {
				if chebyshevDis(p, p0) <= 2 {
					items = append(items, p)
				}
			}

			// 普通移动一步
			// 切比雪夫距离 <= 2 的物品（包括自己）都移动一步
			for dIdx, dir := range directions4 {
				x, y, z := p0.x+dir.x, p0.y+dir.y, p0.z+dir.z
				if !d.isValidPos(point{x, y, z}) {
					continue
				}
				if len(items) > 1 {
					slices.SortFunc(items, func(a, b point) int {
						if dir.x != 0 {
							return int(b.x*dir.x - a.x*dir.x)
						}
						return int(b.y*dir.y - a.y*dir.y)
					})
				}

				newData := d
				unmovedItems := []point{}
				movedItems := []point{}
				for _, oldP := range items {
					// item 往前移动一格
					newP := oldP.add(dir)
					if !d.isValidPos(newP) { // 无法移动
						unmovedItems = append(unmovedItems, oldP)
						continue
					}
					// 大地图物品不能出界 
					// todo 目前只实现了诗人的逻辑
					if isBigMap && oldP != p0 && strings.ContainsRune("ATWCDB789", rune(levelMap[0][newP.x][newP.y])) {
						unmovedItems = append(unmovedItems, oldP)
						continue
					}
					// 尝试移动
					if chebyshevDis(newP, p0) > 2 { // item 是力场最前面的点
						if slices.Contains(allMovableObjs, newP) { // 不能与力场外的对象碰撞
							unmovedItems = append(unmovedItems, oldP)
							continue
						}
					} else if slices.Contains(unmovedItems, newP) { // 力场后面的点，不能与前面移动失败的对象碰撞
						unmovedItems = append(unmovedItems, oldP)
						continue
					}
					movedItems = append(movedItems, oldP)
					newData.changePos(oldP, newP, math.MaxUint8)
				}

				if !slices.Contains(unmovedItems, p0) {
					if newData.bard[:][0] != (point{x, y, z}) {
						panic("诗人移动错误，代码有误")
					}

					// 特性：如果诗人脚下是物品，且该物品移动了，那么诗人可以再走一格
					// todo 对于物品叠物品的情况，也是同样的规则？
					if slices.Contains(movedItems, point{p0.x, p0.y, p0.z - 1}) {
						nxtP := point{x + dir.x, y + dir.y, z + dir.z}
						if d.isValidPos(nxtP) && !slices.Contains(unmovedItems, nxtP) {
							newData.bard[:][0] = nxtP
						}
						// todo （待确认）如果 z-2 也移动了，那么再再走一格
					}

					newData.bigMapForceSwapChar(p0, newData.bard[:][0])
					add(d, newData, dir4String[dIdx])
				}
			}
		case charDruid:
			p0 := d.druid[:][0]
			for dIdx, dir := range directions4 {
				newP := p0.add(dir)
				// 草变水晶
				if len(d.grass) > 0 {
					if i := slices.Index(d.grass[:], newP); i >= 0 {
						newData := d
						newData.crystals[:][0] = newData.grass[i] // 加个切片避免报错
						newData.grass[i] = noPos
						add(d, newData, dir4String[dIdx]+"C") // trans
						continue
					}
				}

				// 水晶变草
				if len(d.crystals) > 0 && druidCrystalToGrass {
					if i := slices.Index(d.crystals[:], newP); i >= 0 {
						newData := d
						newData.grass[:][0] = newData.crystals[i]
						newData.crystals[i] = noPos
						add(d, newData, dir4String[dIdx]+"C") // trans
						continue
					}
				}

				// 哥布林 <-> 水晶
				if len(d.goblins) > 0 && priestNumberInit > 0 && d.cleric[:][0] != noPos {
					if i := pdIndex(d.goblins[:], newP); i >= 0 {
						newData := d
						newData.goblins[i].dir ^= dirCrystalDelta
						add(d, newData, dir4String[dIdx]+"C") // trans
						continue
					}
				}

				// 喷火龙 <-> 水晶
				if len(d.dragons) > 0 {
					if i := pdIndex(d.dragons[:], newP); i >= 0 {
						newData := d
						newData.dragons[i].dir ^= dirCrystalDelta
						add(d, newData, dir4String[dIdx]+"C") // trans
						continue
					}
				}

				// 人 <-> 水晶
				// 目前只支持单人
				if druidTransMan {
					if tp := d.getCharType(newP); tp > 0 {
						newData := d
						newData.isCrystalMask[:][0] ^= 1 << (tp - 1)
						add(d, newData, dir4String[dIdx]+"C") // trans
						continue
					}
				}

				// 普通移动一步
				if !d.isValidPos(newP) || slices.Contains(allMovableObjs, newP) {
					continue
				}
				newData := d
				newData.druid[:][0] = newP
				newData.bigMapForceSwapChar(p0, newP)
				add(d, newData, dir4String[dIdx]) // move
			}
		case charExplorer:
			// 普通移动一步
			p0 := d.explorer[:][0]
			for dIdx, dir := range directions4 {
				newP := p0.add(dir)
				if !d.isValidPos(newP) {
					if mapSizeH > 1 {
						// 如果头上有喷火龙或者镜子，修改其朝向
						if i := pdIndex(d.dragons[:], point{p0.x, p0.y, p0.z + 1}); i >= 0 {
							newData := d
							newData.dragons[i].dir &^= 7
							newData.dragons[i].dir |= uint8(dIdx)
							add(d, newData, dir4String[dIdx])
						}
						// todo 镜子
					}
					continue // 枚举另一个方向
				}

				newData := d
				if allowExplorerPushItem {
					if i := slices.Index(allMovableObjs, newP); i >= 0 {
						// 推物品
						nxt2 := newP.add(dir)
						if !d.isValidPos(nxt2) || slices.Contains(allMovableObjs, nxt2) {
							continue // 枚举另一个方向
						}
						newData.changePos(newP, nxt2, math.MaxUint8)
					}
				} else if slices.Contains(allMovableObjs, newP) {
					continue
				}

				if mapSizeH > 1 {
					oldTop := point{p0.x, p0.y, p0.z + 1}
					// 如果原位置头上有喷火龙或者镜子，修改其位置和朝向
					if i := pdIndex(newData.dragons[:], oldTop); i >= 0 {
						newTop := newP
						newTop.z++
						if !d.isValidPos(newTop) || slices.Contains(allMovableObjs, newTop) {
							continue // todo 暂时禁止喷火龙落地 
						}
						// todo 如果喷火龙和人的方向不同呢？
						newData.dragons[i] = pointWithDir{newTop, uint8(dIdx)}
					} else if slices.Contains(allMovableObjs, oldTop) {
						newTop := newP
						newTop.z++
						newData.changePos(oldTop, newTop, uint8(dIdx))
					}
					// todo 镜子
				}

				newData.explorer[:][0] = newP
				newData.bigMapForceSwapChar(p0, newP)
				add(d, newData, dir4String[dIdx])
			}
		case charSailor:
			// 普通移动一步
			p0 := d.sailor[:][0]
			for dIdx, dir := range directions4 {
				newP := p0.add(dir)
				if !d.isValidPos(newP) {
					if mapSizeH > 1 {
						// 如果头上有喷火龙或者镜子，修改其朝向
						if i := pdIndex(d.dragons[:], point{p0.x, p0.y, p0.z + 1}); i >= 0 {
							newData := d
							newData.dragons[i].dir &^= 7
							newData.dragons[i].dir |= uint8(dIdx)
							add(d, newData, dir4String[dIdx])
						}
						// todo 镜子
					}
					continue // 枚举另一个方向
				}

				newData := d
				if allowAllPushItem {
					if i := slices.Index(allMovableObjs, newP); i >= 0 {
						// 推物品
						nxt2 := newP.add(dir)
						if !d.isValidPos(nxt2) || slices.Contains(allMovableObjs, nxt2) {
							continue // 枚举另一个方向
						}
						newData.changePos(newP, nxt2, math.MaxUint8)
					}
				} else if slices.Contains(allMovableObjs, newP) {
					continue
				}

				if mapSizeH > 1 {
					oldTop := point{p0.x, p0.y, p0.z + 1}
					// 如果原位置头上有喷火龙或者镜子，修改其位置和朝向
					if i := pdIndex(newData.dragons[:], oldTop); i >= 0 {
						newTop := newP
						newTop.z++
						if !d.isValidPos(newTop) || slices.Contains(allMovableObjs, newTop) {
							continue // todo 暂时禁止喷火龙落地 
						}
						newData.dragons[i] = pointWithDir{newTop, uint8(dIdx)}
					} else if slices.Contains(allMovableObjs, oldTop) {
						newTop := newP
						newTop.z++
						newData.changePos(oldTop, newTop, uint8(dIdx))
					}
					// todo 镜子
				}

				newData.sailor[:][0] = newP
				newData.bigMapForceSwapChar(p0, newP)
				add(d, newData, dir4String[dIdx])
			}
		case charMerchant:
			// 多控
			// 普通移动一步
			for dIdx, dir := range directions4 {
				newData := d
				oldMerchant := newData.merchant
				man := newData.merchant[:]
				if len(newData.merchant) > 1 {
					slices.SortFunc(man, func(a, b point) int {
						if dir.x != 0 {
							return int(b.x*dir.x - a.x*dir.x)
						}
						return int(b.y*dir.y - a.y*dir.y)
					})
				}

				unmovedMan := []point{}
				moved := false
				for manIdx, p0 := range man {
					if p0 == noPos {
						continue
					}
					nxt := p0.add(dir)
					// 无法移动（注意岸边也是无法移动的）
					if !d.isValidPos(nxt) || d.isFallIntoWater(nxt) || slices.Contains(unmovedMan, nxt) {
						unmovedMan = append(unmovedMan, p0)
						continue
					}
					// 如果前面是物品，则推动（能移动的人已经移动了）
					if !slices.Contains(oldMerchant[:], nxt) && slices.Contains(allMovableObjs, nxt) {
						nxt2 := nxt.add(dir)
						// 无法推动前面的物品
						if !d.isValidPos(nxt2) ||
							!slices.Contains(oldMerchant[:], nxt2) && slices.Contains(allMovableObjs, nxt2) ||
							slices.Contains(unmovedMan, nxt2) {
							unmovedMan = append(unmovedMan, p0)
							continue
						}
						newData.changePos(nxt, nxt2, math.MaxUint8)
					}
					moved = true
					man[manIdx] = nxt // 移走！
				}
				if !moved { // 没人动
					continue
				}
				add(d, newData, dir4String[dIdx])
			}
		case charDefault:
			panic("代码有误，当前角色不能为 charDefault")
		default:
			// 跳石
			//oriChar := d.curCharTypeNum - skippingStoneDelta
			//_ = oriChar

		}

	afterSwitch:
		// 换成其他人
		if !isBigMap {
			for _, char := range validChars {
				if char == d.curCharTypeNum {
					continue
				}
				// 不能是水晶
				if len(d.isCrystalMask) > 0 && d.isCrystalMask[:][0]>>(char-1)&1 > 0 {
					continue
				}
				newData := d
				newData.curCharTypeNum = char
				var info string
				if len(allChars) > 2 {
					info = digits[char : char+1]
				} else {
					info = "c"
				}
				if d.curCharTypeNum == charBard {
					info = "B" + info // 等一下再换人
				}
				add(d, newData, info)
			}
		}
	}

	// 无解
	return nil
}

const digits = "0123456789"

const (
	charDefault = iota // 仅占位，不使用
	charWarrior
	charThief
	charWizard
	charCleric
	charDruid
	charBard
	charExplorer
	charSailor   // 同大地图角色
	charMerchant // Trader
)

var charNumToName = [...]byte{
	charWarrior:  'A',
	charThief:    'T',
	charWizard:   'W',
	charCleric:   'C',
	charBard:     'B',
	charDruid:    'D',
	charExplorer: '7',
	charSailor:   '8',
	charMerchant: '9',
}

const (
	beamDefault  = iota
	beamOpen     // 红 1
	beamDouble   // 橙 2
	beamSmash    // 黄 3
	beamThrough  // 绿 4
	beamStrong   // 青 5
	beamUnknown  // 蓝 6 todo
	beamEndpoint // 紫 7
	beamActive   // 白 8
	beamModify   // 彩 9
)

// 跳石，无法操纵，只能原地等待
// 当跳石被推动后，额外进入该角色
// 当跳石停止移动后，换回原来的角色（用 skippingStoneDelta + 原来的角色编号表示跳石的情况）
//const skippingStoneDelta = 1 << 6

const dirCrystalDelta = 1 << 6
