package model

import "fmt"

type BinaryTree struct {
	Key       int
	LeftTree  *BinaryTree
	RightTree *BinaryTree
}

func Find(root *BinaryTree, key int) *BinaryTree {
	if root == nil {
		return nil
	}

	if root.Key == key {
		return root
	} else if root.Key > key {
		return Find(root.LeftTree, key)
	} else {
		return Find(root.RightTree, key)
	}
}

func InsertNode(root *BinaryTree, key int) *BinaryTree {
	if root == nil {
		return &BinaryTree{Key: key}
	}

	if root.Key > key {
		root.LeftTree = InsertNode(root.LeftTree, key)
	} else if root.Key < key {
		root.RightTree = InsertNode(root.RightTree, key)
	} else {
		return root
	}

	return root
}

func findMin(root *BinaryTree) *BinaryTree {
	if root == nil {
		return nil
	}
	for root.LeftTree != nil {
		root = root.LeftTree
	}
	return root
}

func DeleteNode(root *BinaryTree, key int) *BinaryTree {
	if root == nil {
		return nil
	}
	if key < root.Key {
		root.LeftTree = DeleteNode(root.LeftTree, key)
		return root
	} else if key > root.Key {
		root.RightTree = DeleteNode(root.RightTree, key)
		return root
	}

	if root.LeftTree != nil && root.RightTree != nil {
		successor := findMin(root.RightTree)
		root.Key = successor.Key
		root.RightTree = DeleteNode(root.RightTree, successor.Key)
		return root
	}
	if root.LeftTree == nil {
		return root.RightTree
	}
	return root.LeftTree
}

func PreorderPrint(root *BinaryTree) {
	if root == nil {
		return
	}

	fmt.Printf("%d ", root.Key)
	PreorderPrint(root.LeftTree)
	PreorderPrint(root.RightTree)
}

func FindUnbalancedNodes(root *BinaryTree) []int {
	var res []int
	var postorderVisit func(root *BinaryTree) int

	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}

	postorderVisit = func(root *BinaryTree) int {
		if root == nil {
			return 0
		}

		leftCount := postorderVisit(root.LeftTree)
		rightCount := postorderVisit(root.RightTree)

		if abs(leftCount-rightCount) == 1 {
			res = append(res, root.Key)
		}

		return leftCount + rightCount + 1
	}

	postorderVisit(root)
	return res
}

func HeightIterativeDFS(root *BinaryTree) int {
	if root == nil {
		return 0
	}

	type NodeLevel struct {
		Node  *BinaryTree
		Level int
	}

	stack := []*NodeLevel{{
		Node:  root,
		Level: 1,
	}}
	maxLevel := 1

	for len(stack) != 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if curr.Level > maxLevel {
			maxLevel = curr.Level
		}

		if curr.Node.RightTree != nil {
			stack = append(stack, &NodeLevel{
				Node:  curr.Node.RightTree,
				Level: curr.Level + 1,
			})
		}

		if curr.Node.LeftTree != nil {
			stack = append(stack, &NodeLevel{
				Node:  curr.Node.LeftTree,
				Level: curr.Level + 1,
			})
		}
	}

	return maxLevel
}

func KthLeafInorder(root *BinaryTree, k int) (*BinaryTree, bool) {
	counter := 0
	var result *BinaryTree
	var inorderVisit func(root *BinaryTree)

	inorderVisit = func(root *BinaryTree) {
		if root == nil {
			return
		}

		inorderVisit(root.LeftTree)
		if root.LeftTree == nil && root.RightTree == nil {
			counter++
			if counter == k {
				result = root
				return
			}
		}
		inorderVisit(root.RightTree)
	}
	inorderVisit(root)

	if result != nil {
		return result, true
	}

	return nil, false
}
