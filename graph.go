package cpic // import "github.com/ggaaooppeenngg/cpic"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]
import (
	"bytes"
	"fmt"

	"github.com/ggaaooppeenngg/util"
	"github.com/ggaaooppeenngg/util/container"
)

type graph struct {
	*container.Graph                           //Graph container
	width            map[*container.Vertex]int //width of vertex in this graph
	height           map[*container.Vertex]int //height of vertex in this graph
	index            map[*container.Vertex]int //index of vertex in Graph.Vertecies
}

//getMax returns max number of a and b,if a and b are both smaller thanl,returns l.
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

//newGraph returns new graph and initiates the maps.
func newGraph() *graph {
	return &graph{container.NewGraph(),
		make(map[*container.Vertex]int),
		make(map[*container.Vertex]int),
		make(map[*container.Vertex]int)}
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
	fmt.Println("scale", width, height)
	return width, height
}

func (g *graph) sort() {
	for i, k := range g.Graph.Vertices {
		g.index[k] = i
	}
}

func sort(vtxg, vtxs []*container.Vertex, l, h int) {
	if h-l < 2 || vtxs == nil {
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
	vtxs[index], vtxs[l] = vtxs[l], vtxs[index]
	sort(vtxg, vtxs, l, index+1)
	sort(vtxg, vtxs, index+1, h)
}

func (g graph) divide(vtx *container.Vertex, vtxs []*container.Vertex) (up, down []*container.Vertex) {
	sort(g.Graph.Vertices, vtxs, 0, len(vtxs))
	divider := -1
	for i, v := range vtxs {
		if g.index[v] > g.index[vtx] {
			divider = i
			break
		}
	}
	if divider == -1 {
		divider = len(vtxs)
	} else {
		down = vtxs[divider:]
	}
	up = vtxs[:divider]
	fmt.Println("divider", divider)
	return
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
	g.sort()
	fmt.Println("draw")
	//一趟向下
	for index1, v1 := range g.Graph.Vertices {
		fmt.Println("loop")
		fmt.Println(m.output())
		outbound := g.Graph.OutBound(v1)
		//divide to down and up outbound
		//get divider
		//对于每个点的出度,都划分成上下两部分,
		//上部分是在全局数组g.Graph.Vertices中在v1,前面的点
		//下部分是在全局数组里面
		outboundUp, outboundDown := g.divide(v1, outbound)

		var index2 int                       // index of outbound vertext
		var xS, yS, xR, yR, xE, yE, x, y int //start,end,temporary position
		for k1, v2 := range outboundDown {
			fmt.Println("down")
			x, y = g.pos(index1)
			index2 = g.index[v2]
			//右边的出度,最接近,最下面
			//第k1个出度
			//起点.
			xS = x + g.width[v1]
			yS = y + g.height[v1] - 1 - k1 //k1 counts from 0
			fmt.Println(yS, y, g.height[v1], k1)
			inbound := g.Graph.InBound(v2)
			up, _ := g.divide(v1, inbound)
			//对于v2,v1的位置.,最后得出v1到v2的路线.
			k2 := util.IndexOf(v1, up)
			x, y = g.pos(index2)
			xE = x + g.width[v2] - 1 - k2 //k2 counts from 0
			fmt.Println(xE, x, g.width[v2], k2)
			yE = y - 1
			yR = yS
			xR = xE
			l := xR - xS
			//画"---"线
			fmt.Println("paint")
			fmt.Println("paint -", l, "times", "start and end", xS, yS)
			m.paintA([][]byte{bytes.Repeat([]byte("-"), l)}, xS, yS)
			h := yE - yR + 1
			//画"+-->"线,并且转置成竖线
			m.paintT([][]byte{append([]byte("+"), append(bytes.Repeat([]byte("|"), h-2), byte('v'))...)}, xR, yR)
		}
		for k1, v2 := range outboundUp {
			fmt.Println("up")
			//k1越小离得越远
			x, y = g.pos(index1)
			xS = x - 1
			yS = y + g.height[v1] - 1 - k1
			fmt.Println(yS, y, g.height[v2], k1)
			index2 := g.index[v2]
			x, y = g.pos(index2)
			inbound := g.Graph.InBound(v2)
			//入度的坐标也要切
			_, down := g.divide(v2, inbound)
			k2 := util.IndexOf(v1, down)
			//越小就离右边离得越近
			xE = x + g.width[v2] - 1 - k2
			fmt.Println(xE, x, g.width[v2], k2)
			yE = y + g.height[v2]
			fmt.Println(yE, y, g.height[v2])
			yR = yS
			xR = xE
			l := xS - xR
			m.paintA([][]byte{bytes.Repeat([]byte("-"), l)}, xR+1, yR)
			h := yR - yE + 1
			//画"<---+"线,并且转置
			m.paintT([][]byte{append([]byte("^"), append(bytes.Repeat([]byte("|"), h-2), byte('+'))...)}, xE, yE)
		}
	}
}
