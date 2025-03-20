package main

import (
	"math/rand"
	"sync"
	"sysu.com/task3/algo"
	"sysu.com/task3/im"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fileName := "BA200"
	im.Init("network/" + fileName + ".txt")
	findSeed(fileName)
}

func findSeed(fileName string) {
	// evalType=1, 优化node攻击；evalType=2，优化edge攻击; evalType=else,优化node+edge攻击
	evalType := 1
	var savePath string
	switch evalType {
	case 1:
		savePath = "data/node-edge/" + fileName + "-"
	case 2:
		savePath = "data/edge-node/" + fileName + "-"
	default:
		savePath = "data/" + fileName + "-"
	}
	wg := new(sync.WaitGroup)
	wg.Add(3)
	go new(algo.GA).FindBestSeed(savePath, evalType, wg)
	go new(algo.MA).FindBestSeed(savePath, evalType, wg)
	go new(algo.NMA).FindBestSeed(savePath, evalType, wg)
	wg.Wait()

}
