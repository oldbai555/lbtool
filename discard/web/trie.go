package web

import "strings"

// node 路由节点
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// matchChild 找到第一个匹配成功的节点，用于插入
// params part 通过路径找到符合的子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 两种情况,
		// 1. 该 node 下的 child 存在精确匹配的路径
		// 2. 该 node 下的 child 存在模糊匹配的路径参数
		if child.part == part || child.isWild {
			return child
		}
	}
	// 都找不到,表示是新插入的节点
	return nil
}

// matchChildren 所有匹配成功的节点，用于查找
// params part 通过路径找到符合的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {

		// 两种情况,
		// 1. 该 node 下的 child 存在精确匹配的路径
		// 2. 该 node 下的 child 存在模糊匹配的路径参数
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 注册节点
// params pattern 路由路径 例如 /p/a/:name
// params parts 解析后的路由 例如 ["p","a",":name"]
// params height 路由深度,从 0 开始
func (n *node) insert(pattern string, parts []string, height int) {
	// 将解析后的路由 parts 都给遍历完,到最后一位，就结束插入
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	// 找不到子节点,表示是新的节点
	if child == nil {
		// part[0] == ':' || part[0] == '*' 用于判断字符串第一位是否是需要 匹配路径参数
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 递归插入
	child.insert(pattern, parts, height+1)
}

// search 匹配
// params parts 解析后的路由 例如 ["p","a",":name"]
// params height 路由深度,从 0 开始
func (n *node) search(parts []string, height int) *node {
	// 两种情况
	// 1. parts 解析的路由遍历结束
	// 2. 遇到了通配符
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 遍历结束也找不到匹配的子节点
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		// 一直往下找，直到找到第一个符合条件的路径
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
