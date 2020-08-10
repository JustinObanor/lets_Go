package main

import "fmt"

type Node struct {
	value int
	left  *Node
	right *Node
}

type BinarySearchTree struct {
	root *Node
}

func newNode(value int) *Node {
	return &Node{
		value: value,
	}
}

func (b *BinarySearchTree) insert(value int) {
	node := newNode(value)
	if b.root == nil {
		b.root = node
	} else {
		currNode := b.root
		for {
			if value < currNode.value {
				if currNode.left == nil {
					currNode.left = node
					return
				}
				currNode = currNode.left
			} else {
				if currNode.right == nil {
					currNode.right = node
					return
				}
				currNode = currNode.right
			}
		}
	}
}

func (b *BinarySearchTree) lookup(value int) bool {
	currNode := b.root
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

func (b *BinarySearchTree) remove(value int) {
	currNode := b.root
	parentNode := currNode

	for currNode != nil {
		if value < currNode.value {
			parentNode = currNode
			currNode = currNode.left
		} else if value > currNode.value {
			parentNode = currNode
			currNode = currNode.right
		} else if value == currNode.value {
			if currNode.right == nil {
				if parentNode == nil {
					b.root = currNode
				} else {
					if parentNode.value < currNode.value {
						parentNode.right = currNode.left
					} else {
						parentNode.left = currNode.left
					}
				}

			} else if currNode.right.left == nil {
				if parentNode == nil {
					b.root = currNode
				} else {
					if parentNode.value < currNode.value {
						parentNode.right = currNode.right
					} else {
						parentNode.left = currNode.right
					}
				}

			} else {
				leftMost := currNode.right.left
				leftMostParent := currNode.right

				for leftMost != nil {
					leftMostParent = leftMost
					leftMost = leftMost.left
				}

				leftMostParent.left = leftMost.right
				leftMost.left = currNode.left
				leftMost.right = currNode.right

				if parentNode == nil {
					b.root = leftMost
				} else {
					if parentNode.value < currNode.value {
						parentNode.right = leftMost
					} else {
						parentNode.left = leftMost
					}
				}
			}

		}
	}
}

func main() {
	var bst BinarySearchTree
	bst.insert(9)
	bst.insert(7)
	bst.insert(10)
	bst.insert(6)
	bst.insert(8)
	fmt.Println(bst.lookup(6))
}
