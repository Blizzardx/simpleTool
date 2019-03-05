package main

import "fmt"

//----------------stack----------------
type StackNode struct {
	Data interface{}
	Next *StackNode
}
type Stack struct {
	Top *StackNode
}

func (stack *Stack) Push(data interface{}) {
	node := &StackNode{Data: data}
	node.Next = stack.Top
	stack.Top = node
}
func (stack *Stack) Pop() (interface{}, bool) {
	if stack.Top == nil {
		return nil, false
	}
	node := stack.Top
	stack.Top = node.Next
	return node.Data, true
}

//----------------queue----------------
type QueueNode struct {
	Data interface{}
	Next *QueueNode
	Pre  *QueueNode
}
type Queue struct {
	First *QueueNode
	Last  *QueueNode
}

func (queue *Queue) Enqueue(data interface{}) {
	nodeData := &QueueNode{Data: data}
	currentFirst := queue.First
	queue.First = nodeData
	nodeData.Next = currentFirst
	nodeData.Pre = nil
	if currentFirst == nil {
		queue.Last = nodeData
	} else {
		currentFirst.Pre = nodeData
	}
}
func (queue *Queue) Dequeue() (interface{}, bool) {
	nodeData := queue.Last
	if nodeData == nil {
		return nil, false
	}
	queue.Last = nodeData.Pre
	if queue.Last != nil {
		queue.Last.Next = nil
	} else {
		queue.First = nil
	}

	return nodeData.Data, true
}

//----------------tree----------------
type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Data  int
}
type Tree struct {
	Root *TreeNode
}

func printTree(node *TreeNode, orderId int) {
	if node == nil {
		return
	}
	//
	if orderId == 0 {
		printTreeData(node)
	}
	printTree(node.Left, orderId)

	if orderId == 1 {
		printTreeData(node)
	}
	printTree(node.Right, orderId)

	if orderId == 2 {
		printTreeData(node)
	}
}
func printTreeData(node *TreeNode) {
	if node == nil {
		return
	}
	fmt.Print(node.Data, " ")
}
func reversTree(node *TreeNode) {
	if node == nil {
		return
	}
	tmpNode := node.Left
	node.Left = node.Right
	node.Right = tmpNode
	reversTree(node.Left)
	reversTree(node.Right)
}

func (tree *Tree) reversTree() {
	reversTree(tree.Root)
}
func (tree *Tree) printTree(orderId int) {
	printTree(tree.Root, orderId)
}
func (tree *Tree) printTreeByQueue() {
	tmpQueue := &Queue{}
	tmpNode := tree.Root
	for tmpNode != nil {
		if tmpNode.Left != nil {
			tmpQueue.Enqueue(tmpNode.Left)
		}
		if tmpNode.Right != nil {
			tmpQueue.Enqueue(tmpNode.Right)
		}
		printTreeData(tmpNode)

		tmpData, ok := tmpQueue.Dequeue()
		if !ok {
			tmpNode = nil
		} else {
			tmpNode = tmpData.(*TreeNode)
		}
	}
}
func (tree *Tree) printTreeByStack() {
	tmpStack := &Stack{}
	tmpNode := tree.Root
	for tmpNode != nil {
		if tmpNode.Right != nil {
			tmpStack.Push(tmpNode.Right)
		}
		if tmpNode.Left != nil {
			tmpStack.Push(tmpNode.Left)
		}
		printTreeData(tmpNode)

		tmpData, ok := tmpStack.Pop()
		if !ok {
			tmpNode = nil
		} else {
			tmpNode = tmpData.(*TreeNode)
		}
	}
}

//----------------test----------------
func createTestTree() *Tree {
	node0 := &TreeNode{Data: 0}
	node1 := &TreeNode{Data: 1}
	node2 := &TreeNode{Data: 2}
	node3 := &TreeNode{Data: 3}
	node4 := &TreeNode{Data: 4}
	node5 := &TreeNode{Data: 5}
	node6 := &TreeNode{Data: 6}
	node7 := &TreeNode{Data: 7}
	node8 := &TreeNode{Data: 8}
	node9 := &TreeNode{Data: 9}
	node10 := &TreeNode{Data: 10}
	node11 := &TreeNode{Data: 11}
	node12 := &TreeNode{Data: 12}
	node13 := &TreeNode{Data: 13}
	node14 := &TreeNode{Data: 14}
	node15 := &TreeNode{Data: 15}
	node16 := &TreeNode{Data: 16}
	node17 := &TreeNode{Data: 17}

	node0.Left = node1
	node0.Right = node2
	node1.Left = node3
	node1.Right = node4
	node2.Left = node5
	node2.Right = node6
	node3.Left = node7
	node3.Right = node8
	node4.Left = node9
	node5.Left = node10
	node5.Right = node11
	node6.Right = node12
	node9.Left = node13
	node9.Right = node14
	node11.Left = node15
	node14.Left = node16
	node14.Right = node17

	tree := &Tree{Root: node0}

	return tree
}

func main() {
	tree := createTestTree()
	tree.printTree(0)
	fmt.Println("")
	fmt.Println("------------------------------------")
	tree.printTree(1)
	fmt.Println("")
	fmt.Println("------------------------------------")
	tree.printTree(2)
	fmt.Println("")
	fmt.Println("------------------------------------")
	tree.printTreeByStack()
	fmt.Println("")
	fmt.Println("------------------------------------")
	tree.printTreeByQueue()
	fmt.Println("")
	fmt.Println("------------------------------------")
	tree.reversTree()
	tree.printTreeByQueue()
	fmt.Println("")
	fmt.Println("------------------------------------")
}
