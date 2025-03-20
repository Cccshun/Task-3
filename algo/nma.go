package algo

import (
	"fmt"
	"math/rand"
	"sort"
	"sysu.com/task3/im"
)

type NMA struct {
	MA
	NicheList [][]im.Seed
}

func (n *NMA) FindBestSeed(savePath string, evalType int) {
	n.Init(evalType)
	file := im.CreateDataPath(savePath, "nma")
	defer file.Close()

	for i := 0; i < im.MaxGen; i++ {
		n.AllocateNiches()
		n.NicheCrossover(evalType)
		n.Mutate(evalType)
		n.DeduplicateSeed()
		n.PopAlign(evalType)
		n.LocalSearch(evalType)
		n.Select()
		if evalType == 1 {
			im.SaveData(file, n.Pop[0].Fit, im.GetAvgFit(n.Pop[0].Nodes, im.CalRobustInfluenceByEdge))
		} else if evalType == 2 {
			im.SaveData(file, n.Pop[0].Fit, im.GetAvgFit(n.Pop[0].Nodes, im.CalRobustInfluenceByNode))
		} else {
			im.SaveData(file, n.Pop[0].Fit, im.GetAvgFit(n.Pop[0].Nodes, im.CalRobustInfluenceByNode), im.GetAvgFit(n.Pop[0].Nodes, im.CalRobustInfluenceByEdge))
		}
		fmt.Printf("gen-%d: %s\n", i, n.ExportBestSeed())
	}
}

func (n *NMA) AllocateNiches() {
	var nicheSize = im.PopSize / im.NicheNumber
	n.NicheList = make([][]im.Seed, im.NicheNumber)

	var tmpPop = im.DeepCopyPop(n.Pop)
	for i := 0; i < im.NicheNumber; i++ {
		// 确定第i个niche leader
		sort.Slice(tmpPop, func(i, j int) bool {
			return tmpPop[i].Fit < tmpPop[j].Fit
		})
		var nicheLeader = tmpPop[len(tmpPop)-1]
		tmpPop = tmpPop[:len(tmpPop)-1]
		n.NicheList[i] = []im.Seed{nicheLeader}

		// 填充第i个niche
		sort.Slice(tmpPop, func(i, j int) bool {
			return nicheLeader.SimilarityCalculate(tmpPop[i]) < nicheLeader.SimilarityCalculate(tmpPop[j])
		})
		if i == im.NicheNumber-1 {
			n.NicheList[i] = append(n.NicheList[i], tmpPop[:]...)
		} else {
			n.NicheList[i] = append(n.NicheList[i], tmpPop[len(tmpPop)-nicheSize+1:]...)
		}
		tmpPop = tmpPop[:len(tmpPop)-nicheSize+1]
	}

	n.NewPop = make([]im.Seed, 0) // todo
	for i := 0; i < im.NicheNumber; i++ {
		n.NewPop = append(n.NewPop, n.NicheList[i]...)
	}
}

func (n *NMA) NicheCrossover(evalType int) {
	for i := 0; i < 3; i++ {
		for currentNiche := 0; currentNiche < im.NicheNumber; {
			targetNiche := rand.Intn(im.NicheNumber)
			if targetNiche == currentNiche {
				continue
			}
			x, y := rand.Intn(len(n.NicheList[currentNiche])), rand.Intn(len(n.NicheList[targetNiche]))
			n.doCrossover(&n.NicheList[currentNiche][x], &n.NicheList[targetNiche][y])
			n.wg.Add(2)
			go im.EvaluateSeedSync(&n.NicheList[currentNiche][x], &n.wg, evalType)
			go im.EvaluateSeedSync(&n.NicheList[targetNiche][y], &n.wg, evalType)
			n.wg.Wait()
			currentNiche++
		}
	}

	var tmp []im.Seed
	for _, list := range n.NicheList {
		tmp = append(tmp, list...)
	}
	fmt.Printf("交叉后:%d   ", similarity(n.Pop, tmp))
}

func (n *NMA) DeduplicateSeed() {
	var tmpPop []im.Seed
	for i := 0; i < len(n.NewPop); i++ {
		for j := 0; j < len(n.Pop); j++ {
			if !n.NewPop[i].Equal(n.Pop[j]) {
				tmpPop = append(tmpPop, n.NewPop[i])
			}
		}
	}
	n.NewPop = tmpPop

	fmt.Printf("去重后:%d  ", similarity(n.Pop, n.NewPop))
}

func (n *NMA) PopAlign(evalType int) {
	for len(n.NewPop) < im.PopSize {
		n.NewPop = append(n.NewPop, im.NewSeed(evalType))
	}
}

func similarity(source, other []im.Seed) int {
	cnt := 0
	for _, i := range source {
		for _, j := range other {
			if i.Equal(j) {
				cnt++
			}
		}
	}
	return cnt
}
