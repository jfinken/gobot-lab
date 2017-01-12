package main

import (
	"log"
	"math/rand"
	"os"

	"github.com/ajstarks/svgo"
	net "github.com/jfinken/gobot-lab/adabot/network"
)

func rn(n int) int { return rand.Intn(n) }

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
func main() {
	//world := new(GobotWorld)
	diagonalCost := 5.0

	//-------------------------------------------------------------------------
	// Declare A* nodes and edges
	// TODO:
	//	- new HTTP route: POST a network that is written to key/value storage
	//	- new HTTP route: that given a route-ID, reads storage, serves up SVG
	//-------------------------------------------------------------------------
	scale := 10
	networkID := "jf_home"
	nStart := net.AddNode(1*scale, 1*scale, networkID, "START")
	n1 := net.AddNode(1*scale, 2*scale, networkID, "n1")
	n2 := net.AddNode(0*scale, 2*scale, networkID, "n2")
	n3 := net.AddNode(2*scale, 1*scale, networkID, "n3")
	n4 := net.AddNode(3*scale, 1*scale, networkID, "n4")
	n5 := net.AddNode(2*scale, 2*scale, networkID, "n5")
	n6 := net.AddNode(3*scale, 2*scale, networkID, "n6")
	n7 := net.AddNode(2*scale, 3*scale, networkID, "n7")
	n8 := net.AddNode(3*scale, 3*scale, networkID, "n8")
	n9 := net.AddNode(3*scale, 4*scale, networkID, "n9")
	nEnd := net.AddNode(3*scale, 5*scale, networkID, "END")

	nodes := []*net.Node{nStart, n1, n2, n3, n4, n5, n6, n7, n8, n9, nEnd}

	//-------------------------------------------------------------------------
	// POC: store and retrieve nodes
	//-------------------------------------------------------------------------
	store, err := net.OpenStore()
	if err != nil {
		log.Printf("Network Store err: %s\n", err.Error())
	}
	err = store.Update(nodes)
	if err != nil {
		log.Printf("Network Store err: %s\n", err.Error())
	}
	var storedNodes []*net.Node
	err = store.Query(storedNodes, networkID)
	if err != nil {
		log.Printf("Network Store err: %s\n", err.Error())
	}
	err = store.Close()
	if err != nil {
		log.Printf("Network Store err: %s\n", err.Error())
	}

	// effectively set width and height relative to MAX_X, MAX_Y of Node locs
	width := 0
	height := 0
	for _, n := range storedNodes {
		if width <= n.X {
			width = n.X
		}
		if height <= n.Y {
			height = n.Y
		}
	}
	// scale out the dimensions
	width = width * 2
	height = height * 2

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

	// Generate the path.  p is the slice of nodes
	//p, dist, found := net.GeneratePath(nStart, nEnd)

	//-------------------------------------------------------------------------
	// Render nodes and edges SVG
	//-------------------------------------------------------------------------
	//fmt.Printf("WIDTH: %d, HEIGHT: %d\n", width, height)
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:black")

	// draw fence
	fenX := []int{0, width, width, 0, 0}
	fenY := []int{0, 0, height, height, 0}
	canvas.Polyline(fenX, fenY, `fill="none"`, `stroke="red"`, `stroke-width:3`)

	// render nodes
	nodeDim := 3
	for _, n := range nodes {
		fill := "fill:white"
		if n.Label == "START" {
			fill = "fill:green"
		} else if n.Label == "END" {
			fill = "fill:red"
		}
		canvas.Rect(n.X, n.Y, nodeDim, nodeDim, fill)
	}
	canvas.End()
}
