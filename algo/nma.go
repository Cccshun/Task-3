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
	n.Init()
	file := im.CreateDataPath(savePath, "nma")
	defer file.Close()

	for i := 0; i < im.MaxGen; i++ {
		n.AllocateNiches()
		n.NicheCrossover()
		n.Mutate(evalType)
		n.DeduplicateSeed()
		n.PopAlign()
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

	n.NewPop = make([]im.Seed, 0)
	for i := 0; i < im.NicheNumber; i++ {
		n.NewPop = append(n.NewPop, n.NicheList[i]...)
	}
}

func (n *NMA) NicheCrossover() {
	for i := 0; i < 3; i++ {
		for currentNiche := 0; currentNiche < im.NicheNumber; {
			targetNiche := rand.Intn(im.NicheNumber - 1)
			if targetNiche == currentNiche {
				continue
			}
			x, y := rand.Intn(len(n.NicheList[currentNiche])), rand.Intn(len(n.NicheList[targetNiche]))
			n.doCrossover(&n.NicheList[currentNiche][x], &n.NicheList[targetNiche][y])
			currentNiche++
		}
	}
}

func (n *NMA) DeduplicateSeed() {
	for i := 0; i < len(n.NewPop); i++ {
		for j := i + 1; j < len(n.NewPop); j++ {
			if n.NewPop[i].Equal(n.NewPop[j]) {
				n.NewPop = append(n.NewPop[:j], n.NewPop[j+1:]...)
			}
		}
	}
}

func (n *NMA) PopAlign() {
	for len(n.NewPop) < im.PopSize {
		n.NewPop = append(n.NewPop, im.NewSeed())
	}
}
