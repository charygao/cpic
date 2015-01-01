package cpic // import "github.com/ggaaooppeenngg/cpic"

//inspired by gyuho's goraph [https://github.com/gyuho/goraph/blob/refactor/graph/graph.go]
import (
	"fmt"
	"log"
	"sync"
)

// 当函数要修改数据的时候才用指针作为接收者.
// graph contains graph datas,implemented by adjacency list.
type graph struct {
	vertices []*vertex

	sync.Mutex

	outEdge   map[*vertex][]*edge // outBound maps vertex to outgoing edges
	inEdge    map[*vertex][]*edge // inBound maps vertex to outgoing edges
	vertexIDs map[string]bool     // prevent duplicating IDs
}

// newGraph inits a new graph.
func newGraph() *graph {
	return &graph{
		vertices:  []*vertex{},
		outEdge:   make(map[*vertex][]*edge),
		inEdge:    make(map[*vertex][]*edge),
		vertexIDs: make(map[string]bool),
	}
}

//vertex is a vertex(node) in graph
type vertex struct {
	id  string //for debug
	ele interface{}
}

//newVertex inits a new vertext.
func newVertext(id string) *vertex {
	return &vertex{
		id: id,
	}
}

type edge struct {
	//source or destination
	vtx    *vertex
	weight float64
}

//for debug
func (g graph) String() string {
	var output string
	for _, vtx := range g.vertices {
		output += fmt.Sprintf("[%s]:\n", vtx.id)
		for _, edge := range g.outEdge[vtx] {
			output += fmt.Sprintf("-> %f %s\n", edge.weight, edge.vtx.id)
		}
	}
	return output
}

func (g *graph) addVertex(vtx *vertex) error {
	if _, ok := g.vertexIDs[vtx.id]; ok {
		return fmt.Errorf("'%s' already exists", vtx.id)
	}
	g.Mutex.Lock()
	//map 不是并发安全的
	g.vertexIDs[vtx.id] = true
	g.Mutex.Unlock()
	g.vertices = append(g.vertices, vtx)
	return nil
}

func (g graph) findVertexById(id string) *vertex {
	for _, vtx := range g.vertices {
		if vtx.id == id {
			return vtx
		}
	}
	return nil
}

func (g *graph) connect(src, dst *vertex, weight float64) {
	err := g.addVertex(src)
	if err != nil {
		log.Printf("'%s' was previously added to graph\n", src.id)
		src = g.findVertexById(src.id)
	} else {
		log.Printf("'%s' is added to graph\n", src.id)
	}
	err = g.addVertex(dst)
	if err != nil {
		log.Printf("'%s' was previously added to graph\n", dst.id)
		dst = g.findVertexById(dst.id)
	} else {
		log.Printf("'%s' is added to graph\n", dst.id)
	}
	edgeSrc := &edge{
		vtx:    src,
		weight: weight,
	}
	edgeDst := &edge{
		vtx:    dst,
		weight: weight,
	}
	g.Mutex.Lock()
	if _, ok := g.outEdge[src]; !ok {
		g.outEdge[src] = []*edge{edgeDst}
	} else {
		isDuplicate := false
		for _, edge := range g.outEdge[src] {
			if edge.vtx == dst {
				log.Println("Duplicate(Parallel) Edge Found. Overwriting the Weight value.")
				log.Printf("%v --> %v + %v\n", edge.weight, edge.weight, weight)
				edge.weight += weight
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			g.outEdge[src] = append(g.outEdge[src], edgeDst)
		}
	}
	if _, ok := g.inEdge[dst]; !ok {
		g.inEdge[dst] = []*edge{edgeSrc}
	} else {
		isDuplicate := false
		for _, edge := range g.inEdge[dst] {
			if edge.vtx == src {
				log.Println("Duplicate(Parallel) Edge Found. Overwriting the Weight value.")
				log.Printf("%v --> %v + %v\n", edge.weight, edge.weight, weight)
				edge.weight += weight
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			g.inEdge[dst] = append(g.inEdge[dst], edgeSrc)
		}

	}
	g.Mutex.Unlock()

}
