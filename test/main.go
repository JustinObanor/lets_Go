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

func (b *BinarySearchTree) lookup(value int) bool {
	if b.root == nil {
		return false
	}

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
	if b.root == nil {
		return
	}

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

			//no right child
			if currNode.right == nil {
				if parentNode == nil {
					b.root = currNode.left
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
					b.root = currNode.right
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
					b.root = leftMost
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

/*
			 10
		7		   18
	 2    9     17	  25
*/

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

/*
					10

		7						18
		
2				9		17				25		
*/

func main() {
	var bst BinarySearchTree
	bst.insert(10)
	bst.insert(7)
	bst.insert(18)
	bst.insert(2)
	bst.insert(9)
	bst.insert(17)
	bst.insert(25)

	// bst.inOderTraversal(bst.root)
	bst.preOderTraversal(bst.root)
	// bst.postOderTraversal(bst.root)

}
