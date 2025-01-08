package im

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Edge [2]int
type Load [2]float64            // <load, capacity>
type EdgeWithLoad map[Edge]Load //记录edge的负载load

// CalRobustInfluenceByNode 节点攻击下鲁棒影响力
func CalRobustInfluenceByNode(seeds []int) float64 {
	graph := AssignLoad(AdjList)
	var totalEdges = len(graph)
	sumFit, attackCnt := 0.0, 0
	for float32(len(graph)) > 0.1*float32(totalEdges) {
		node := graph.findMaxLoadNode()
		//node := graph.findRandomNode()
		graph.AttackNode(node)
		graph.removeOverloadEdge()
		sumFit += CalInfluence(seeds, graphToList(graph))
		attackCnt++
	}
	return sumFit / float64(attackCnt)
}

// 链路攻击下鲁棒影响力
func CalRobustInfluenceByEdge(seeds []int) float64 {
	graph := AssignLoad(AdjList)
	var totalEdges = len(graph)
	sumFit, attackCnt := 0.0, 0
	for float32(len(graph)) > 0.1*float32(totalEdges) {
		edge := graph.findMaxLoadEdge()
		//edge := graph.findRandomEdge()
		graph.AttackEdge(edge)
		graph.removeOverloadEdge()
		sumFit += CalInfluence(seeds, graphToList(graph))
		attackCnt++
	}
	return sumFit / float64(attackCnt)
}

// 图网络转化为邻接表
func graphToList(graph map[Edge]Load) [][]int {
	adjList := make([][]int, NetworkSize)
	for k := range graph {
		adjList[k[0]] = append(adjList[k[0]], k[1])
	}
	return adjList
}

// AssignLoad 分配负载
func AssignLoad(adjList [][]int) EdgeWithLoad {
	g := EdgeWithLoad{}
	for i := 0; i < NetworkSize; i++ {
		iDegree := len(adjList[i])
		for _, j := range adjList[i] {
			jDegree := len(adjList[j])
			load := math.Pow(float64(iDegree)*float64(jDegree), Alpha) //负载定义
			capacity := load * Beta                                    // 容积定义
			//分别记录<i,j>和<j,i>
			g[Edge{i, j}] = Load{load, capacity}
			g[Edge{j, i}] = Load{load, capacity}
		}
	}
	return g
}

func (g EdgeWithLoad) Attack() {
	edgeNum := len(g)
	attackCnt := 0
	for len(g) > 1 {
		if rand.Float32() < NodeAttackPer {
			node := g.findMaxLoadNode()
			//node := g.findRandomNode()
			g.AttackNode(node)
		} else {
			edge := g.findMaxLoadEdge()
			//edge := g.findRandomEdge()
			g.AttackEdge(edge)
		}
		g.removeOverloadEdge()
		attackCnt++
		fmt.Printf("故障次数:%d, scale:%.5f\n", attackCnt, float64(len(g))/float64(edgeNum))
	}
}

// AttackNode 攻击节点
func (g EdgeWithLoad) AttackNode(node int) {
	for _, val := range AdjList[node] {
		g.doRemoveEdge(Edge{node, val})
	}
}

// AttackEdge 攻击链路
func (g EdgeWithLoad) AttackEdge(edge Edge) {
	g.doRemoveEdge(edge)
}

// 移除边
func (g EdgeWithLoad) doRemoveEdge(edge Edge) {
	leftNode, rightNode := edge[0], edge[1]
	removedLoad := g[edge] // <load, cap>
	if removedLoad == [2]float64{0, 0} {
		return
	}
	delete(g, [2]int{leftNode, rightNode})
	delete(g, [2]int{rightNode, leftNode})

	adjCapacity := 0.0 //相邻负载总和
	for k, v := range g {
		if k[0] == leftNode || k[1] == leftNode || k[0] == rightNode || k[1] == rightNode {
			adjCapacity += v[1] / 2
		}
	}
	//分配转移负载,按capacity等比分配。《转移负载时会损失65%负载》
	for k, v := range g {
		if k[0] == leftNode || k[1] == leftNode || k[0] == rightNode || k[1] == rightNode {
			g[k] = Load{v[0] + (v[1]/adjCapacity)*removedLoad[0]*0.65, v[1]}
		}
	}
}

// 删除过载的边，并返回删除的边数
func (g EdgeWithLoad) removeOverloadEdge() {
	for {
		var overloadEdges []Edge
		for edge, load := range g {
			if load[0] > load[1] {
				overloadEdges = append(overloadEdges, edge)
			}
		}

		if len(overloadEdges) == 0 {
			break
		}

		sort.Slice(overloadEdges, func(i, j int) bool {
			if overloadEdges[i][0] != overloadEdges[j][0] {
				return overloadEdges[i][0] < overloadEdges[j][0]
			}
			return overloadEdges[i][1] < overloadEdges[j][1]
		})
		for _, edge := range overloadEdges {
			g.doRemoveEdge(edge)
		}
	}
}

// 找出网络中负载最大的边
func (g EdgeWithLoad) findMaxLoadEdge() Edge {
	var targetEdges []Edge
	var maxLoad float64 = -1
	for k, v := range g {
		if v[0] > maxLoad {
			targetEdges = []Edge{k}
			maxLoad = v[0]
		} else if v[0] == maxLoad {
			targetEdges = append(targetEdges, k)
		}
	}

	// 选择序号最小的edge
	sort.Slice(targetEdges, func(i, j int) bool {
		if targetEdges[i][0] != targetEdges[j][0] {
			return targetEdges[i][0] < targetEdges[j][0]
		}
		return targetEdges[i][1] < targetEdges[j][1]
	})
	return targetEdges[0]
}

// 随机选择一条链路
func (g EdgeWithLoad) findRandomEdge() Edge {
	randomIndex, cnt := rand.Intn(len(g)), 0
	for k := range g {
		if cnt == randomIndex {
			return k
		}
		cnt++
	}
	return Edge{0, 0}
}

// 找出网络中负载最大的节点
func (g EdgeWithLoad) findMaxLoadNode() int {
	//节点负载为相邻链路负载之和
	nodeMap := make(map[int]float64)
	for k, v := range g {
		nodeMap[k[0]] += v[0]
		nodeMap[k[1]] += v[0]
	}

	var maxNode int
	var maxLoad float64 = -1
	for k, v := range nodeMap {
		if v > maxLoad {
			maxNode = k
			maxLoad = v
		}
	}
	return maxNode
}

// 随机选择一个节点
func (g EdgeWithLoad) findRandomNode() int {
	randomIndex := rand.Intn(len(g))
	cnt := 0
	for k := range g {
		if cnt == randomIndex {
			return k[0]
		}
		cnt++
	}
	return 0
}
