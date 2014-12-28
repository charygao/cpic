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
	m.height = 0
	m.width = 0

	//这个应该还是一个分治问题
	var getscale func(n *node) (width int, height int)
	getscale = func(n *node) (width int, height int) {
		fmt.Println("get scale")
		fmt.Println(n.ele)
		if n.ele == nil {
			height = 0
			width = 0
		} else {
			height = 1
			width = len(n.ele.(Token).lit)
		}
		if n.leftChild != nil {
			//该节点右子节点至少要画一个竖线,高度还要加1,NIL的根节点加一用来放一个tree标记
			height++
			var htTmp, whTmp, withSpaceCount int

			for l := n.leftChild; l != nil; l = l.nextSibling {
				//递归获取每个子树的规模,取最高的高度加到高度的规模上.
				//宽度要叠加,还要算上子树之间的一个空白符,最后和节点本身的宽度比较,决定大小.
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
		//如果没有子节点,那么宽度就是这个节点的宽度,高度就是1
		return
	}
	m.width, m.height = getscale(n)
	m.py = make([][]uint8, m.height)
	for i := 0; i < m.height; i++ {
		m.py[i] = make([]uint8, m.width)
	}
	return m
}

//paint (x,y)
func (m *matrix) paint(r uint8, x, y int) {
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
