package network

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"crypto/md5"

	"github.com/ajstarks/svgo"
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
type RawGraph struct {
	Nodes []*RawNode `json:"nodes"`
	Edges []*RawEdge `json:"edges"`
	NetID string
}
type NetworkGraph struct {
	Graph map[string]*Node
}

// RawEdge is london dataset-specific but is intended to be marshaled
// to and from key/value storage (no circular references).
type RawEdge struct {
	ID   string `json:"key"`
	St   string `json:"end1key"`
	End  string `json:"end2key"`
	Kind string `json:"kind"`
}

// RawNode is london dataset-specific but is intended to be marshaled
// to and from key/value storage (no circular references).
type RawNode struct {
	ID   string  `json:"key"`
	Name string  `json:"name"`
	Lat  float64 `json:"latitude"`
	Lng  float64 `json:"longitude"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

// Edge type connects two Nodes with a cost.
type Edge struct {
	Raw  *RawEdge
	From *Node
	To   *Node
	Cost float64
}

// A Node is a place in a grid which implements Pather.  ID is a unique
// identifier auto-generated on construction.  NetID allows nodes to
// simply be "grouped" together for the purposes of persistent storage.
// X and Y are the coordinates of the node.  OutTo is a slice of type
// Edge going to other nodes.
type Node struct {
	Raw   *RawNode
	ID    string
	X     int
	Y     int
	OutTo []Edge
	Label string
}

// Space, Wall and Furn here mirror the static ints in TangoPolygon.class
const (
	Space     = iota // 0
	Wall             // 1
	Furniture        // 2
)

// A Polygon holds floor plan polygons
type Polygon struct {
	Area     float64     `json:"area"`
	Layer    int         `json:"layer"`
	IsClosed bool        `json:"isClosed"`
	Verts    [][]float64 `json:"vertices2d"`
}

// A Floorplan defines the polygons that make up a 2D floor plan representation.
type Floorplan struct {
	Polygons []Polygon
}

// AddNode constructs a new Node.
func AddNode(x, y int, rawNode *RawNode, label string) *Node {

	t1 := &Node{X: x, Y: y, Raw: rawNode, Label: label}

	now := time.Now().UnixNano()
	data := []byte(strconv.FormatInt(now, 10))
	t1.ID = fmt.Sprintf("%x", md5.Sum(data))[0:7]

	return t1
}

// AddEdge constructs a new Edge from t1 to t2.
func AddEdge(t1, t2 *Node, cost float64) *Edge {
	edge1 := new(Edge)
	edge1.Cost = cost
	edge1.From = t1
	edge1.To = t2

	t1.OutTo = append(t1.OutTo, *edge1)

	return edge1
}

// PathNeighbors returns the neighbors of the Node.
func (t *Node) PathNeighbors() []astar.Pather {

	neighbors := []astar.Pather{}

	for _, edgeElement := range t.OutTo {
		neighbors = append(neighbors, astar.Pather(edgeElement.To))
	}
	return neighbors
}

// PathNeighborCost returns the cost of the edge leading to Node.
func (t *Node) PathNeighborCost(to astar.Pather) float64 {

	for _, edgeElement := range (t).OutTo {
		if astar.Pather((edgeElement.To)) == to {
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

// GeneratePath invokes the A* path generation function and returns a
// slice of Pather (Nodes), distance and whether or not a successful
// path was found.
func GeneratePath(from, to astar.Pather) (path []astar.Pather, distance float64, found bool) {
	return astar.Path(from, to)
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

// RenderPath will render the nodes and edges to SVG
/*
func RenderPath(w io.Writer, graph *RawGraph) {

	processed := processRawGraph(graph)

	width := 2048
	height := 2048
	vbMinX := 1024 * 1024 * 1024 //math.MaxInt64
	vbMinY := 1024 * 1024 * 1024 //math.MaxInt64
	vbMaxX := 0
	vbMaxY := 0
	// set viewbox min_x,y and max_x, y
	for _, n := range graph.Graph {
		if n.X > 0 {
			vbMinX = min(vbMinX, n.X)
			vbMinY = min(vbMinY, n.Y)

			vbMaxX = max(vbMaxX, n.X)
			vbMaxY = max(vbMaxY, n.Y)
		}
	}
	//fmt.Printf("ViewBox: %d, %d, %d, %d\n", vbMinX, vbMinY, vbMaxX, vbMaxY)

	canvas := svg.New(w)
	// Given the wide-ranging data extents, specify a viewbox: minX, minY, vbWidth, vbHeight
	viewBox := fmt.Sprintf(`viewBox="%d %d %d %d"`,
		vbMinX, vbMinY, (vbMaxX - vbMinX), (vbMaxY - vbMinY))
	aspect := `preserveAspectRatio="xMidYMid meet"`
	canvas.Start(width, height, viewBox, aspect)
	canvas.Rect(vbMinX, vbMinY, (vbMaxX - vbMinX), (vbMaxY - vbMinY), "fill:dimgray")

	// draw fence
	fenX := []int{vbMinX, vbMaxX, vbMaxX, vbMinX, vbMinX}
	fenY := []int{vbMinY, vbMinY, vbMaxY, vbMaxY, vbMinY}
	canvas.Polyline(fenX, fenY, `fill="none"`, `stroke="white"`, `stroke-width:2`)

	// EDGES
	//	for _, edge := range graph.Edges {
	//		if edge.Kind == "Connection" {
	//			from := graph.Graph[edge.St]
	//			to := graph.Graph[edge.End]
	//			//fmt.Printf("[%s, %s]\n", from, to)
	//			canvas.Line(from.X, from.Y, to.X, to.Y, `stroke="skyblue"`, `stroke-width:1`)
	//		}
	//	}

	// Processed NODES and EDGES
	nodeDim := 3
	for _, node := range graph.Graph {
		fill := "fill:white"
		if node.Label == "START" {
			fill = "fill:green"
		} else if node.Label == "END" {
			fill = "fill:red"
		}
		// Draw the node
		//canvas.Rect(n.X, n.Y, nodeDim, nodeDim, fill)
		canvas.Circle(node.X, node.Y, nodeDim, fill)
		// Node label
		//canvas.Text(n.X, n.Y, n.Name, `font-size="8px"`, `fill="red"`)

		// Draw the edges
		for _, e := range node.OutTo {
			canvas.Line(e.From.X, e.From.Y, e.To.X, e.To.Y, `stroke="skyblue"`, `stroke-width:1`)
		}
	}
	canvas.End()
}
*/
func (data *Floorplan) Render(w io.Writer) {
	canvas := svg.New(w)
	canvas.Start(500, 500)
	canvas.Circle(250, 250, 125, "fill:none;stroke:black")
	canvas.End()
}
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
