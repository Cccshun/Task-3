package main

import (
	"fmt"
	"sysu.com/task3/algo"
	"sysu.com/task3/im"
	"time"
)

func main() {
	im.Init("network/B_BA_200.txt")
	startTime := time.Now()

	//ma := &algo.Ma{}
	//ma.FindSeed()

	ga := &algo.Ga{}
	ga.FindSeed()

	elapsedTime := time.Since(startTime)
	fmt.Printf("运行时间: %s\n", elapsedTime)
	//im.Test()
}
