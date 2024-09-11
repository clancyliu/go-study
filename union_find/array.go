package union_find

import "fmt"

// 初始化并查集
func initialize(n int) []int {
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	return parent
}

// 查找根节点
func find(parent []int, x int) int {
	return parent[x]
}

// 合并两个集合
func union(parent []int, x, y int) {
	rootX := find(parent, x)
	rootY := find(parent, y)
	if rootX != rootY {
		for i := 0; i < len(parent); i++ {
			if parent[i] == rootY {
				parent[i] = rootX
			}
		}
	}
}

func main() {
	parent := initialize(10)
	union(parent, 2, 3)
	union(parent, 3, 4)
	fmt.Println(find(parent, 4)) // 输出 2，表示节点 4 和 2 是连通的
}
