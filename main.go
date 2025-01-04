package main

import (
	"fmt"
	"sysu.com/task3/algo"
	"sysu.com/task3/im"
	"time"
)

func main() {
	im.Init("network/BA200.txt")
	savePath := "data/BA200/tmp"

	startTime := time.Now()
	// evalType=1, 优化node攻击；evalType=2，优化edge攻击; evalType=else,优化node+edge攻击
	var ma = new(algo.Ma)
	ma.FindSeed(savePath, 3)
	var ga = new(algo.Ga)
	ga.FindSeed(savePath, 3)
	elapsedTime := time.Since(startTime)

	fmt.Printf("运行时间: %s\n", elapsedTime)
}
