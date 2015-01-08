package cpic

import (
	"testing"

	"github.com/ggaaooppeenngg/util/container"
)

func TestGraph(t *testing.T) {
	g := newGraph()
	a := &container.Vertex{Id: "A"}
	b := &container.Vertex{Id: "B"}
	c := &container.Vertex{Id: "C"}
	g.Graph.Connect(a, b, 0)
	g.Graph.Connect(b, c, 0)
	t.Log(g.scale())
}
