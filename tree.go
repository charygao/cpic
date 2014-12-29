package cpic

import (
	"container/list"
	"fmt"
)

/*
tree
*/

type relation int

//position to add to a node.
const (
	rLEFT relation = iota
	rRIGHT
	rPREV
	rNEXT
	rPAREN
)

type tree struct {
	root *node
}

//dfs
func (n *node) walk(f func(n *node)) {
	if n != nil {
		f(n)
		for l := n.leftChild; l != nil; l = l.nextSibling {
			l.walk(f)
		}
	}
}

//bfs
func (n *node) bfs(f func(n *node)) {
	var queue = list.New()
	queue.PushBack(n)
	for queue.Len() != 0 {
		frt := queue.Front()
		n := frt.Value.(*node)
		queue.Remove(frt)
		f(n)
		for l := n.leftChild; l != nil; l = l.nextSibling {
			queue.PushBack(l)
		}
	}
}

type node struct {
	leftChild   *node
	rightChild  *node
	parent      *node
	prevSibling *node
	nextSibling *node
	depth       int
	ele         interface{}
}

func (n *node) String() string {
	return fmt.Sprintf("%.10s", n.ele)
}

//bfs get widest width
func (n *node) width() int {
	var width = 0
	var queue = list.New()
	queue.PushBack(n)
	for queue.Len() != 0 {
		frt := queue.Front()
		n := frt.Value.(*node)
		queue.Remove(frt)
		count := 0
		for l := n.leftChild; l != nil; l = l.nextSibling {
			count++
			queue.PushBack(l)
		}
		if count > width {
			width = count
		}
	}
	return width
}

//dfs get highest highet
func (n *node) height() int {
	var height = 0
	n.walk(func(n *node) {
		if n.depth > height {
			height = n.depth
		}
	})
	return height
}
func (n *node) rm() {
	//it's not root
	if n.parent != nil {
		if n.nextSibling != nil && n.prevSibling != nil {
			//remove itself and its children
			n.prevSibling.nextSibling = n.nextSibling
			n.nextSibling.prevSibling = n.prevSibling
			return
		}
		//sigle node
		if n.prevSibling == nil && n.nextSibling == nil {
			n.parent.rightChild = nil
			n.parent.leftChild = nil
			return
		}
		//left child
		if n.prevSibling == nil && n.nextSibling != nil {
			n.parent.leftChild = n.nextSibling
			n.nextSibling.prevSibling = nil
		}
		//right child
		if n.prevSibling != nil && n.nextSibling == nil {
			n.parent.rightChild = n.prevSibling
			n.prevSibling.nextSibling = nil
		}
	}
}

//can not add parent
func (n *node) add(newNode *node, typ relation) *node {
	switch typ {
	case rLEFT:
		newNode.parent = n
		if n.leftChild == nil {
			n.leftChild = newNode
			n.rightChild = newNode
		} else {
			newNode.nextSibling = n.leftChild
			n.leftChild.prevSibling = newNode
			n.leftChild = newNode
		}
		newNode.depth = n.depth + 1
		return newNode
	case rRIGHT:
		newNode.parent = n
		if n.rightChild == nil {
			n.leftChild = newNode
			n.rightChild = newNode
		} else {
			newNode.prevSibling = n.rightChild
			n.rightChild.nextSibling = newNode
			n.rightChild = newNode
		}
		newNode.depth = n.depth + 1
		return newNode
	case rPREV:
		//if it's lefChild
		if n.prevSibling == nil {
			n.parent.leftChild = newNode
			newNode.parent = n.parent
			newNode.nextSibling = n
			n.prevSibling = newNode
		} else {
			newNode.parent = n.parent
			newNode.prevSibling = n.prevSibling
			newNode.nextSibling = n
			n.prevSibling.nextSibling = newNode
			n.prevSibling = newNode
		}
		newNode.depth = n.depth
		return newNode
	case rNEXT:
		//if it's right child
		if n.nextSibling == nil {
			n.parent.rightChild = newNode
			newNode.parent = n.parent
			newNode.prevSibling = n
			n.nextSibling = newNode
		} else {
			newNode.parent = n.parent
			newNode.nextSibling = n.nextSibling
			newNode.prevSibling = n
			n.nextSibling.prevSibling = newNode
			n.nextSibling = newNode
		}
		newNode.depth = n.depth
		return newNode
	}
	return newNode

}
