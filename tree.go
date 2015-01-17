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

//implements painter
type tree struct {
	root *node
}

func newTree() *tree {
	return &tree{
		root: new(node),
	}
}

func (t tree) scale() (width int, height int) {
	return getscale(t.root)
}
func _draw(m *matrix, n *node, offsetX, offsetY int) {
	//fmt.Println("draw x,y", offsetX, offsetY)
	//根节点打印TREE
	//print "TREE" for root nil node.
	if n.ele == nil {
		//m.paintA([][]byte{[]byte("TREE")}, offsetX, offsetY)
		//offsetY++
		if n.leftChild != nil {
			_draw(m, n.leftChild, offsetX, offsetY)
		}
	} else {
		//fmt.Println("is not nil")
		m.paintA([][]byte{[]byte(n.ele.(token).lit)}, offsetX, offsetY)
		offsetY++
		//如果有子节点就要向下画一层辅助线
		//if there is child,paint  auxiliary line.
		if n.leftChild != nil {
			//fmt.Println("child is not nil", "parent is", n.ele.(token))
			m.paint(byte('/'), offsetX, offsetY)
			//fmt.Println("draw right down child  x y", offsetX, offsetY+1)
			_draw(m, n.leftChild, offsetX, offsetY+1)
			//fmt.Println("get out right down child")
			w, _ := getscale(n.leftChild)
			//fmt.Println("child width->", w)
			//fmt.Println("before offsetX is", offsetX)
			offsetX += w
			//fmt.Println("begin loop****")
			for l := n.leftChild.nextSibling; l != nil; l = l.nextSibling {
				//fmt.Println("draw \\ -> x,y:", offsetX, offsetY)
				m.paint(byte('\\'), offsetX, offsetY)
				_draw(m, l, offsetX+1, offsetY+1)
				w, _ = getscale(l)
				offsetX += w + 1
				//fmt.Println("after offsetX is", offsetX)
			}
			//fmt.Println("end loop****")
		}

	}
}

func (t tree) draw(m *matrix) {
	_draw(m, t.root, 0, 0)
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

//getscale gets width and hieght of a root node n on behalf of a tree.
func getscale(n *node) (width int, height int) {
	if n.ele == nil {
		//height = 1
		//width = 4 // len("TREE")
	} else {
		height = 1
		width = len(n.ele.(token).lit)
	}
	if n.leftChild != nil {
		//该节点右子节点至少要画一个竖线,高度还要加1,root的NIL结点不加
		//one more line for auxiliary line( '\' | '|'),root的NIL结点不加
		if n.ele != nil {
			height++
		}
		var htTmp, whTmp, withSpaceCount int

		for l := n.leftChild; l != nil; l = l.nextSibling {
			//递归获取每个子树的规模,取最高的高度加到高度的规模上.
			//宽度要叠加,还要算上子树之间的一个空白符,最后和节点本身的宽度比较,决定大小.
			//get highest hight of a tree to represent the whole tree.
			//accumulate all the width of nodes on the same level as well as withspace between.
			w, h := getscale(l)
			if h > htTmp {
				htTmp = h
			}
			whTmp += w
			withSpaceCount++
		}
		whTmp += withSpaceCount - 1
		if whTmp > width {
			width = whTmp
		}
		height += htTmp
	}
	return
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
