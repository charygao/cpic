package cpic // import "github.com/ggaaooppeenngg/cpic"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]

import (
	"bytes"
	"strconv"

	"github.com/ggaaooppeenngg/util"
	"github.com/ggaaooppeenngg/util/container"
)

var precision int

type info struct {
	x, y, w, h int
}

//graph wraps container.Graph
type graph struct {
	*container.Graph                            //Graph container
	info             map[*container.Vertex]info //postion of vertex in this graph
	index            map[*container.Vertex]int  //index of vertex in Graph.Vertecies
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
		make(map[*container.Vertex]info),
		make(map[*container.Vertex]int)}
}

//TODO:graph要有某种合理的排序使得空间不那么拥挤.
func (g graph) scale() (width int, height int) {
	var offsetX, offsetY int //offsetX is previous node's right most x,offsetY is previous node's bottom most y.
	for k1 := 0; k1 < len(g.Vertices); k1++ {
		var (
			v1           = g.Vertices[k1]
			weightLenMax = 0          // max length of weight stringifed.
			d, r         = 0, 0       // bottom inbound,right outbound.
			t, l         = 0, 0       // top inbound, left outbound
			idLen        = len(v1.Id) // length of Id.
		)

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
		for _, v2 := range g.Graph.Vertices[:k1] {
			if g.Graph.EdgeExists(v2, v1) {
				t++
				weight := g.Graph.GetEdgeWeight(v2, v1)
				if weight != nil {
					if weightLenMax < len(strconv.FormatFloat(*weight, 'f', 3, 64)) {
						weightLenMax = len(strconv.FormatFloat(*weight, 'f', 3, 64))
					}
				}
			}
			if g.Graph.EdgeExists(v1, v2) {
				l++
				weight := g.Graph.GetEdgeWeight(v1, v2)
				if weight != nil {
					if weightLenMax < len(strconv.FormatFloat(*weight, 'f', 3, 64)) {
						weightLenMax = len(strconv.FormatFloat(*weight, 'f', 3, 64))
					}
				}
			}

		}

		//TODO:add lock to map

		info := g.info[v1] //用于存储定位后的结果
		//宽度最大值
		w := getMax(1, t, d, idLen) //at least width one
		info.w = w
		//长度最大值
		h := getMax(1, l, r) //at least height one
		info.h = h
		if k1 != 0 {
			info.x = offsetX + 2 + weightLenMax
		}
		offsetX = info.x + info.w - 1
		if k1 != 0 {
			info.y = offsetY + 3
		}
		offsetY = info.y + info.h - 1
		g.info[v1] = info
	}

	return offsetX + 1, offsetY + 1
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

//quick sort
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

//divide divides vertecies into two parts,left vtxs and right vtxs of vtx.
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
	return
}

//opsition get the node x,y position at the index of Vertices.
func (g graph) position(index int) (x, y int) {
	x = g.info[g.Vertices[index]].x
	y = g.info[g.Vertices[index]].y
	return
}

//draw implement painter interface,paint chars a matrix
func (g graph) draw(m *matrix) {
	g.sort() //TODO:使得分布比较合理

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

	//draw from top to down
	for index1, v1 := range g.Graph.Vertices {
		var xS, yS, xR, yR, xE, yE, x, y int //start,end,temporary position
		x, y = g.position(index1)
		m.paintA([][]byte{[]byte(v1.Id)}, x, y)
		outbound := g.Graph.OutBound(v1)
		//divide into down and up outbound
		//把每个出度分成在这个点之前后之后的,因为绘制箭头的方式不一样.
		outboundUp, outboundDown := g.divide(v1, outbound)

		var index2 int // index of outbound vertext in global g.Graph.Vertices.
		for k1, v2 := range outboundDown {
			//get postion of first '-'
			x, y = g.position(index1)
			//越近就越在下面
			xS = x + g.info[v1].w
			yS = y + g.info[v1].h - 1 - k1 //k1 counts from 0.if k1 = 0,outbound is the most bottom outbound.

			//get position of 'v'
			inboundOfV2 := g.Graph.InBound(v2)
			up, _ := g.divide(v1, inboundOfV2)
			k2 := util.IndexOf(v1, up)
			index2 = g.index[v2]
			x, y = g.position(index2)
			//越近就越在左边
			xE = x + g.info[v2].w - 1 - k2 //k2 counts from 0.if k2 = 0,inbound is the most left inbound.
			yE = y - 1

			//postion of '+'
			yR = yS
			xR = xE

			l := xR - xS // l >= 1
			weight := g.GetEdgeWeight(v1, v2)
			var weightString string
			if weight != nil {
				weightString = strconv.FormatFloat(*weight, 'f', 3, 64)
			}
			m.paintA([][]byte{append([]byte("-"+weightString), bytes.Repeat([]byte("-"), l-(1+len(weightString)))...)}, xS, yS)
			//draw vertical line
			h := yE - yR + 1
			m.paintT([][]byte{append([]byte("+"), append(bytes.Repeat([]byte("|"), h-2), byte('v'))...)}, xR, yR)
		}
		for k1, v2 := range outboundUp {
			//get position of first '-'
			x, y = g.position(index1)
			xS = x - 1
			yS = y + g.info[v1].h - 1 - k1 //k1 counts from 0,if k1 = 0,outbound is the most bottom outboound.

			//get postion of '^'
			index2 := g.index[v2]
			x, y = g.position(index2)
			inbound := g.Graph.InBound(v2)
			_, down := g.divide(v2, inbound)
			k2 := util.IndexOf(v1, down)
			xE = x + g.info[v2].w - 1 - k2 //k2 counts from 0,if k2 = 0,inbound is the most right inbound.
			yE = y + g.info[v2].h

			//get postion of '+'
			yR = yS
			xR = xE

			//draw horizontal line.
			l := xS - xR // l >= 1.
			weight := g.GetEdgeWeight(v1, v2)
			var weightString string
			if weight != nil {
				weightString = strconv.FormatFloat(*weight, 'f', 3, 64)
			}
			m.paintA([][]byte{append(bytes.Repeat([]byte("-"), l-(1+len(weightString))), []byte(weightString+"-")...)}, xR+1, yR)
			//draw vertical line.
			h := yR - yE + 1
			m.paintT([][]byte{append([]byte("^"), append(bytes.Repeat([]byte("|"), h-2), byte('+'))...)}, xE, yE)
		}
	}
}
