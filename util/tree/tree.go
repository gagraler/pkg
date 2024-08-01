package tree

import (
	"sort"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-08-01 16:29:09
 * @file: tree.go
 * @description: tree util 处理任意类型的树形数据结构
 */

// Predicate 对元素进行判断
type Predicate[E any] func(E) bool

// BiFunction 接收两个元素，返回一个布尔值
type BiFunction[E any] func(E, E) bool

// BiConsumer 接收两个元素，返回一个元素
type BiConsumer[E any] func(E, []E)

// Function 接收一个元素，返回一个list
type Function[E any] func(E) []E

// Consumer 接收一个元素
type Consumer[E any] func(E)

// MakeTree 从list中构建tree
func MakeTree[E any](list []E, rootCheck Predicate[E], parentCheck BiFunction[E], setSubChildren BiConsumer[E]) []E {
	var result []E
	for _, x := range list {
		if rootCheck(x) {
			setSubChildren(x, makeChildren(x, list, parentCheck, setSubChildren))
			result = append(result, x)
		}
	}
	return result
}

// makeChildren 递归构建tree
func makeChildren[E any](parent E, allData []E, parentCheck BiFunction[E], setSubChildren BiConsumer[E]) []E {
	var children []E
	for _, x := range allData {
		if parentCheck(parent, x) {
			subChildren := makeChildren(x, allData, parentCheck, setSubChildren)
			setSubChildren(x, subChildren)
			children = append(children, x)
		}
	}
	return children
}

// Flat 将tree打平为list
func Flat[E any](tree []E, getSubChildren Function[E], setSubChildren Consumer[E]) []E {
	var res []E
	forPostOrder(tree, func(item E) {
		setSubChildren(item)
		res = append(res, item)
	}, getSubChildren)
	return res
}

// ForPreOrder 对tree进行前序遍历
func ForPreOrder[E any](tree []E, consumer Consumer[E], getSubChildren Function[E]) {
	for _, l := range tree {
		consumer(l)
		es := getSubChildren(l)
		if len(es) > 0 {
			ForPreOrder(es, consumer, getSubChildren)
		}
	}
}

// ForLevelOrder 对tree进行层序遍历
func ForLevelOrder[E any](tree []E, consumer Consumer[E], getSubChildren Function[E]) {
	queue := append([]E{}, tree...)
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		consumer(item)
		childList := getSubChildren(item)
		if len(childList) > 0 {
			queue = append(queue, childList...)
		}
	}
}

// ForPostOrder 对tree进行后序遍历
func ForPostOrder[E any](tree []E, consumer Consumer[E], getSubChildren Function[E]) {
	forPostOrder(tree, consumer, getSubChildren)
}

// forPostOrder 递归后序遍历
func forPostOrder[E any](tree []E, consumer Consumer[E], getSubChildren Function[E]) {
	for _, item := range tree {
		childList := getSubChildren(item)
		if len(childList) > 0 {
			forPostOrder(childList, consumer, getSubChildren)
		}
		consumer(item)
	}
}

// Sort 根据给定的比较器对tree的所有子节点进行排序
func Sort[E any](tree []E, comparator func(E, E) bool, getChildren Function[E]) []E {
	for _, item := range tree {
		childList := getChildren(item)
		if len(childList) > 0 {
			Sort(childList, comparator, getChildren)
		}
	}
	sort.Slice(tree, func(i, j int) bool {
		return comparator(tree[i], tree[j])
	})
	return tree
}
