package cpic

import (
	"fmt"
	"testing"
)

func walkFunc(n *node) {
	for i := 0; i < n.depth; i++ {
		fmt.Print(".")
	}
	fmt.Println(n.ele)
}
func TestMain(t *testing.T) {
	root := &node{}
	node3 := &node{ele: 3}
	node2 := &node{ele: 2}
	root.add(node3, LEFT)
	root.add(node2, LEFT)
	node4 := &node{ele: 4}
	r := root.add(node4, RIGHT)
	node5 := &node{ele: 5}
	node6 := &node{ele: 6}
	node7 := &node{ele: 7}
	r.add(node5, NEXT)
	r.add(node6, LEFT)
	r.add(node7, RIGHT)
	root.ele = int(1)
	//fmt.Println(root.height())
	//root.walk(walkFunc)
}
