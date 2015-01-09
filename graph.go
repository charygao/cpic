package cpic // import "github.com/ggaaooppeenngg/cpic"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]
import (
	"bytes"

	"github.com/ggaaooppeenngg/util"
	"github.com/ggaaooppeenngg/util/container"
)

type graph struct {
	*container.Graph
	width  map[*container.Vertex]int
	height map[*container.Vertex]int
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
	return &graph{container.NewGraph(), make(map[*container.Vertex]int), make(map[*container.Vertex]int)}
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
		w := getMax(t, d, 1)
		g.width[v1] = w
		h := getMax(l, r, 1)
		g.height[v1] = h
		if k1 != len(g.Vertices)-1 {
			width += w + 1
			height += h + 2
		} else {
			width += w
			height += h
		}
	}
	return width, height
}

func sort(vtxg, vtxs []*container.Vertex, l, h int) {
	if 1 == h-l || vtxs == nil {
		return
	}
	var pivot = vtxs[l]
	var index = l
	for i := l; i < h; i++ {
		if util.IndexOf(vtxs[i], vtxg) < util.IndexOf(pivot, vtxg) {
			vtxs[index+1], vtxs[i] = vtxs[i], vtxs[index+1]
			index++
		}
	}
	sort(vtxg, vtxs, l, index+1)
	sort(vtxg, vtxs, index+1, l)
}

//获得左上角的位置.
func (g graph) pos(index int) (x, y int) {
	for _, v := range g.Graph.Vertices[:index] {
		x += g.width[v] + 1  //空格
		y += g.height[v] + 2 //下箭头
	}
	return
}

func (g graph) draw(m *matrix) {
	//一趟向下
	for index1, v1 := range g.Graph.Vertices {
		outbound := g.Graph.OutBound(v1)
		sort(g.Graph.Vertices, outbound, 0, len(outbound))
		if index1 != len(g.Graph.Vertices)-1 {
			outbound = outbound[index1+1:]
		}
		//右边的出度,最接近,最下面
		for k1, v2 := range outbound {
			//第k1个出度
			index2 := util.IndexOf(g.Graph.Vertices, v2)
			//起点.
			var xS, yS, xR, yR, xE, yE, x, y int //start,end,temporary position
			x, y = g.pos(index1)
			xS = x + g.width[v1]
			yS = y + g.height[v1] - k1
			inbound := g.Graph.InBound(v2)
			sort(g.Graph.Vertices, inbound, 0, len(inbound))
			//对于v2,v1的位置.,最后得出v1到v2的路线.
			k2 := util.IndexOf(inbound, v1)
			x, y = g.pos(index2)
			xE = x + g.width[v2] - k2
			yE = y - 1
			yR = yS
			xR = xE
			l := xR - xS
			//画"---"线
			m.paintA([][]byte{bytes.Repeat([]byte("-"), l)}, xS, yS)
			h := yE - yR + 1
			//画"+-->"线,并且转置成竖线
			m.paintT([][]byte{append([]byte("+"), append(bytes.Repeat([]byte("-"), h-2), byte('v'))...)}, xR, yR)
		}
	}
}
