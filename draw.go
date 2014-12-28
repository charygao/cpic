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
	node   *node
	py     [][]byte // v y
}

//这个应该还是一个分治问题
func getscale(n *node) (width int, height int) {
	if n.ele == nil {
		height = 0
		width = 4 // len("TREE")
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

func NewMatrix(n *node) *matrix {

	m := new(matrix)
	m.node = n
	m.height = 0
	m.width = 0

	//这个应该还是一个分治问题
	var getscale func(n *node) (width int, height int)
	getscale = func(n *node) (width int, height int) {
		if n.ele == nil {
			height = 0
			width = 4 // len("TREE")
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
	m.py = make([][]byte, m.height)
	for i := 0; i < m.height; i++ {
		m.py[i] = make([]byte, m.width)
	}
	for i := 0; i < len(m.py); i++ {
		for j := 0; j < len(m.py[i]); j++ {
			m.py[i][j] = ' '
		}
	}
	return m
}

//paint (x,y)
func (m *matrix) paint(r byte, x, y int) {
	m.py[y][x] = r
}

//paint an area
func (m *matrix) paintA(content [][]byte, px, py int) {
	for y, xl := range content {
		for x, r := range xl {
			m.py[py+y][px+x] = r
		}
	}
}

//递归绘制TREE
func (m *matrix) draw() {
	var _draw func(n *node, offsetX, offsetY int)
	_draw = func(n *node, offsetX, offsetY int) {
		fmt.Println("draw x,y", offsetX, offsetY)
		//根节点打印TREE
		if n.ele == nil {
			m.paintA([][]byte{[]byte("TREE")}, offsetX, offsetY)
			offsetY++
			if n.leftChild != nil {
				_draw(n.leftChild, offsetX, offsetY)
			}
		} else {
			fmt.Println("is not nil")
			m.paintA([][]byte{[]byte(n.ele.(Token).lit)}, offsetX, offsetY)
			offsetY++
			//如果有子节点就要向下画一层辅助线
			if n.leftChild != nil {
				fmt.Println("child is nil")
				m.paint(byte('|'), offsetX, offsetY)
				fmt.Println("draw right down child  x y", offsetX, offsetY+1)
				_draw(n.leftChild, offsetX, offsetY+1)
				fmt.Println("get out right down child")
				w, _ := getscale(n.leftChild)
				fmt.Println("child width->", w)
				offsetX += w
				for l := n.leftChild.nextSibling; l != nil; l = l.nextSibling {
					fmt.Println("draw \\ -> x,y:", offsetX, offsetY)
					m.paint(byte('\\'), offsetX, offsetY)
					_draw(n.leftChild, offsetX+1, offsetY+1)
					w, _ := getscale(n.leftChild)
					offsetX += w + 1
				}
			}

		}
	}
	_draw(m.node, 0, 0)
}

//输出字符串
func (m *matrix) output() string {
	for _, xl := range m.py {
		for _, p := range xl {
			fmt.Printf("%c", p)
		}
		fmt.Println()
	}
	return ""
}
