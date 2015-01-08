package cpic // import "github.com/ggaaooppeenngg/cpic"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]
import (
	"github.com/ggaaooppeenngg/util/container"
)

type graph struct {
	*container.Graph
}

//getMax returns max number of a and b,if a and b are both smaller thanl,returns l
func getMax(a, b, l int) int {
	if a < l && b < l {
		return l
	} else {
		if a > b {
			return a
		} else {
			return b
		}
	}
}
func newGraph() *graph {
	return &graph{container.NewGraph()}
}

//TODO:graph要有某种合理的排序使得空间不那么拥挤.
//TODO:考虑边上要加上权重的值.
func (g graph) scale() (width int, height int) {
	//确定每个结点矩形的大小
	//矩形的上下边取决于两个入度的最大值,不能小于1.
	//矩形的左右边取决于两个出度的最大值,不能小于1.
	//遍历每个结点
	//整个图形的骨架就是.
	//**
	//**-+
	// ^ |
	// | v
	// | ***
	// | ***
	// +-***-+
	//     ^ |
	//     | v
	//     +-*
	for k1 := 0; k1 < len(g.Vertices); k1++ {
		var v1 = g.Vertices[k1]
		var d, r int // bottom inbound,right outbound
		//来自于下面的入度和右面的出度
		if k1 != len(g.Graph.Vertices)-1 {
			for _, v2 := range g.Graph.Vertices[k1+1:] {
				if g.Graph.EdgeExists(v2, v1) {
					d++
				}
				if g.Graph.EdgeExists(v1, v2) {
					r++
				}
			}
		}
		var t, l int //top inbound,left outbound
		//来自于上面的入度和左边的出度
		for _, v2 := range g.Graph.Vertices[:k1] {
			if g.Graph.EdgeExists(v2, v1) {
				t++
			}
			if g.Graph.EdgeExists(v1, v2) {
				l++
			}
		}
		//accumulate width and height
		if k1 != len(g.Vertices)-1 {
			width += getMax(t, d, 1) + 1
			height += getMax(l, r, 1) + 2
		} else {
			width += getMax(t, d, 1)
			height += getMax(l, r, 1)
		}
	}
	return width, height
}

func (g graph) draw(m *matrix) {
}
