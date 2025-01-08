package lc

import (
	"math"
	"sort"
)

// SortedListToBST 有序链表转BST
func SortedListToBST(head *ListNode) *TreeNode {
	var f func(*ListNode) *TreeNode
	f = func(node *ListNode) *TreeNode {
		if node == nil {
			return nil
		}

		mid := getMiddleNodeFromList(node)
		var root = &TreeNode{Val: mid.Val}

		if node == mid {
			root.Left = nil
		} else {
			var prev = node
			for prev.Next != mid {
				prev = prev.Next
			}
			prev.Next = nil
			root.Left = f(node)
		}

		next := mid.Next
		mid.Next = nil
		root.Right = f(next)

		return root
	}
	return f(head)
}

func ArrayToListNode(nums []int) *ListNode {
	if nums == nil || len(nums) == 0 {
		return nil
	}
	dummy := &ListNode{}
	var p = dummy
	for _, num := range nums {
		p.Next = &ListNode{Val: num}
		p = p.Next
	}
	return dummy.Next
}

// 获取链表中点
func getMiddleNodeFromList(head *ListNode) *ListNode {
	var dummy = &ListNode{Val: -1, Next: head}
	slowPtr, fastPtr := dummy, dummy
	for fastPtr != nil {
		fastPtr = fastPtr.Next
		if fastPtr != nil {
			fastPtr = fastPtr.Next
			slowPtr = slowPtr.Next
		}
	}
	return slowPtr
}

// 有序数组转BST
func sortedArrayToBST(nums []int) *TreeNode {
	var f func(int, int) *TreeNode
	f = func(l, r int) *TreeNode {
		if l > r {
			return nil
		}

		var m = (l + r) / 2
		var root = &TreeNode{Val: nums[m]}
		root.Left = f(l, m-1)
		root.Right = f(m+1, r)
		return root
	}
	return f(0, len(nums)-1)
}

func convertBST(root *TreeNode) *TreeNode {
	var preTotal int
	var dfs func(*TreeNode)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}

		dfs(root.Right)
		preTotal += root.Val
		root.Val = preTotal
		dfs(root.Left)
	}
	dfs(root)
	return root
}

func latestTimeCatchTheBus(buses []int, passengers []int, capacity int) int {
	sort.Ints(buses)
	sort.Ints(passengers)
	var pos, space int
	for _, bus := range buses {
		space = capacity
		for space > 0 && pos < len(passengers) && passengers[pos] <= bus {
			space--
			pos++
		}
	}

	pos--
	var lastCatchTime int
	if space > 0 {
		lastCatchTime = buses[len(buses)-1]
	} else {
		lastCatchTime = passengers[pos]
	}
	for pos >= 0 && lastCatchTime == passengers[pos] {
		pos--
		lastCatchTime--
	}
	return lastCatchTime
}

func numBusesToDestination(routes [][]int, source int, target int) int {
	cache := make(map[int][]int)
	for routeIdx, route := range routes {
		for _, station := range route {
			cache[station] = append(cache[station], routeIdx)
		}
	}

	ans := math.MaxInt
	visited := make(map[int]bool)
	var backtrack func(int, int)
	backtrack = func(currStation, cnt int) {
		if currStation == target {
			ans = min(ans, cnt)
			return
		}

		for _, routeIdx := range cache[currStation] {
			for _, nextStation := range routes[routeIdx] {
				if !visited[nextStation] {
					visited[nextStation] = true
					backtrack(nextStation, cnt+1)
					visited[nextStation] = false
				}
			}
		}
	}
	backtrack(source, 0)
	if ans == math.MaxInt {
		return -1
	}
	return ans
}

func uniquePathsIII(grid [][]int) int {
	n, m := len(grid), len(grid[0])
	var startX, startY, st int
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 1 {
				startX, startY = i, j
			} else if grid[i][j] == -1 {
				st |= 1 << (i*m + j)
			}
		}
	}

	cache := make(map[int]int)
	var backtrack func(int, int, int) int
	backtrack = func(x, y int, visited int) int {
		rank := x*m + y
		if x < 0 || x >= n || y < 0 || y >= m || (visited&(1<<rank)) == 1 {
			return 0
		}
		visited |= 1 << rank

		if grid[x][y] == 2 {
			if visited == (1<<(n*m) - 1) {
				return 1
			}
			return 0
		}

		key := rank<<(n*m) | visited
		if val, ok := cache[key]; !ok {
			ans := backtrack(x+1, y, visited) +
				backtrack(x-1, y, visited) +
				backtrack(x, y+1, visited) +
				backtrack(x, y-1, visited)
			cache[key] = ans
			return ans
		} else {
			return val
		}
	}

	return backtrack(startX, startY, st)
}

func uniquePathsWithObstacles(obstacleGrid [][]int) int {
	n := len(obstacleGrid)
	m := len(obstacleGrid[0])
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, m)
	}
	if obstacleGrid[0][0] == 1 || obstacleGrid[n-1][m-1] == 1 {
		return 0
	}
	dp[0][0] = 1
	for i := 1; i < n && obstacleGrid[i][0] != 1; i++ {
		dp[i][0] += dp[i-1][0]
	}
	for i := 1; i < m && obstacleGrid[0][i] != 1; i++ {
		dp[0][i] += dp[0][i-1]
	}

	for i := 1; i < n; i++ {
		for j := 1; j < m; j++ {
			if obstacleGrid[i][j] != 1 {
				dp[i][j] = dp[i-1][j] + dp[i][j-1]
			}
		}
	}
	return dp[n-1][m-1]
}

func uniquePaths(m int, n int) int {
	prev := make([]int, n)
	for i := 0; i < n; i++ {
		prev[i] = 1
	}

	for i := 1; i < m; i++ {
		curr := make([]int, n)
		curr[0] = 1
		for j := 1; j < n; j++ {
			curr[j] = prev[j] + curr[j-1]
		}
		prev = curr
	}
	return prev[n-1]
}
