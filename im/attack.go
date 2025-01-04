package im

import (
	"fmt"
	"math"
	"math/rand"
)

type Edge [2]int
type Load [2]float64            // <load, capacity>
type EdgeWithLoad map[Edge]Load //记录edge的负载load

func Test() {
	g := AssignLoad(AdjList)
	g.Attack()
}

// CalRobustInfluenceByNode 节点攻击下鲁棒影响力
func CalRobustInfluenceByNode(seeds []int) float64 {
	g := AssignLoad(AdjList)
	sumFit := 0.0
	cnt := 0 // 故障次数
	for len(g) > 1 {
		node := g.findMaxLoadNode()
		//node := g.findRandomNode()
		g.AttackNode(node)
		g.removeOverloadEdge()
		adjList := graphToList(g)
		sumFit += CalInfluence(seeds, adjList)
		cnt++
	}
	return sumFit / float64(cnt)
}

// 链路攻击下鲁棒影响力
func CalRobustInfluenceByEdge(seeds []int) float64 {
	g := AssignLoad(AdjList)
	sumFit := 0.0
	cnt := 0 // 故障次数
	for len(g) > 1 {
		edge := g.findMaxLoadEdge()
		//edge := g.findRandomEdge()
		g.AttackEdge(edge)
		g.removeOverloadEdge()
		adjList := graphToList(g)
		sumFit += CalInfluence(seeds, adjList)
		cnt++
	}
	return sumFit / float64(cnt)
}

// 图网络转化为邻接表
func graphToList(graph map[Edge]Load) [][]int {
	adjList := make([][]int, NetworkSize)
	for k := range graph {
		adjList[k[0]] = append(adjList[k[0]], k[1])
	}
	return adjList
}

// 分配负载
func AssignLoad(adjList [][]int) EdgeWithLoad {
	g := EdgeWithLoad{}
	for i := 0; i < NetworkSize; i++ {
		i_degree := len(adjList[i])
		for _, j := range adjList[i] {
			j_degree := len(adjList[j])
			load := math.Pow(float64(i_degree)*float64(j_degree), Alpha) //负载定义
			cap := load * Beta                                           // 容积定义
			//分别记录<i,j>和<j,i>
			g[Edge{i, j}] = Load{load, cap}
			g[Edge{j, i}] = Load{load, cap}
		}
	}
	return g
}

func (g EdgeWithLoad) Attack() {
	edgeNum := len(g)
	cnt := 0
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
		cnt++
		fmt.Printf("故障次数:%d, scale:%.5f\n", cnt, float64(len(g))/float64(edgeNum))
	}
}

// 攻击节点
func (g EdgeWithLoad) AttackNode(node int) {
	for _, val := range AdjList[node] {
		g.doRemoveEdge(Edge{node, val})
	}
}

// 攻击链路
func (g EdgeWithLoad) AttackEdge(edge Edge) {
	g.doRemoveEdge(edge)
}

// 移除边
func (g EdgeWithLoad) doRemoveEdge(edge Edge) {
	left, right := edge[0], edge[1]
	removedLoad := g[edge] // <load, cap>
	if removedLoad == [2]float64{0, 0} {
		return
	}
	delete(g, [2]int{left, right})
	delete(g, [2]int{right, left})

	adjCap := 0.0 //相邻负载总和
	for k, v := range g {
		if k[0] == left || k[1] == left || k[0] == right || k[1] == right {
			adjCap += v[1] / 2
		}
	}
	//分配转移负载,按cap等比分配。《转移负载时会损失65%负载》
	for k, v := range g {
		if k[0] == left || k[1] == left || k[0] == right || k[1] == right {
			g[k] = Load{v[0] + (v[1]/adjCap)*removedLoad[0]*0.65, v[1]}
		}
	}
}

// 删除过载的边，并返回删除的边数
func (g EdgeWithLoad) removeOverloadEdge() {
	set := make(map[Edge]Load)
	for edge, load := range g {
		if load[0] > load[1] {
			set[edge] = load
		}
	}
	if len(set) == 0 {
		return
	}
	for edge := range set {
		g.doRemoveEdge(edge)
	}
	g.removeOverloadEdge()
}

// 找出网络中负载最大的边
func (g EdgeWithLoad) findMaxLoadEdge() Edge {
	var maxEdge Edge
	var maxLoad float64 = -1
	for k, v := range g {
		if v[0] > maxLoad {
			maxEdge = k
			maxLoad = v[0]
		}
	}
	return maxEdge
}

// 随机选择一条链路
func (g EdgeWithLoad) findRandomEdge() Edge {
	randomIndex := rand.Intn(len(g))
	cnt := 0
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
