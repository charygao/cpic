package cpic

import (
	"fmt"
	//"log"
	"testing"
)

func TestDraw(t *testing.T) {
	/*
		p := NewParser(`tree:
		-> black
			-> red
			-> red`)
		n := p.Parse()
	*/
	//m := NewMatrix(n)
	//TREE
	//black
	//|  \
	//red red   width=7,height=3
	/*
		if m.width != 7 || m.height != 4 {
			log.Fatal(m.width, m.height)
		}
		m.draw()
		log.Println(m.output())
	*/
	p := NewParser(`tree:
	-> black
		->red
			->red
		->red`)
	n := p.Parse()
	var count int
	n.walk(func(nd *node) {
		for i := 0; i < nd.depth; i++ {
			fmt.Print("*")
		}
		count++
		if nd.ele != nil {
			fmt.Println(count, nd.ele.(Token))
		} else {
			fmt.Println(count, "TREE")
		}
	})
	//TREE
	//black
	//|  \   \
	//red red red
	//|
	//red    width=11,height=6
	/*
		m = NewMatrix(n)
		if m.width != 11 || m.height != 6 {
			log.Fatal(m.width, m.height)
		}
		m.draw()
		log.Println(m.output())
	*/

}
