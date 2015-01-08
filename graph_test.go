package cpic

import (
	"testing"

	"github.com/ggaaooppeenngg/util/container"
)

func TestGraph(t *testing.T) {
	g := newGraph()
	g.Graph.Connect(&container.Vertex{Id: "A"}, &container.Vertex{Id: "B"}, 0)
	t.Log(g.scale())
}
