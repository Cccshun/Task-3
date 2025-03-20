package algo

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"

	"sysu.com/task3/im"
)

type GA struct {
	wg     sync.WaitGroup
	Pop    []im.Seed
	NewPop []im.Seed
}

func (g *GA) Init(evalType int) {
	g.Pop = make([]im.Seed, im.PopSize)
	for i := 0; i < im.PopSize; i++ {
		g.Pop[i] = im.NewSeed(evalType)
	}
}

// Crossover 交叉
func (g *GA) Crossover(evalType int) {
	g.NewPop = im.DeepCopyPop(g.Pop)
	rand.Shuffle(len(g.NewPop), func(i, j int) {
		g.NewPop[i], g.NewPop[j] = g.NewPop[j], g.NewPop[i]
	})
	for i := 0; i < im.PopSize; i += 2 {
		if rand.Float32() < im.PC {
			g.doCrossover(&g.NewPop[i], &g.NewPop[i+1])
			im.RemoveDuplicateGene(g.NewPop[i])
			im.RemoveDuplicateGene(g.NewPop[i+1])
			g.wg.Add(2)
			go im.EvaluateSeedSync(&g.NewPop[i], &g.wg, evalType)
			go im.EvaluateSeedSync(&g.NewPop[i+1], &g.wg, evalType)
		}
	}
	g.wg.Wait()
}

// 均匀交叉
func (g *GA) doCrossover(seed1, seed2 *im.Seed) {
	for i := 0; i < im.SeedSize; i++ {
		// todo 0.5?
		if rand.Float32() < 0.5 {
			seed1.Nodes[i], seed2.Nodes[i] = seed2.Nodes[i], seed1.Nodes[i]
		}
	}
}

// Mutate 变异
func (g *GA) Mutate(evalType int) {
	for i := range g.NewPop {
		g.doMutate(&g.NewPop[i])
		im.RemoveDuplicateGene(g.NewPop[i])
		g.wg.Add(1)
		go im.EvaluateSeedSync(&g.NewPop[i], &g.wg, evalType)
	}
	g.wg.Wait()
}

// 单点变异
func (g *GA) doMutate(seed *im.Seed) {
	for i := range seed.Nodes {
		if rand.Float32() < im.PM {
			seed.Nodes[i] = im.NewGene()
		}
	}
}

// Select 选择策略，轮盘赌
func (g *GA) Select() {
	mergedPop := im.DeepCopyPop(g.Pop)
	mergedPop = append(mergedPop, g.NewPop...)

	g.Pop = im.RouletteSelection(mergedPop)
	sort.Sort(im.BySeed(g.Pop))
}

func (g *GA) FindBestSeed(savePath string, evalType int, wg *sync.WaitGroup) {
	defer wg.Done()
	g.Init(evalType)
	file := im.CreateDataPath(savePath, "ga")
	defer file.Close()
	for i := 0; i < im.MaxGen; i++ {
		g.Crossover(evalType)
		g.Mutate(evalType)
		g.Select()
		switch evalType {
		case 1:
			im.SaveData(file, g.Pop[0].Fit, im.GetAvgFit(g.Pop[0].Nodes, im.CalRobustInfluenceByEdge))
		case 2:
			im.SaveData(file, g.Pop[0].Fit, im.GetAvgFit(g.Pop[0].Nodes, im.CalRobustInfluenceByNode))
		case 3:
			im.SaveData(file, g.Pop[0].Fit, im.GetAvgFit(g.Pop[0].Nodes, im.CalRobustInfluenceByNode), im.GetAvgFit(g.Pop[0].Nodes, im.CalRobustInfluenceByEdge))
		default:
			_, _ = fmt.Fprintln(os.Stderr, "parameter error")
		}
		fmt.Printf("ga-gen-%d: %s\n", i, g.ExportBestSeed())
	}
}

func (g *GA) ExportBestSeed() string {
	return fmt.Sprintf("%v, NodeAttack:%.3f, EdgeAttack:%.3f",
		g.Pop[0], im.GetAvgFit(g.Pop[0].Nodes, im.CalRobustInfluenceByNode), im.GetAvgFit(g.Pop[0].Nodes, im.CalRobustInfluenceByEdge))
}

func (g *GA) ExportPop() {
	for idx, elem := range g.Pop {
		fmt.Printf("Pop--%d:%+v\n", idx, elem)
	}
}

func (g *GA) ExportNewPop() {
	for idx, elem := range g.NewPop {
		fmt.Printf("NewPop--%d:%+v\n", idx, elem)
	}
}

// 输出种群个体适应度信息
func (g *GA) ExportEvolutionInfo() string {
	str := "[ "
	for idx, elem := range g.Pop {
		str += fmt.Sprintf("%d:%.2f ", idx, elem.Fit)
	}
	str += fmt.Sprintf("]")
	return str
}

// 输出种群中不同个体的比例
func (g *GA) ExportScale() float64 {
	hashTable := map[uint64]struct{}{}
	// 统计Pop种群中不重复的个体数量
	for _, seed := range g.Pop {
		hashTable[seed.Hash()] = struct{}{}
	}
	return float64(len(hashTable)) / (im.PopSize)
}
