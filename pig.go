/*
 * Copyright 2011 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package main

import (
	"fmt"
	"rand"
)

const (
	赢得数  = 100 // The winning 评分 in a game of Pig
	系列比赛 = 10 // The number of 竞赛 per series to simulate
)

type 评分 struct {
	选手, 对手, 总数 int
}

type 动作 func(动作前 评分) (动作后 评分, 回合结束 bool)

func 掷骰(这回 评分) (评分, bool) {
	结果 := rand.Intn(6) + 1 // [1, 6] 之间的随机数
	if 结果 == 1 {
		return 评分{这回.对手, 这回.选手, 0}, true
	}
	return 评分{这回.选手, 这回.对手, 结果 + 这回.总数}, false
}

func 逗留(这回 评分) (评分, bool) {
	return 评分{这回.对手, 这回.选手 + 这回.总数, 0}, true
}

type 战略 func(评分) 动作

func 限定逗留(限度 int) 战略 {
	return func(这回 评分) 动作 {
		if 这回.总数 >= 限度 {
			return 逗留
		}
		return 掷骰
	}
}

func 竞赛(战略A, 战略B 战略) int {
	所有战略 := []战略{战略A, 战略B}
	var 这回 评分
	var 回合结束 bool
	选手 := rand.Intn(2) // 随机选择第一选手
	for 这回.选手+这回.总数 < 赢得数 {
		动作 := 所有战略[选手](这回)
		if 动作 != 掷骰 && 动作 != 逗留 {
			panic(fmt.Sprintf("选手 %d is cheating", 选手))
		}
		这回, 回合结束 = 动作(这回)
		if 回合结束 {
			选手 = (选手 + 1) % 2 // 换竞赛者
		}
	}
	return 选手
}

func 循环赛(所有战略 []战略) ([]int, int) {
	取胜次数 := make([]int, len(所有战略))
	for A := 0; A < len(所有战略); A++ {
		for B := A + 1; B < len(所有战略); B++ {
			for 比赛 := 0; 比赛 < 系列比赛; 比赛++ {
				胜利者 := 竞赛(所有战略[A], 所有战略[B])
				if 胜利者 == 0 {
					取胜次数[A]++
				} else {
					取胜次数[B]++
				}
			}
		}
	}
	战略竞赛数 := 系列比赛 * (len(所有战略) - 1) // 不许比赛自己
	return 取胜次数, 战略竞赛数
}

func 比率(数组 ...int) string {
	总数 := 0
	for _, 数 := range 数组 {
		总数 += 数
	}
	串 := ""
	for _, 数 := range 数组 {
		if 串 != "" {
			串 += ", "
		}
		百分比 := 100 * float64(数) / float64(总数)
		串 += fmt.Sprintf("%d/%d (%0.1f%%)", 数, 总数, 百分比)
	}
	return 串
}

func main() {
	所有战略 := make([]战略, 赢得数)
	for 限度 := range 所有战略 {
		所有战略[限度] = 限定逗留(限度 + 1)
	}
	取胜次数, 竞赛总数 := 循环赛(所有战略)

	for 限度 := range 所有战略 {
		fmt.Printf("取胜次数,失败次数: 逗留限度为% 4d: %s\n",
			限度+1, 比率(取胜次数[限度], 竞赛总数-取胜次数[限度]))
	}
}