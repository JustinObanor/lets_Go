package main

import (
	"fmt"
)

//Node ...
type Node struct {
	value int
	left  *Node
	right *Node
}

//BinarySearchTree ...
type BinarySearchTree struct {
	root *Node
}

func newNode(value int) *Node {
	return &Node{
		value: value,
	}
}

func (bst *BinarySearchTree) insert(value int) {
	node := newNode(value)
	if bst.root == nil {
		bst.root = node
	} else {
		currNode := bst.root
		for {
			if value > currNode.value {
				if currNode.right == nil {
					currNode.right = node
					return
				}
				currNode = currNode.right
			} else {
				if currNode.left == nil {
					currNode.left = node
					return
				}
				currNode = currNode.left
			}
		}
	}
}

//binary search(Olog(n))
func (bst *BinarySearchTree) lookup(value int) bool {
	if bst.root == nil {
		return false
	}

	currNode := bst.root

	for currNode != nil {
		if value < currNode.value {
			currNode = currNode.left
		} else if value > currNode.value {
			currNode = currNode.right
		} else if value == currNode.value {
			return true
		}
	}
	return false
}

func (bst *BinarySearchTree) remove(value int) {
	if bst.root == nil {
		return
	}

	currNode := bst.root

	parentNode := currNode
	for currNode != nil {
		if value < currNode.value {
			parentNode = currNode
			currNode = currNode.left
		} else if value > currNode.value {
			parentNode = currNode
			currNode = currNode.right
		} else if value == currNode.value {

			//no right child
			if currNode.right == nil {
				if parentNode == nil {
					bst.root = currNode.left
				} else {
					if currNode.value > parentNode.value {
						parentNode.right = currNode.left
					} else if currNode.value < parentNode.value {
						parentNode.left = currNode.left
					}
				}

				//right child which doesnt have a left child
			} else if currNode.right.left == nil {
				currNode.right.left = currNode.left
				if parentNode == nil {
					bst.root = currNode.right
				} else {
					if currNode.value < parentNode.value {
						parentNode.left = currNode.right
					} else if currNode.value > parentNode.value {
						parentNode.right = currNode.right
					}
				}

				//Right child that has a left child
			} else {
				//find the Right child's left most child

				leftMost := currNode.right.left
				leftMostParent := currNode.right
				for leftMost != nil {
					leftMostParent = leftMost
					leftMost = leftMost.left
				}

				//basically replacing currNode w/ leftMosts

				leftMostParent.left = leftMost.right
				leftMost.left = currNode.left
				leftMost.right = currNode.right

				if parentNode == nil {
					bst.root = leftMost
				} else {
					if currNode.value < parentNode.value {
						parentNode.left = leftMost
					} else if currNode.value > parentNode.value {
						parentNode.right = leftMost
					}
				}
			}
		}
	}
}

func (bst *BinarySearchTree) preOderTraversal(currNode *Node) {
	if currNode == nil {
		return
	}

	fmt.Println(currNode.value)

	bst.preOderTraversal(currNode.left)

	bst.preOderTraversal(currNode.right)
}

func (bst *BinarySearchTree) inOderTraversal(currNode *Node) {
	if currNode == nil {
		return
	}

	bst.inOderTraversal(currNode.left)

	fmt.Println(currNode.value)

	bst.inOderTraversal(currNode.right)
}

func (bst *BinarySearchTree) postOderTraversal(currNode *Node) {
	if currNode == nil {
		return
	}

	bst.postOderTraversal(currNode.left)

	bst.postOderTraversal(currNode.right)

	fmt.Println(currNode.value)
}

func (bst *BinarySearchTree) breadthFirstSearch() (list []int) {
	queue := []*Node{}

	currNode := bst.root
	queue = append(queue, currNode)
	for len(queue) > 0 {
		currNode, queue = queue[0], queue[1:]
		list = append(list, currNode.value)

		if currNode.left != nil {
			queue = append(queue, currNode.left)
		}

		if currNode.right != nil {
			queue = append(queue, currNode.right)
		}
	}

	return list
}

func (bst *BinarySearchTree) breadthFirstSearchRecursive(queue []*Node, list []int) []int {
	if len(queue) == 0 {
		return list
	}

	currNode, queue := queue[0], queue[1:]
	list = append(list, currNode.value)

	if currNode.left != nil {
		queue = append(queue, currNode.left)
	}

	if currNode.right != nil {
		queue = append(queue, currNode.right)
	}

	return bst.breadthFirstSearchRecursive(queue, list)
}

/*
				9
		4				20
	1		6		15		170
bfs = 9 4 20 1 6 15 170 (more memory required)
dfs = 9 4 1 6 20 15 170 (lower memory required)
*/

func main() {
	var bst BinarySearchTree
	bst.insert(9)
	bst.insert(4)
	bst.insert(6)
	bst.insert(20)
	bst.insert(170)
	bst.insert(15)
	bst.insert(1)
	fmt.Println(bst.breadthFirstSearch())
	fmt.Println(bst.breadthFirstSearchRecursive([]*Node{bst.root}, []int{}))
}

//If you know a solution is not far from the root of the tree:BFS

//If the tree is very deep and solutions are rare: BFS (dfs too slow(recursive functions))

//If the tree is very wide:DFS(BFS will take too much memory)

//If solutions are frequent but located deep in the tree:DFS

//Determining whether a path exists between two nodes:DFS

//Finding the shortest path:BFS
