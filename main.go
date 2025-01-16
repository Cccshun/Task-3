package main

import (
	"fmt"
	"sysu.com/task3/algo"
	"sysu.com/task3/im"
	"time"
)

func main() {
	im.Init("network/BA200.txt")
	startTime := time.Now()
	findSeed()
	elapsedTime := time.Since(startTime)
	fmt.Printf("运行时间: %s\n", elapsedTime)
}

func findSeed() {
	// evalType=1, 优化node攻击；evalType=2，优化edge攻击; evalType=else,优化node+edge攻击
	savePath := "data/BA200/tmp"
	//var ma = new(algo.MA)
	//ma.FindBestSeed(savePath, 3)
	//var ga = new(algo.GA)
	//ga.FindBestSeed(savePath, 3)
	var nma = new(algo.NMA)
	nma.FindBestSeed(savePath, 3)
}
