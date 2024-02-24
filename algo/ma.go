package algo

import (
	"fmt"
	"math/rand"
	"sync"

	"sysu.com/task3/im"
)

type Ma struct {
	Ga
}

// 局部搜索
func (m *Ma) LocalSearch() {
	for i := 0; i < im.PopSize; i++ {
		if rand.Float32() < im.PL {
			m.wg.Add(1)
			go m.doLocalSearch(&m.NewPop[i])
		}
	}
	m.wg.Wait()
}

func (m *Ma) doLocalSearch(seed *im.Seed) {
	defer m.wg.Done()
	for i := 0; i < im.SeedSize; i++ {
		CompareAndSwap(seed, i)
	}
}

// 搜索2-hop内最优种子
func CompareAndSwap(seed *im.Seed, idx int) {
	nodes := im.Get2HopNodes(seed.Nodes[idx])
	var wg sync.WaitGroup
	var mu sync.Mutex

	//用nodes中的节点替换seed.Nodes[idx],并评估适应度
	seedsMap := make(map[*im.Seed]float64)
	for node := range nodes {
		seedCopy := im.DeepCopySeed(*seed)
		seedCopy.Nodes[idx] = node
		wg.Add(1)
		go im.EvaluteSeedAsync(&seedCopy, seedsMap, &wg, &mu)
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

func (m *Ma) FindSeed() {
	m.Init()
	for i := 0; i < im.MaxGen; i++ {
		m.Crossover()
		m.Mutate()
		m.LocalSearch()
		m.Select()
		fmt.Printf("gen--%d: %s, sacle:%.2f\n", i, m.ExportEvolutionInfo(), m.ExportScale())
	}
}
