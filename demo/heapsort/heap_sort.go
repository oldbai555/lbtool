package heapsort

import "fmt"

// SortMax 最小堆排序
// 创建最小堆
// 调整堆
// 找到最大的那个值,放到首位
// 交换首尾节点(为了维持一个完全二叉树才要进行收尾交换)
func SortMax(arr []int, length int) []int {
	// length := len(arr)
	if length <= 1 {
		return arr
	}
	// 二叉树深度
	depth := length/2 - 1
	// 从最后一个 三角形区域(父节点、左右子节点中的三角区域)进行筛选最大值
	for i := depth; i >= 0; i-- {
		// 假定最大的位置就在i的位置
		topMax := i
		// 左子节点
		leftChild := 2*i + 1
		// 右子节点
		rightChild := 2*i + 2
		// 防止越过界限 拿到父节点、左右子节点中的三者的最大值
		if leftChild <= length-1 && arr[leftChild] > arr[topMax] {
			topMax = leftChild
		}
		// 防止越过界限 拿到父节点、左右子节点中的三者的最大值
		if rightChild <= length-1 && arr[rightChild] > arr[topMax] {
			topMax = rightChild
		}
		// 交换元素
		if topMax != i {
			arr[i], arr[topMax] = arr[topMax], arr[i]
		}
	}
	return arr
}

// SortMin 最大堆排序
// 创建最大堆
// 调整堆
// 找到最小的那个值,放到首位
// 交换首尾节点(为了维持一个完全二叉树才要进行收尾交换)
func SortMin(arr []int, length int) []int {
	// length := len(arr)
	if length <= 1 {
		return arr
	}
	// 二叉树深度
	depth := length/2 - 1
	// 从最后一个 三角形区域(父节点、左右子节点中的三角区域)进行筛选最小值
	for i := depth; i >= 0; i-- {
		// 假定最大的位置就在i的位置
		topMin := i
		// 左子节点
		leftChild := 2*i + 1
		// 右子节点
		rightChild := 2*i + 2
		// 防止越过界限 拿到父节点、左右子节点中的三者的最小值
		if leftChild <= length-1 && arr[leftChild] < arr[topMin] {
			topMin = leftChild
		}
		// 防止越过界限 拿到父节点、左右子节点中的三者的最小值
		if rightChild <= length-1 && arr[rightChild] < arr[topMin] {
			topMin = rightChild
		}
		// 交换元素
		if topMin != i {
			arr[i], arr[topMin] = arr[topMin], arr[i]
		}
	}
	return arr
}

// HeapSort 堆排序
func HeapSort(arr []int) []int {
	// 数组长度
	length := len(arr)
	for i := 0; i < length; i++ {
		lastLen := length - i
		SortMax(arr, lastLen)
		if i < length {
			arr[0], arr[lastLen-1] = arr[lastLen-1], arr[0]
		}
		fmt.Println(arr)
	}
	return arr
}
