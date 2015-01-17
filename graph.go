package cpic // import "github.com/ggaaooppeenngg/cpic"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/ggaaooppeenngg/util"
	"github.com/ggaaooppeenngg/util/container"
)

//graph wraps container.Graph
type graph struct {
	*container.Graph                           //Graph container
	width            map[*container.Vertex]int //width of vertex in this graph
	height           map[*container.Vertex]int //height of vertex in this graph
	offsetX          map[*container.Vertex]int //x pixl offset of a vertex
	index            map[*container.Vertex]int //index of vertex in Graph.Vertecies
}

//getMax returns max number of a and b,if a and b are both smaller thanl,returns l.
func getMax(min int, nums ...int) int {
	var max = -2147483648 //minnest number of int32
	for i := 0; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	if min > max {
		max = min
	}
	return max
}

//newGraph returns new graph and initiates the maps.
func newGraph() *graph {
	return &graph{container.NewGraph(),
		make(map[*container.Vertex]int),
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
	//有了代价的时候,终点的位置会相应改变,只要获得x值最大的那个,就是对应的终点的坐标.
	//要考虑up和down方向.
	var weightLenMax int = 0 // max length of weight stringifed.
	for k1 := 0; k1 < len(g.Vertices); k1++ {
		fmt.Println("index", k1)
		var v1 = g.Vertices[k1]
		var d, r int               // bottom inbound,right outbound.
		var idLen int = len(v1.Id) // length of Id.
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
		var t, l int // top inbound, left outbound
		for _, v2 := range g.Graph.Vertices[:k1] {
			if g.Graph.EdgeExists(v2, v1) {
				t++
				weight := g.Graph.GetEdgeWeight(v2, v1)
				if weight != nil {
					if weightLenMax < len(strconv.FormatFloat(*weight, 'f', 3, 64)) {
						weightLenMax = len(strconv.FormatFloat(*weight, 'f', 3, 64))
					}
					//fmt.Println(strconv.FormatFloat(*weight, 'f', 3, 64))
				}
			}
			if g.Graph.EdgeExists(v1, v2) {
				l++
				weight := g.Graph.GetEdgeWeight(v1, v2)
				if weight != nil {
					if weightLenMax < len(strconv.FormatFloat(*weight, 'f', 3, 64)) {
						weightLenMax = len(strconv.FormatFloat(*weight, 'f', 3, 64))
					}
					//fmt.Println(strconv.FormatFloat(*weight, 'f', 3, 64))
				}
			}

		}
		//fmt.Println("weight length max", weightLenMax)
		g.offsetX[v1] = weightLenMax
		//accumulate width and height
		w := getMax(1, t, d, idLen) //at least width one
		g.width[v1] = w
		h := getMax(1, l, r) //at least height one
		g.height[v1] = h
		fmt.Println("weight len max", weightLenMax)
		if k1 != len(g.Vertices)-1 {
			width += w + 1
			height += h + 2
		} else {
			width += w
			height += h
		}
	}
	//danshi zhengti de guimo hai buzhi
	//fmt.Println("scale", width, height)
	return width + weightLenMax, height
}

func (g *graph) sort() {
	for i, k := range g.Graph.Vertices {
		g.index[k] = i
	}
}
func isSorted(vtxg, vtxs []*container.Vertex) bool {
	/*
		for i := 0; i < len(vtxs)-1; i++ {
			if util.IndexOf(vtxs[i], vtxg) > util.IndexOf(vtxs[i+1], vtxg) {
				return false
			}
		}
	*/
	return true
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
	if !isSorted(g.Graph.Vertices, vtxs) {
		sort(g.Graph.Vertices, vtxs, 0, len(vtxs))
	}
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
	//fmt.Println("divider", divider)
	return
}

//获得左上角的位置.
func (g graph) pos(index int) (x, y int) {
	for _, v := range g.Graph.Vertices[:index] {
		x += g.width[v] + 1  //空格
		y += g.height[v] + 2 //下箭头
	}
	x += g.offsetX[g.Vertices[index]]
	return
}

func (g graph) draw(m *matrix) {
	g.sort()
	width, height := m.scale()
	fmt.Println("draw weigth width height", width, height)
	//一趟向下
	for index1, v1 := range g.Graph.Vertices {
		//fmt.Println("loop", v1.Id)
		var xS, yS, xR, yR, xE, yE, x, y int //start,end,temporary position
		x, y = g.pos(index1)
		m.paintA([][]byte{[]byte(v1.Id)}, x, y)
		outbound := g.Graph.OutBound(v1)
		//divide to down and up outbound
		//get divider
		//对于每个点的出度,都划分成上下两部分,
		//上部分是在全局数组g.Graph.Vertices中在v1,前面的点
		//下部分是在全局数组里面
		outboundUp, outboundDown := g.divide(v1, outbound)

		var index2 int // index of outbound vertext
		for k1, v2 := range outboundDown {
			//fmt.Println("down")
			x, y = g.pos(index1)
			//右边的出度,最接近,最下面
			//第k1个出度
			//起点.
			xS = x + g.width[v1]
			yS = y + g.height[v1] - 1 - k1 //k1 counts from 0
			//fmt.Println(yS, y, g.height[v1], k1)
			//对于v2来说,v1在上方.
			inbound := g.Graph.InBound(v2)
			up, _ := g.divide(v1, inbound)
			//对于v2,v1的位置.,最后得出v1到v2的路线.
			k2 := util.IndexOf(v1, up)
			index2 = g.index[v2]
			x, y = g.pos(index2)
			xE = x + g.width[v2] - 1 - k2 //k2 counts from 0
			//fmt.Println(xE, x, g.width[v2], k2)
			yE = y - 1
			yR = yS
			xR = xE
			l := xR - xS // l >= 1
			//画"---"线
			//fmt.Println("paint")
			//fmt.Println("paint -", l, "times", "start and end", xS, yS)
			weight := g.GetEdgeWeight(v1, v2)
			var weightString string
			if weight != nil {
				weightString = strconv.FormatFloat(*weight, 'f', 3, 64)
			}
			m.paintA([][]byte{append([]byte("-"+weightString), bytes.Repeat([]byte("-"), l-(1+len(weightString)))...)}, xS, yS)
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
			fmt.Println("yS", yS, y, g.height[v2], k1)
			index2 := g.index[v2]
			x, y = g.pos(index2)
			inbound := g.Graph.InBound(v2)
			//入度的坐标也要切
			_, down := g.divide(v2, inbound)
			k2 := util.IndexOf(v1, down)
			//越小就离右边离得越近
			xE = x + g.width[v2] - 1 - k2
			//fmt.Println(xE, x, g.width[v2], k2)
			yE = y + g.height[v2]
			//fmt.Println(yE, y, g.height[v2])
			yR = yS
			xR = xE
			l := xS - xR //l >=1
			weight := g.GetEdgeWeight(v1, v2)
			var weightString string
			if weight != nil {
				weightString = strconv.FormatFloat(*weight, 'f', 3, 64)
			}
			m.paintA([][]byte{append(bytes.Repeat([]byte("-"), l-(1+len(weightString))), []byte(weightString+"-")...)}, xR+1, yR)
			h := yR - yE + 1
			//画"<---+"线,并且转置
			m.paintT([][]byte{append([]byte("^"), append(bytes.Repeat([]byte("|"), h-2), byte('+'))...)}, xE, yE)
		}
		fmt.Println(m.output())
	}
}
