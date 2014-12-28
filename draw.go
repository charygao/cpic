package cpic

import (
	"fmt"
)

func _placeholder() {
	fmt.Println("place_holder")
}

//> x , v y
type matrix struct {
	width  int
	height int
	py     [][]uint8 // v y
}

func NewMatrix(n *node) *matrix {
	m := new(matrix)
	m.height = n.height()
	m.width = n.width()
	return m
}

//paint (x,y)
func (m *matrix) paint(r uint8, x, y int) {
	if x >= m.width || x < 0 {
		panic("x out")
	}
	if y >= m.height || y < 0 {
		panic("y out")
	}
	m.py[y][x] = r
}

//paint an area
func (m *matrix) paintA(content [][]uint8, px, py int) {
	for y, xl := range content {
		for x, r := range xl {
			m.py[py+y][px+x] = r
		}
	}
}

func (m *matrix) draw() {
	for _, xl := range m.py {
		fmt.Printf("%s\n", xl)
	}
}
