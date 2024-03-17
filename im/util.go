package im

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

// 去除个体中对重复基因
func RemoveDuplicateGene(src Seed) {
	hashTable := map[int]struct{}{}
	for _, val := range src.Nodes {
		hashTable[val] = struct{}{}
	}

	// 补充基因
	for len(hashTable) < SeedSize {
		hashTable[NewGene()] = struct{}{}
	}

	idx := 0
	for k := range hashTable {
		src.Nodes[idx] = k
		idx++
	}
}

// 深拷贝种子
func DeepCopySeed(src Seed) Seed {
	nodes := make([]int, SeedSize)
	for i := range nodes {
		nodes[i] = src.Nodes[i]
	}
	return Seed{Nodes: nodes, Fit: src.Fit}
}

// 深拷贝种群
func DeepCopyPop(src []Seed) []Seed {
	dis := make([]Seed, PopSize)
	for i := range dis {
		dis[i] = DeepCopySeed(src[i])
	}
	return dis
}

// 轮盘赌选择
func RouletteSelection(src []Seed) []Seed {
	sort.Sort(BySeed(src))
	totalFit := float64(0)
	for _, seed := range src {
		totalFit += seed.Fit
	}

	//精英选择,始终保留最优解
	dist := [PopSize]Seed{src[0]}
	selectedSeed := map[*Seed]bool{&src[0]: true}

	rand.Shuffle(len(src), func(i, j int) { src[i], src[j] = src[j], src[i] })
	for i := 1; i < PopSize; i++ {
		randomNumber := rand.Float64() * totalFit
		accumulatedFit := float64(0)
		for _, seed := range src {
			accumulatedFit += seed.Fit
			if accumulatedFit >= randomNumber && !selectedSeed[&seed] {
				dist[i] = DeepCopySeed(seed)
				selectedSeed[&seed] = true
				break
			}
		}
	}
	return dist[:]
}

// 选择node在网络G中2-hop领域的相邻节点
func Get2HopNodes(node int) map[int]struct{} {
	set := make(map[int]struct{})
	for adj1 := range AdjList[node] { //选择node的1-hop领域
		for adj2 := range AdjList[adj1] { //选择node的2-hop领域
			set[adj2] = struct{}{}
		}
	}
	return set
}

func GetAvgFit(seeds []int, f func([]int) float64) float64 {
	sumFit := 0.0
	for i := 0; i < RepeatTime; i++ {
		sumFit += f(seeds)
	}
	return sumFit / RepeatTime
}

func CreateDataPath(path, algoName string) *os.File {
	file, err := os.Create(path + algoName + "-" + strconv.Itoa(SeedSize) + ".txt")
	if err != nil {
		fmt.Printf("创建文件错误:{%s}\n", err)
	}
	return file
}

func SaveData(file *os.File, datas ...float64) {
	for idx, data := range datas {
		if _, err := file.WriteString(strconv.FormatFloat(data, 'f', 3, 64)); err != nil {
			fmt.Printf("写文件错误: {%s}\n", err)
		}
		if idx != len(datas)-1 {
			if _, err := file.WriteString(" "); err != nil {
				fmt.Printf("写文件错误: {%s}\n", err)
			}
		}
	}

	if _, err := file.WriteString("\n"); err != nil {
		fmt.Printf("写文件错误: {%s}\n", err)
	}
}
