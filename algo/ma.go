package algo

import (
	"fmt"
	"math/rand"
	"sync"

	"sysu.com/task3/im"
)

type MA struct {
	GA
}

// 局部搜索
func (m *MA) LocalSearch(evalType int) {
	for i := 0; i < im.PopSize; i++ {
		if rand.Float32() < im.PL {
			m.wg.Add(1)
			go m.doLocalSearch(&m.NewPop[i], evalType)
		}
	}
	m.wg.Wait()

	fmt.Printf("局部搜索后:%d   ", similarity(m.Pop, m.NewPop))
}

func (m *MA) doLocalSearch(seed *im.Seed, evalType int) {
	defer m.wg.Done()
	for i := 0; i < im.SeedSize; i++ {
		CompareAndSwap(seed, i, evalType)
	}
}

// 搜索2-hop内最优种子
func CompareAndSwap(seed *im.Seed, idx int, evalType int) {
	nodes := im.Get2HopNodes(seed.Nodes[idx])
	var wg sync.WaitGroup
	var mu sync.Mutex

	//用nodes中的节点替换seed.Nodes[idx],并评估适应度
	seedsMap := make(map[*im.Seed]float64)
	for node := range nodes {
		seedCopy := im.DeepCopySeed(*seed)
		seedCopy.Nodes[idx] = node
		im.RemoveDuplicateGene(seedCopy)
		wg.Add(1)
		go im.EvaluateSeedAsync(seedCopy, seedsMap, &wg, &mu, evalType)
	}
	wg.Wait()

	//找出最优解
	for key, val := range seedsMap {
		if val > seed.Fit {
			seed.Nodes = key.Nodes
			seed.Fit = val
		}
	}
}

func (m *MA) FindBestSeed(savePath string, evalType int) {
	m.Init(evalType)
	file := im.CreateDataPath(savePath, "ma")
	defer file.Close()
	for i := 0; i < im.MaxGen; i++ {
		m.Crossover(evalType)
		m.Mutate(evalType)
		m.LocalSearch(evalType)
		m.Select()
		if evalType == 1 {
			im.SaveData(file, m.Pop[0].Fit, im.GetAvgFit(m.Pop[0].Nodes, im.CalRobustInfluenceByEdge))
		} else if evalType == 2 {
			im.SaveData(file, m.Pop[0].Fit, im.GetAvgFit(m.Pop[0].Nodes, im.CalRobustInfluenceByNode))
		} else {
			im.SaveData(file, m.Pop[0].Fit, im.GetAvgFit(m.Pop[0].Nodes, im.CalRobustInfluenceByNode), im.GetAvgFit(m.Pop[0].Nodes, im.CalRobustInfluenceByEdge))
		}
		fmt.Printf("gen-%d: %s\n", i, m.ExportBestSeed())
	}
}
