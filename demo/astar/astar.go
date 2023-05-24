package astar

import "container/heap"

// Node 节点
// 使用 g 值和 h 值可以计算出从起点到目标节点的最短路径，这是 A* 算法的核心思想。
// 通过使用 g 值，可以确保扩展的节点是当前已知的最短路径上的节点。
// 通过使用 h 值，可以估计扩展节点到目标节点的最短距离，并向着目标节点前进。
// 使用 f 值可以将这两个因素结合起来，并选择当前最佳的节点进行扩展。
type Node struct {
	x int // 节点的横坐标
	y int // 节点的纵坐标

	// g 值表示起点到当前节点的距离。具体地，对于节点 n，g(n) 是从起点到 n 的实际距离，通常是通过从起点到 n 的路径上的边权之和来计算的。在 A* 算法中，g 值是在搜索过程中动态计算的。
	g int

	// h 值表示当前节点到终点的距离估算值。具体地，对于节点 n，h(n) 是从 n 到目标节点的估计距离。
	//   估价函数必须满足两个条件：
	//  	第一，它必须始终小于或等于从 n 到目标节点的实际距离；
	//      第二，它必须始终大于或等于从任意节点到目标节点的实际距离。在 A* 算法中，h 值是在搜索过程中动态计算的。
	h int

	// f 值是节点的总代价，即 f(n) = g(n) + h(n)。在 A* 算法中，f 值用于决定节点的优先级。A* 算法在扩展节点时总是选择 f 值最小的节点。
	f int

	// 父节点
	parent *Node
}

// PriorityQueue 优先队列
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// 求绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 它接受两个节点，并返回它们之间的启发式估价函数值。
// 该函数计算两个节点之间的曼哈顿距离，并将其作为启发式估价函数的值返回。
// 曼哈顿距离是指两点在矩形网格上沿着网格线行走的距离之和，这里是用节点的坐标差的绝对值之和来计算的。
func heuristic(a, b *Node) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

// CheckNodeLogic 检查该节点
type CheckNodeLogic func(neighbor *Node) (stop bool)

// Astar 实现了 A* 搜索算法，它接受一个起始节点、一个目标节点和一个二维网格数组，并返回一条从起始节点到目标节点的最短路径。
// A* 算法的主要思路如下：
// 1.将起点放入开放列表，并将其 f 值设置为起点的启发式估价函数值。同时将其 g 值（起点到该点的实际代价）设置为 0 , h 值（该点到终点的启发式估价函数值）。
// 2.从开放列表中选取 f 值最小的节点，将其从开放列表中移除并加入到封闭列表（Closed List）中。将该节点的相邻节点添加到开放列表中，如果该相邻节点尚未被探索过，则计算其 g 值和 h 值，并设置其父节点为当前节点。如果该相邻节点已经在开放列表中，则检查更新其 g 值和父节点，如果当前路径比原路径更短，则更新。
// 3.重复步骤 2，直到找到终点或开放列表为空。
// 4.如果找到终点，则从终点开始追踪父节点，直到回溯到起点，即可得到最短路径。
func Astar(start, end *Node, grid [][]int) []*Node {
	openList := make(PriorityQueue, 0)
	closedList := make(map[*Node]bool)
	start.g = 0
	start.h = heuristic(start, end)
	start.f = start.g + start.h
	heap.Push(&openList, start)

	for len(openList) > 0 {
		// 出栈 拿到最小f的node
		current := heap.Pop(&openList).(*Node)
		if current == end {
			path := make([]*Node, 0)
			for current != nil {
				path = append(path, current)
				current = current.parent
			}
			for i := 0; i < len(path)/2; i++ {
				path[i], path[len(path)-i-1] = path[len(path)-i-1], path[i]
			}
			return path
		}
		closedList[current] = true

		for _, neighbor := range getNeighbors(current, grid) {

			if closedList[neighbor] {
				continue
			}

			// 算一下 g
			tentativeGScore := current.g + 1

			if !inNodeList(neighbor, openList) || tentativeGScore < neighbor.g {
				neighbor.g = tentativeGScore
				neighbor.h = heuristic(neighbor, end)
				neighbor.f = neighbor.g + neighbor.h
				neighbor.parent = current

				if !inNodeList(neighbor, openList) {
					heap.Push(&openList, neighbor)
				}
			}
		}
	}
	return nil
}

// 它接受一个节点和一个节点指针的切片，用于检查该节点是否已经在开放列表中。
// 该函数使用一个循环遍历开放列表中的每个节点。
// 对于每个节点，它检查其 x 和 y 坐标是否与给定的节点相同。
// 如果找到一个与给定节点匹配的节点，则说明该节点已经在开放列表中，函数返回 true。
// 如果遍历完整个列表后没有找到匹配的节点，则说明该节点不在开放列表中，函数返回 false。
func inNodeList(n *Node, list []*Node) bool {
	for _, node := range list {
		if n.x == node.x && n.y == node.y {
			return true
		}
	}
	return false
}

// 它接受一个节点和一个二维网格数组，并返回一个包含所有相邻节点的指针数组。
// 这个函数首先创建一个空的 neighbors 数组，然后检查节点左、右、上、下四个方向的相邻节点是否可达。
// 如果可达，则将该相邻节点添加到 neighbors 数组中，并返回最终的相邻节点数组。
// 注意，这里我们假设所有可达的节点在网格上用数字 0 表示。
func getNeighbors(n *Node, grid [][]int) []*Node {
	neighbors := make([]*Node, 0)

	// 左边的节点
	if n.x > 0 && grid[n.x-1][n.y] == 0 {
		neighbors = append(neighbors, &Node{n.x - 1, n.y, 0, 0, 0, nil})
	}

	// 右边的节点
	if n.x < len(grid)-1 && grid[n.x+1][n.y] == 0 {
		neighbors = append(neighbors, &Node{n.x + 1, n.y, 0, 0, 0, nil})
	}

	// 上面的节点
	if n.y > 0 && grid[n.x][n.y-1] == 0 {
		neighbors = append(neighbors, &Node{n.x, n.y - 1, 0, 0, 0, nil})
	}

	// 下面的节点
	if n.y < len(grid[0])-1 && grid[n.x][n.y+1] == 0 {
		neighbors = append(neighbors, &Node{n.x, n.y + 1, 0, 0, 0, nil})
	}

	return neighbors
}
