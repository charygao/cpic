package cpic

//二维树结构的字符画绘制方法:
//开头第一行是特殊行,标记"TREE",对应root结点,ele是nil.
//之后可以把每个树和子树都抽象成一个矩形,高度是树的高度,宽度是广度遍历的最宽宽度.
//绘制一个结点时如果有子结点的时候,第一个子节点添加一条辅助线('|'),然后空出该子树宽度的空格,
//当有第二个或者更多子树的时候,之间要用空格隔开一个,对应的位置作辅助线('\'),之后一个位置开始绘制子树,
//然后递归绘制.
//结构类似类似于下面这张结构图
//
//TREE
//┌──┐
//└──┘
//│   ╲      ╲
//┌──┐ ┌────┐ ┌────────────┐
//│  │ │    │ │            │
//│  │ │    │ │            │
//│  │ └────┘ └────────────┘
//└──┘
//
//
import (
	"errors"
	"fmt"
	"strings"
)

var (
	placeHolder byte = ' ' //for debug
)

func _placeholder() {
	fmt.Println("place_holder")
}

//2D matrix
//-------> X
//|
//|
//v
//Y
type matrix struct {
	width  int
	height int
	node   *node    //root node
	py     [][]byte //every py's element is a x line
}

//getscale gets width and hieght of a root node n on behalf of a tree.
func getscale(n *node) (width int, height int) {
	if n.ele == nil {
		height = 1
		width = 4 // len("TREE")
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

//Newmatrix parses root node scale's width and height to
//initate underlying [][]uint8 and return matrix.
func newMatrix(n *node) *matrix {

	m := new(matrix)
	m.node = n
	m.height = 0
	m.width = 0

	m.width, m.height = getscale(n)
	//init m.py

	m.py = make([][]byte, m.height)
	for i := 0; i < m.height; i++ {
		m.py[i] = make([]byte, m.width)
	}
	//set placHolder for every point
	for i := 0; i < len(m.py); i++ {
		for j := 0; j < len(m.py[i]); j++ {
			m.py[i][j] = placeHolder
		}
	}
	return m
}

//paint (x,y)
func (m *matrix) paint(r byte, x, y int) {
	m.py[y][x] = r
}

//paint an area(px,py,lenX,lenY)
func (m *matrix) paintA(content [][]byte, px, py int) {
	for y, xl := range content {
		for x, r := range xl {
			m.py[py+y][px+x] = r
		}
	}
}

//递归绘制TREE
//paint tree on the matrix.
func (m *matrix) draw() {
	var _draw func(n *node, offsetX, offsetY int)
	_draw = func(n *node, offsetX, offsetY int) {
		//fmt.Println("draw x,y", offsetX, offsetY)
		//根节点打印TREE
		//print "TREE" for root nil node.
		if n.ele == nil {
			m.paintA([][]byte{[]byte("TREE")}, offsetX, offsetY)
			offsetY++
			if n.leftChild != nil {
				_draw(n.leftChild, offsetX, offsetY)
			}
		} else {
			//fmt.Println("is not nil")
			m.paintA([][]byte{[]byte(n.ele.(token).lit)}, offsetX, offsetY)
			offsetY++
			//如果有子节点就要向下画一层辅助线
			//if there is child,paint  auxiliary line.
			if n.leftChild != nil {
				//fmt.Println("child is not nil", "parent is", n.ele.(token))
				m.paint(byte('|'), offsetX, offsetY)
				//fmt.Println("draw right down child  x y", offsetX, offsetY+1)
				_draw(n.leftChild, offsetX, offsetY+1)
				//fmt.Println("get out right down child")
				w, _ := getscale(n.leftChild)
				//fmt.Println("child width->", w)
				//fmt.Println("before offsetX is", offsetX)
				offsetX += w
				//fmt.Println("begin loop****")
				for l := n.leftChild.nextSibling; l != nil; l = l.nextSibling {
					//fmt.Println("draw \\ -> x,y:", offsetX, offsetY)
					m.paint(byte('\\'), offsetX, offsetY)
					_draw(l, offsetX+1, offsetY+1)
					w, _ = getscale(l)
					offsetX += w + 1
					//fmt.Println("after offsetX is", offsetX)
				}
				//fmt.Println("end loop****")
			}

		}
	}
	_draw(m.node, 0, 0)
}

//输出字符串
//output returns matrix as a string.
func (m *matrix) output() string {
	var o string
	for _, xl := range m.py {
		o += fmt.Sprintf("%s\n", xl)
	}
	return o
}

//Gen inputs source text and output data structure char picture.
func Gen(src string) (string, error) {
	p := newParser(src)
	n := p.parse()
	if n == nil {
		return "", errors.New(strings.Join(p.errors, "\n"))
	}
	m := newMatrix(n)
	m.draw()
	return m.output(), nil
}
