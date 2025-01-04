package im

// 计算seeds的影响力
func CalInfluence(seeds []int, adjList [][]int) float64 {
	influ := 0.0

	for _, seed := range seeds {
		influ += 1
		sumTemp := 0.0

		for _, seedAdjacency := range adjList[seed] {
			sumTemp += 1
			sumTemp += float64(len(adjList[seedAdjacency])) * ActivationProb
		}

		influ += ActivationProb * sumTemp
	}

	secSumTemp := 0.0
	thrSumTemp := 0.0

	for i := 0; i < len(seeds); i++ {
		seed := seeds[i]
		Cs := adjList[seed]
		var CsSimi []int
		findSimi(&CsSimi, Cs, seeds)

		sumTemp := 0.0
		for j := 0; j < len(CsSimi); j++ {
			var temp float64
			node := CsSimi[j]
			temp += 1
			temp += ActivationProb * float64(len(adjList[node]))
			sumTemp += temp * ActivationProb
		}
		secSumTemp += sumTemp

		var CsDisSimi []int
		var CsD []int
		findThird(&CsDisSimi, &CsD, Cs, seeds, seed)

		for i2 := 0; i2 < len(CsDisSimi); i2++ {
			for j := 0; j < len(CsD); j++ {
				if contains(adjList[CsDisSimi[i2]], CsD[j]) {
					thrSumTemp += ActivationProb * ActivationProb
				}
			}
		}
	}

	return (influ - secSumTemp - thrSumTemp)
}

func findSimi(CsSimi *[]int, Cs, seeds []int) {
	for i := 0; i < len(seeds); i++ {
		tempNode := seeds[i]
		for j := 0; j < len(Cs); j++ {
			if Cs[j] == tempNode {
				*CsSimi = append(*CsSimi, tempNode)
			}
		}
	}
}

func findThird(CsDisSimi, CsD *[]int, Cs, seeds []int, seed int) {
	CsDisSimiTemp := make([]int, len(Cs))
	copy(CsDisSimiTemp, Cs)

	for i := 0; i < len(CsDisSimiTemp); i++ {
		for j := 0; j < len(seeds); j++ {
			if CsDisSimiTemp[i] == seeds[j] {
				CsDisSimiTemp[i] = -1
			}
		}
		if CsDisSimiTemp[i] != -1 {
			*CsDisSimi = append(*CsDisSimi, CsDisSimiTemp[i])
		}
	}

	var CsSimi []int
	for i := 0; i < len(seeds); i++ {
		tempNode := seeds[i]
		for j := 0; j < len(Cs); j++ {
			if Cs[j] == tempNode {
				CsSimi = append(CsSimi, tempNode)
			}
		}
	}

	for i := 0; i < len(CsSimi); i++ {
		if CsSimi[i] != seed {
			*CsD = append(*CsD, CsSimi[i])
		}
	}
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
