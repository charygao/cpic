package cpic

//二维树结构的字符画绘制方法:
//开头第一行是特殊行,标记"TREE",对应root结点,ele是nil.
//之后可以把每个树和子树都抽象成一个矩形,高度是树的高度,宽度是广度遍历的最宽宽度.
//绘制一个结点时如果有子结点的时候,第一个子节点添加一条辅助线('/'),然后空出该子树宽度的空格,
//当有第二个或者更多子树的时候,之间要用空格隔开一个,对应的位置作辅助线('\'),之后一个位置开始绘制子树,
//然后递归绘制.
//结构类似类似于下面这张结构图
//修改:这种方式,在表现二叉树的时候无法体现左右子树的性质,第一个子树为左斜线比较合适,可以体现出左子树的性质.
//
//
//TREE
//┌──┐
//└──┘
///   ╲      ╲
//┌──┐ ┌────┐ ┌────────────┐
//│  │ │    │ │            │
//│  │ │    │ │            │
//│  │ └────┘ └────────────┘
//└──┘
//
//树的宽度是子树之和+n-1(n为子树的个数)

//图的绘制思路:
//线的交叉用+表示,只用直线,+表示转弯,如果线和线之间交叉后画的盖过前画的.
//不推荐生成比较复杂的图,生成出来的图太容易覆盖掉了.
//首先选择一个方向拓扑排序生成对应的位置,一次从左上到左下排开.
//当结点数是3的时候,以右边为例,左边是对称的,1上面空出一行,3的右边空出一行.
//_----+
//1-+  |
//  |  |
//  v  |
//  2-+|
//    ||
//	  vv
//	  3_
//当结点数是4的时候,以右边的为例,每次向下上面的横线少一个,右边的横线多一个.但是本身算一个横线也算一个竖线,如果要算位置的话.
//最近的在最左边,还没想清楚用自顶向下还是自底向上.
// ___--------+
// ___-----+  |
// __1-+   |  |
// ^^^ |   |  |           //TODO:总结出这种画法得出的规模.
// ||| v   |  |
// |||___----+|
// ||+_2_-+| ||
// || ___ || ||
// || ^^  vv ||
// || |+- 3 +||
// |+----   |||
// |  |  ^  |||
// |  |  |  vvv
// |  |  +--4__
// |  +-----___
// +--------___
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

type painter interface {
	draw(m *matrix)
	scale() (width int, height int)
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
	node   painter  //root node
	py     [][]byte //every py's element is a x line
}

//Newmatrix parses root node scale's width and height to
//initate underlying [][]uint8 and return matrix.
func newMatrix(node painter) *matrix {
	m := new(matrix)
	m.node = node

	m.width, m.height = node.scale()
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

//paintA paints an area(px,py,lenX,lenY)
func (m *matrix) paintA(content [][]byte, px, py int) {
	for y, xl := range content {
		for x, r := range xl {
			m.py[py+y][px+x] = r
		}
	}
}

//paintT paints a transpose matrix
func (m *matrix) paintT(content [][]byte, px, py int) {
	var transpose [][]byte
	if len(content) != 0 {
		transpose = make([][]byte, len(content[0]))
	}
	for i, _ := range transpose {
		transpose[i] = make([]byte, len(content))
	}
	for i, bs := range content {
		for j, byte := range bs {
			transpose[j][i] = byte
		}
	}
	m.paintA(transpose, px, py)
}

//递归绘制图形
//paint tree on the matrix.
func (m *matrix) draw() {
	m.node.draw(m)
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
	t := newTree()
	t.root = n
	m := newMatrix(t)
	m.draw()
	return m.output(), nil
}
