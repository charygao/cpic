package cpic

import (
	"testing"

	"github.com/ggaaooppeenngg/util/container"
)

func TestGraph(t *testing.T) {
	var g *graph
	var width, height int
	g = newGraph()
	a := &container.Vertex{Id: "AC"}
	b := &container.Vertex{Id: "B"}
	c := &container.Vertex{Id: "C"}
	g.Graph.Connect(a, b)
	g.Graph.Connect(b, c)
	width, height = g.scale()
	if width != 6 || height != 7 {
		t.Log(height, width)
		t.Fatal("height != 7,width !=6")
	}
	g = newGraph()
	g.Graph.Connect(a, b, 123)
	width, height = g.scale()
	if width != 11 || height != 4 {
		t.Log(height, width)
		t.Fatal("height !=4,width !=11")
	}
}
