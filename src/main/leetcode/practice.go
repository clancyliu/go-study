package leetcode

import (
	"math"
)

func magicalString(n int) int {
	queue := make([]int, 0)
	queue = append(queue, []int{2}...)
	count, preNum := 1, 1
	for len(queue) > 0 && n > 0 {
		curCount := queue[0]
		queue = queue[1:]
		curNum := 1
		if preNum == 1 {
			curNum = 2
		}
		if curCount == 1 {
			count++
		}
		for i := 0; i < curCount; i++ {
			queue = append(queue, curNum)
			n--
		}
		preNum = curNum
	}
	return count
}

var direct = [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
var queue [][2]int

func shortestBridge(grid [][]int) int {
	n := len(grid)
	queue = make([][2]int, 0)

	flag := false
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				visitLand(i, j, grid)
				flag = true
				break
			}
		}
		if flag {
			break
		}
	}
	step := 0
	tmpQueue := make([][2]int, 0)
	for len(queue) > 0 {
		firstItem := queue[0]
		queue = queue[1:]
		for _, item := range direct {
			x, y := firstItem[0]+item[0], firstItem[1]+item[1]
			if x < 0 || y < 0 || x >= n || y >= n || grid[x][y] == -1 {
				continue
			}
			if grid[x][y] == 1 {
				return step
			}
			if grid[x][y] == 0 {
				grid[x][y] = -1
				tmpQueue = append(tmpQueue, [2]int{x, y})
			}
		}
		if len(queue) == 0 && len(tmpQueue) > 0 {
			step++
			queue = tmpQueue
			tmpQueue = make([][2]int, 0)
		}
	}
	return step
}

func visitLand(i, j int, grid [][]int) {
	if i < 0 || j < 0 || i >= len(grid) || j >= len(grid[0]) || grid[i][j] != 1 {
		return
	}
	grid[i][j] = -1
	queue = append(queue, [2]int{i, j})
	for _, item := range direct {
		x, y := i+item[0], j+item[1]
		visitLand(x, y, grid)
	}
}

func partitionDisjoint(nums []int) int {
	leftMax := make([]int, len(nums))
	rightMin := make([]int, len(nums))
	leftMax[0] = nums[0]
	rightMin[len(nums)-1] = nums[len(nums)-1]
	for i := 1; i < len(nums); i++ {
		leftMax[i] = max(leftMax[i-1], nums[i])
	}
	for i := len(nums) - 2; i >= 0; i-- {
		rightMin[i] = min(rightMin[i+1], nums[i+1])
	}

	for i := 0; i < len(nums); i++ {
		if leftMax[i] <= rightMin[i] {
			return i + 1
		}
	}
	return -1
}

func isSubStructure(A *TreeNode, B *TreeNode) bool {
	if A == nil {
		return false
	}
	if B == nil {
		return true
	}
	if A.Val == B.Val && isSubStructureHelper(A, B) {
		return true
	}
	isSubStructure(A.Left, B)
	isSubStructure(A.Right, B)
	return false
}

func isSubStructureHelper(A, B *TreeNode) bool {
	if B == nil {
		return true
	}
	if A == nil || A.Val != B.Val {
		return false
	}
	left := isSubStructureHelper(A.Left, B.Left)
	right := isSubStructureHelper(A.Right, B.Right)

	return left && right
}

func cuttingRope(n int) int {
	dp := make([][]int, n+1)
	ans := 0
	for i := 2; i <= n; i++ {
		dp[i] = make([]int, n+1)
		for j := 2; j <= i; j++ {
			if j == 2 {
				dp[i][j] = maxTwoNum(i)
				continue
			}
			for k := i - 1; k > 1; k-- {
				dp[i][j] = max(dp[i][j], dp[k][j-1]*(i-k))
				ans = max(ans, dp[i][j])
			}
		}
	}
	return ans
}

func maxTwoNum(i int) int {
	ans := 0
	for num := i / 2; num > 0; num-- {
		ans = max(ans, num*(i-num))
	}
	return ans
}

func findNumberIn2DArray(matrix [][]int, target int) bool {
	if len(matrix) == 0 {
		return false
	}
	n, m := len(matrix), len(matrix[0])
	i, j := 0, m-1

	for i < n && j >= 0 {
		if target == matrix[i][j] {
			return true
		}
		if target > matrix[i][j] {
			i++
			continue
		}
		if target < matrix[i][j] {
			j--
			continue
		}
	}
	return false
}

type StockSpanner struct {
	stack [][2]int
	index int
}

func Constructor() StockSpanner {
	stack := make([][2]int, 0)
	stack = append(stack, [2]int{-1, math.MaxInt64})
	return StockSpanner{
		stack: stack,
	}
}

func (this *StockSpanner) Next(price int) int {
	curIndex := this.index + 1
	for this.stack[len(this.stack)-1][1] < price {
		this.stack = this.stack[:len(this.stack)-1]
	}
	ans := curIndex - this.stack[len(this.stack)-1][0]
	this.index = curIndex
	this.stack = append(this.stack, [2]int{curIndex, price})
	return ans
}

func dailyTemperatures(temperatures []int) []int {
	stack := make([]int, 0)
	ans := make([]int, len(temperatures))
	for i, item := range temperatures {
		for len(stack) > 0 && item > temperatures[stack[len(stack)-1]] {
			lastItem := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			ans[lastItem] = i - lastItem
		}
		stack = append(stack, i)
	}
	return ans
}

func countSubstrings(s string) int {
	ans := 0
	for i := 0; i < len(s); i++ {
		ans++
		ans += reverseCount(s, i-1, i)
		ans += reverseCount(s, i-1, i+1)
	}
	return ans
}

func reverseCount(s string, left, right int) int {
	tmpCount := 0
	for left >= 0 && right < len(s) {
		if s[left] != s[right] {
			break
		}
		tmpCount++
		left--
		right++
	}
	return tmpCount
}

type myHeap [][2]int

func (h *myHeap) Less(i, j int) bool {
	if (*h)[i][0] == (*h)[j][0] {
		return (*h)[i][1] > (*h)[j][1]
	}
	return (*h)[i][0] < (*h)[j][0]
}

func (h *myHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *myHeap) Len() int {
	return len(*h)
}

func (h *myHeap) Pop() (v any) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *myHeap) Push(v any) {
	*h = append(*h, v.([2]int))
}

func (h *myHeap) Size() int {
	return len(*h)
}

func leastInterval(tasks []byte, n int) int {
	arrCount := make([]int, 26)
	for _, task := range tasks {
		arrCount[task-'A']++
	}
	tmpMax, maxCount := 0, 0
	for _, item := range arrCount {
		if item > tmpMax {
			tmpMax = item
			maxCount = 1
		} else if item == tmpMax {
			maxCount++
		}
	}
	return max(len(tasks), (tmpMax-1)*(n+1)+maxCount)
}

func mergeTrees(root1 *TreeNode, root2 *TreeNode) *TreeNode {
	if root1 == nil && root2 == nil {
		return nil
	}
	if root1 == nil {
		return root2
	}
	if root2 == nil {
		return root1
	}

	root := &TreeNode{Val: root1.Val + root2.Val}
	root.Left = mergeTrees(root1.Left, root2.Left)
	root.Right = mergeTrees(root1.Right, root2.Right)
	return root
}

func findUnsortedSubarray(nums []int) int {
	iStart, iEnd := 0, len(nums)-1
	for iStart < len(nums)-1 && nums[iStart+1] >= nums[iStart] {
		iStart++
	}
	for iEnd > 0 && nums[iEnd] >= nums[iEnd-1] {
		iEnd--
	}
	tmpMin, tmpMax := math.MaxInt64, math.MinInt64
	for i := iStart; i <= iEnd; i++ {
		tmpMin = min(tmpMin, nums[i])
		tmpMax = max(tmpMax, nums[i])
	}
	for iStart >= 0 && nums[iStart] > tmpMin {
		iStart--
	}
	for iEnd < len(nums) && nums[iEnd] < tmpMax {
		iEnd++
	}
	if iStart < iEnd {
		return iEnd - iStart - 1
	}
	return 0
}

func diameterOfBinaryTree(root *TreeNode) int {
	tmpMaxDis = 0
	treeDeep(root)
	return tmpMaxDis
}

var tmpMaxDis int

func treeDeep(root *TreeNode) int {
	if root == nil {
		return 0
	}
	left := treeDeep(root.Left)
	right := treeDeep(root.Right)
	curMax := max(left, right) + 1
	tmpMaxDis = max(tmpMaxDis, left+right+1)
	return curMax
}

func convertBST(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	dfs(root)

	return root
}

func dfs(root *TreeNode) (sum int) {
	if root == nil {
		return
	}
	sum += dfs(root.Right)
	sum += root.Val
	root.Val = sum
	sum += dfs(root.Left)
	return
}

func findTargetSumWays1(nums []int, target int) int {
	tmpCount = 0
	findTargetSumWaysHelper(nums, 0, target)
	return tmpCount
}

var tmpCount int

func findTargetSumWaysHelper(nums []int, index, target int) {
	if index == len(nums) {
		if target == 0 {
			tmpCount++
		}
		return
	}

	findTargetSumWaysHelper(nums, index+1, target+nums[index])
	findTargetSumWaysHelper(nums, index+1, target-nums[index])
}

func findDisappearedNumbers(nums []int) []int {
	n := len(nums)
	for _, num := range nums {
		nums[num-1] += n
	}
	ans := make([]int, 0)
	for i, num := range nums {
		if num < n {
			ans = append(ans, i+1)
		}
	}
	return ans
}

func findAnagrams(s string, p string) []int {
	ans := make([]int, 0)
	if len(s) < len(p) {
		return ans
	}
	var pArr, sArr [26]int
	for i, ch := range p {
		pArr[ch-'a']++
		sArr[s[i]-'a']++
	}
	if pArr == sArr {
		ans = append(ans, 0)
	}
	for i := len(p); i < len(s); i++ {
		sArr[s[i-len(p)]]--
		sArr[s[i]]++
		if pArr == sArr {
			ans = append(ans, i)
		}
	}
	return ans
}

func pathSum(root *TreeNode, targetSum int) int {
	if root == nil {
		return 0
	}
	ans := pathSumHelper(root, targetSum)
	// 遍历整个树
	ans += pathSum(root.Left, targetSum)
	ans += pathSum(root.Right, targetSum)

	return ans
}

// 计算从root开始到叶子节点，和为targetSum的个数
func pathSumHelper(root *TreeNode, targetSum int) (res int) {
	if root == nil {
		return
	}
	if targetSum == root.Val {
		res++
	}
	res += pathSumHelper(root.Left, targetSum-root.Val)
	res += pathSumHelper(root.Right, targetSum-root.Val)
	return
}

func canPartition(nums []int) bool {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	if sum%2 == 1 {
		return false
	}
	target := sum / 2
	// dp[i][j] 表示nums中前i个元素，是否能构成和为j
	dp := make([][]bool, len(nums))
	for i := 0; i < len(nums); i++ {
		dp[i] = make([]bool, target+1)
		dp[i][0] = true
	}

	for i := 0; i < len(nums); i++ {
		for j := 0; j <= target; j++ {
			if j == 0 {
				dp[i][j] = true
				continue
			}
			if i == 0 {
				dp[i][j] = nums[i] == j
				continue
			}
			dp[i][j] = dp[i-1][j]
			if nums[i] <= j {
				dp[i][j] = dp[i-1][j] || dp[i-1][j-nums[i]]
			}
		}
	}
	return dp[len(nums)-1][target]
}

func subarraySum(nums []int, k int) int {
	count, preSum := 0, 0
	tmpMap := make(map[int]int, 0)
	tmpMap[0] = 1
	for _, num := range nums {
		preSum += num
		tmpMap[preSum]++
		// 第i个preSum - 第j个preSum == k
		if tmpCount, ok := tmpMap[preSum-k]; ok {
			count += tmpCount
		}
	}
	return count
}

func buildArray(target []int, n int) []string {
	ans := make([]string, 0)
	for i, j := 0, 1; i < len(target) && j <= n; j++ {
		if target[i] == j {
			ans = append(ans, "Push")
			i++
		} else {
			ans = append(ans, "Push", "Pop")
		}
	}
	return ans
}

func buildTree(preorder []int, inorder []int) *TreeNode {
	return buildTreeHelper(preorder, inorder, 0, len(preorder), 0, len(inorder))
}

func buildTreeHelper(preorder, inorder []int, pStart, pEnd, iStart, iEnd int) *TreeNode {
	if pStart == pEnd {
		return nil
	}
	rootVal := preorder[pStart]
	rootNode := &TreeNode{Val: rootVal}
	iLeft, iRight := 0, 0
	for i, item := range inorder {
		if item == rootVal {
			iLeft = i
			iRight = i + 1
			break
		}
	}
	pLeftNum := iLeft - iStart
	rootNode.Left = buildTreeHelper(preorder, inorder, pStart+1, pStart+pLeftNum+1, iStart, iLeft)
	rootNode.Right = buildTreeHelper(preorder, inorder, pStart+pLeftNum+1, pEnd, iRight, iEnd)

	return rootNode
}
