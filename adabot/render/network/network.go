package network

import (
	"github.com/jfinken/go-astar"
)

// astar_client.go implements implements astar.Pather for
// the sake of testing.  This functionality forms the back end for
// this command's main.go, and serves as an example for how to use A* for a graph.

// Nodes are called 'Nodes' and they have X, Y coordinates
// Edges are called 'Edges', they connect Nodes, and they have a cost
//
// NOTES:
// 1) There is no grid.  Nodes have arbitrary coordinates.
// 2) Edges are not implied by the grid positions.  Instead edges are explicitly
//    modelled as 'Edges'.
// 3) Manhattan distance is used as the heuristic
// 4) The astar.Pather interface is implemented

// GobotWorld will eventually hold a map of type Node
type GobotWorld struct {
	//	nodes map[int]*Node		// not yet used
}

// Edge type connects two Nodes with a cost
type Edge struct {
	from *Node
	to   *Node
	Cost float64
}

// A Node is a place in a grid which implements Grapher.
type Node struct {

	// X and Y are the coordinates of the node.
	X, Y int

	// array of type Edge going to other nodes
	outTo []Edge

	Label string
}

func AddNode(x int, y int, label string) *Node {
	t1 := new(Node)
	t1.X = x
	t1.Y = y
	t1.Label = label
	return t1
}

func AddEdge(t1, t2 *Node, cost float64) *Edge {
	edge1 := new(Edge)
	edge1.Cost = cost
	edge1.from = t1
	edge1.to = t2

	t1.outTo = append(t1.outTo, *edge1)

	return edge1
}

// PathNeighbors returns the neighbors of the Truck
func (t *Node) PathNeighbors() []astar.Pather {

	neighbors := []astar.Pather{}

	for _, edgeElement := range t.outTo {
		neighbors = append(neighbors, astar.Pather(edgeElement.to))
	}
	return neighbors
}

// PathNeighborCost returns the cost of the edge leading to Node.
func (t *Node) PathNeighborCost(to astar.Pather) float64 {

	for _, edgeElement := range (t).outTo {
		if astar.Pather((edgeElement.to)) == to {
			return edgeElement.Cost
		}
	}
	return 10000000
}

// PathEstimatedCost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *Node) PathEstimatedCost(to astar.Pather) float64 {

	toT := to.(*Node)
	absX := toT.X - t.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	r := float64(absX + absY)

	return r
}

// RenderPath renders a path on top of a Goreland world.
func (w GobotWorld) RenderPath(path []astar.Pather) string {

	s := ""
	for _, p := range path {
		pT := p.(*Node)
		if pT.Label != "END" {
			s = pT.Label + "->" + s
		} else {
			s = s + pT.Label
		}
	}
	return s
}
