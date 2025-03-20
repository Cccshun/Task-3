package main

import (
	"math/rand"
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
		savePath = "data/node-edge/" + fileName + "tmp"
	case 2:
		savePath = "data/edge-node/" + fileName + "tmp"
	default:
		savePath = "data/" + fileName + "tmp"
	}
	var ga = new(algo.GA)
	ga.FindBestSeed(savePath, evalType)
	var ma = new(algo.MA)
	ma.FindBestSeed(savePath, evalType)
	var nma = new(algo.NMA)
	nma.FindBestSeed(savePath, evalType)
}
