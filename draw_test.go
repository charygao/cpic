package cpic

import (
	"log"
	"testing"
)

func TestDraw(t *testing.T) {
	p := NewParser(`tree:
	-> black
		-> red
		-> red`)
	n := p.Parse()
	/*
		n.walk(func(nd *node) {
			log.Println(nd)
		})
	*/
	m := NewMatrix(n)
	//TREE
	//black
	//|  \
	//red red   width=7,height=3
	if m.width != 7 || m.height != 4 {
		log.Fatal(m.width, m.height)
	}
	p = NewParser(`tree:
	-> black
		->red
			->red
		->red
		->red`)
	n = p.Parse()
	//TREE
	//black
	//|  \   \
	//red red red
	//|
	//red    width=11,height=6
	m = NewMatrix(n)
	if m.width != 11 || m.height != 6 {
		log.Fatal(m.width, m.height)
	}

}
