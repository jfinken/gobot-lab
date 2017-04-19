package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/svgo"
	"github.com/golang/geo/s2"

	net "github.com/jfinken/gobot-lab/adabot/network"
)

// POC: Declare A*-ready nodes and edges
//		   	 E
//			 |
//		   	 N9
//		  /	 |
//    	 /	 |
//       N7	 N8
//       |	 |
//N2--N1 N5--N6
//    |/ 	 |
//    S--N3--N4
//
// S=Start at (1,1)
// E=End at (3,5)
// GOAL: draw nodes relative to width, height
// S=Start at (1,1)
// E=End at (3,5)
//
// N1 = (1, 2)
// N2 = (0, 2)
// N3 = (2, 1)
// N4 = (3, 1)
// N5 = (2, 2)
// N6 = (3, 2)
// N7 = (2, 3)
// N8 = (3, 3)
// N9 = (3, 4)
/*
func pocToyExample() []*net.Node {

	//world := new(GobotWorld)
	diagonalCost := 5.0

	scale := 10
	nStart := net.AddNode(1*scale, 1*scale, "START")
	n1 := net.AddNode(1*scale, 2*scale, "n1")
	n2 := net.AddNode(0*scale, 2*scale, "n2")
	n3 := net.AddNode(2*scale, 1*scale, "n3")
	n4 := net.AddNode(3*scale, 1*scale, "n4")
	n5 := net.AddNode(2*scale, 2*scale, "n5")
	n6 := net.AddNode(3*scale, 2*scale, "n6")
	n7 := net.AddNode(2*scale, 3*scale, "n7")
	n8 := net.AddNode(3*scale, 3*scale, "n8")
	n9 := net.AddNode(3*scale, 4*scale, "n9")
	nEnd := net.AddNode(3*scale, 5*scale, "END")

	// Create EDGES.  Note this modifies nodes.
	net.AddEdge(nStart, n1, 1)
	net.AddEdge(nStart, n5, diagonalCost)
	net.AddEdge(nStart, n3, 1)
	net.AddEdge(n3, n4, 1)
	net.AddEdge(n1, n2, 1)
	net.AddEdge(n5, n6, 1)
	net.AddEdge(n5, n7, 1)
	net.AddEdge(n4, n6, 1)
	net.AddEdge(n6, n8, 1)
	net.AddEdge(n8, n9, 1)
	net.AddEdge(n7, n9, diagonalCost)
	net.AddEdge(n9, nEnd, 1)

	nodes := []*net.Node{nStart, n1, n2, n3, n4, n5, n6, n7, n8, n9, nEnd}
	return nodes
}
*/
// POC: store and retrieve nodes via dead-simple JSON flat files
func pocStoreRetrieveNodes(graph *net.RawGraph, networkID string) {

	graph.NetID = networkID
	err := net.StoreGraph(graph, networkID)
	if err != nil {
		log.Printf("Network Store err: %s\n", err.Error())
	}

	storedGraph := &net.RawGraph{}
	err = net.LoadGraph(storedGraph, networkID)
	if err != nil {
		log.Printf("Network Get Store err: %s\n", err.Error())
	}
	log.Printf("LEN of raw Nodes in results: %d\n", len(storedGraph.Nodes))
}

// POC: Unmarshal and render nodes of the London tube
func pocUnmarshalLondon() (*net.RawGraph, *net.NetworkGraph, *net.Node, *net.Node) {

	file, e := ioutil.ReadFile("./london_tube.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var rawGraph net.RawGraph
	var london net.NetworkGraph
	london.Graph = make(map[string]*net.Node) // type alias

	json.Unmarshal(file, &rawGraph)

	// transform spherical to cartesian
	scale := 100000.0
	for _, rawNode := range rawGraph.Nodes {
		pt := s2.PointFromLatLng(s2.LatLngFromDegrees(rawNode.Lat, rawNode.Lng))
		if pt.X >= 1.0 || pt.Y == 0.0 {
			continue
		}
		// Create a processed Node
		node := net.AddNode(int(pt.X*scale), int(pt.Y*scale), rawNode, rawNode.Name)

		// store the processed node at key
		london.Graph[rawNode.ID] = node
	}
	// for A* POC purposes:
	//	- generate a random start and end node
	//	- store the edge with the connecting nodes and assign a cost
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)
	rawStart := rawGraph.Nodes[r.Intn(len(rawGraph.Nodes)-1)]
	rawEnd := rawGraph.Nodes[r.Intn(len(rawGraph.Nodes)-1)]

	nStart := london.Graph[rawStart.ID]
	nEnd := london.Graph[rawEnd.ID]

	// Fascinating: given these nodes and haversine as the cost.
	//nStart := london.Graph["290"] // West Hampstead
	//nEnd := london.Graph["185"]   // North Wembley

	for _, edge := range rawGraph.Edges {
		if edge.Kind == "Connection" {
			from := london.Graph[edge.St]
			to := london.Graph[edge.End]

			//cost := 1.0
			// Set Cost to be the actual distance between connected nodes
			n1 := s2.LatLngFromDegrees(from.Raw.Lat, from.Raw.Lng)
			n2 := s2.LatLngFromDegrees(to.Raw.Lat, to.Raw.Lng)
			cost := n1.Distance(n2).Degrees()

			// AddEdge will modify the parameter Nodes.  NOTE: bidirectionality..
			net.AddEdge(from, to, cost)
			net.AddEdge(to, from, cost)
		}
	}

	return &rawGraph, &london, nStart, nEnd
}
func rn(n int) int { return rand.Intn(n) }

func main() {
	//-------------------------------------------------------------------------
	// POC: Declare A* nodes and edges
	// TODO:
	//	- new HTTP route: POST a network that is written to key/value storage
	//	- new HTTP route: that given a route-ID, reads storage, serves up SVG
	//-------------------------------------------------------------------------
	//nodes := pocToyExample()

	//-------------------------------------------------------------------------
	// POC: store and retrieve nodes of the London tube
	//-------------------------------------------------------------------------
	rawGraph, london, nStart, nEnd := pocUnmarshalLondon()

	//-------------------------------------------------------------------------
	// POC: store and retrieve nodes
	//-------------------------------------------------------------------------
	pocStoreRetrieveNodes(rawGraph, "london1234")

	// Generate the path.  p is the slice of nodes
	pathNodes, dist, pathFound := net.GeneratePath(nStart, nEnd)
	fmt.Fprintf(os.Stderr, "%s to %s, dist: %f\n", nStart.Raw.Name, nEnd.Raw.Name, dist)

	//-------------------------------------------------------------------------
	// Render nodes and edges SVG
	//-------------------------------------------------------------------------
	width := 2048
	height := 2048
	vbMinX := math.MaxInt64
	vbMinY := math.MaxInt64
	vbMaxX := 0
	vbMaxY := 0
	// set viewbox min_x,y and max_x, y
	for _, n := range london.Graph {
		if n.X > 0 {
			vbMinX = min(vbMinX, n.X)
			vbMinY = min(vbMinY, n.Y)

			vbMaxX = max(vbMaxX, n.X)
			vbMaxY = max(vbMaxY, n.Y)
		}
	}
	//fmt.Printf("ViewBox: %d, %d, %d, %d\n", vbMinX, vbMinY, vbMaxX, vbMaxY)

	canvas := svg.New(os.Stdout)
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

	/*
		// EDGES
		for _, edge := range london.RawEdges {
			if edge.Kind == "Connection" {
				from := london.Graph[edge.St]
				to := london.Graph[edge.End]
				//fmt.Printf("[%s, %s]\n", from, to)
				canvas.Line(from.X, from.Y, to.X, to.Y, `stroke="skyblue"`, `stroke-width:1`)
			}
		}
	*/

	// NODES
	nodeDim := 3
	for _, node := range london.Graph {
		fill := "fill:white"
		if node.Label == "START" {
			fill = "fill:green"
		} else if node.Label == "END" {
			fill = "fill:red"
		}
		//canvas.Rect(n.X, n.Y, nodeDim, nodeDim, fill)
		canvas.Circle(node.X, node.Y, nodeDim, fill)
		// Node label
		//canvas.Text(n.X, n.Y, n.Name, `font-size="8px"`, `fill="red"`)

		// Draw the edges
		for _, e := range node.OutTo {
			canvas.Line(e.From.X, e.From.Y, e.To.X, e.To.Y, `stroke="skyblue"`, `stroke-width:1`)
		}
	}

	// A* results
	if pathFound {
		nodeDim = 4.0
		// it is an ordered path
		for i := range pathNodes {
			var from *net.Node
			var to *net.Node
			//canvas.Circle(n.(*net.Node).X, n.(*net.Node).Y, nodeDim, "fill:green")
			if i < len(pathNodes)-1 {
				from = pathNodes[i].(*net.Node)
				to = pathNodes[i+1].(*net.Node)
			} else {
				from = pathNodes[len(pathNodes)-2].(*net.Node)
				to = pathNodes[len(pathNodes)-1].(*net.Node)
			}
			canvas.Line(from.X, from.Y, to.X, to.Y, `stroke="red"`, `stroke-width:1`)
		}
	}
	// Start and End nodes
	canvas.Circle(nStart.X, nStart.Y, 4.0, "fill:green")
	canvas.Text(nStart.X, nStart.Y, nStart.Label, `font-size="8px"`, `fill="white"`)
	canvas.Circle(nEnd.X, nEnd.Y, 4.0, "fill:red")
	canvas.Text(nEnd.X, nEnd.Y, nEnd.Label, `font-size="8px"`, `fill="white"`)
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
