package main

import (
	"fmt"
	"sysu.com/task3/algo"
	"sysu.com/task3/im"
	"time"
)

func main() {
	im.Init("network/BA200.txt")
	savePath := "data/BA200/"

	startTime := time.Now()
	// evalType=1, 优化node攻击；evalType=2，优化edge攻击; evalType=else,优化node+edge攻击
	ma := &algo.Ma{}
	ma.FindSeed(savePath, 3)
	ga := &algo.Ga{}
	ga.FindSeed(savePath, 3)
	elapsedTime := time.Since(startTime)

	fmt.Printf("运行时间: %s\n", elapsedTime)
}
