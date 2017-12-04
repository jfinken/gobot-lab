package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jfinken/gobot-lab/adabot"
	net "github.com/jfinken/gobot-lab/adabot/network"
)

var bot *adabot.Robot

func defaultHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Gobot says: Takes team work to make the dream work.")
}

// HealthHandler is the HTTP handler that is expected to be used by a load
// balancer.  As such it simply returns HTTP-200
func HealthHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Healthy")
}

// TreadHandler is the handler that is expected to receive a direction and duration in seconds.
// examples:
//  curl host:8181/api/v1/tread/dir/stop
//  curl host:8181/api/v1/tread/dir/forward
//  curl host:8181/api/v1/tread/dir/left
//  curl host:8181/api/v1/tread/dir/right
//  curl host:8181/api/v1/tread/dir/backward
func TreadHandler(ctx *gin.Context) {
	dir := ctx.Param("dir")
	//duration := ctx.Param("dur")

	/*
		sec, err := strconv.Atoi(duration)
		if err != nil {
			errMsg := fmt.Sprintf("PARAM: %s", err.Error())
			ctx.String(http.StatusBadRequest, errMsg)
			return
		}
	*/
	// TODO: validate the directions
	noop := -1
	switch dir {
	case "stop":
		bot.Stop()
	case "forward":
		bot.Forward(noop)
	case "backward":
		bot.Backward(noop)
	case "left":
		bot.Left(noop)
	case "right":
		bot.Right(noop)
	}
	//ctx.String(http.StatusOK, fmt.Sprintf("dir: %s, duration: %s\n", dir, duration))
	ctx.String(http.StatusOK, fmt.Sprintf("dir: %s\n", dir))
}

// ServoHandler handles requests to control the two servo motors charged with yaw/pitch
// direction of the phone/camera pod.
//  curl host:8181/api/v1/pod/dir/yaw/func/-1
//  curl host:8181/api/v1/pod/dir/pitch/func/1
func ServoHandler(ctx *gin.Context) {
	dir := ctx.Param("dir")
	// fn is expected to be a signed int
	f := ctx.Param("func")

	fn, err := strconv.Atoi(f)
	if err != nil {
		errMsg := fmt.Sprintf("PARAM: %s", err.Error())
		ctx.String(http.StatusBadRequest, errMsg)
		return
	}
	switch dir {
	case "yaw":
		bot.Yaw(fn)
	case "pitch":
		bot.Pitch(fn)
	}
	ctx.String(http.StatusOK, fmt.Sprintf("dir: %s, func: %d\n", dir, fn))
}

// RenderNetworkHandler handles requests to display the road network, given the netid, in SVG.
func RenderNetworkHandler(ctx *gin.Context) {
	networkID := ctx.Param("netid")
	var graph *net.RawGraph

	// TODO:
	//	- Decide on the expected data model: is a unique ID stored at the graph level, node level?
	// 	- Trouble getting data into bolt due to circular references of the Node and Edge structs.
	//	- So need multiple structs: a Node struct for encoding into Bolt and a Node struct for A*

	err := graph.Load(networkID)

	if err != nil {
		log.Printf("Network Store err: %s\n", err.Error())
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Network Store err.\n"))
		return
	}
	ctx.Writer.Header().Set("Content-Type", "image/svg+xml")

	// FIXME
	//net.RenderGraph(ctx.Writer, graph)
}

// StoreNetworkHandler stores the bound json payload to local storage.
// curl -H "Content-Type: application/json" --data @body.json http://localhost:8181/api/v1/network/:netid"
func StoreNetworkHandler(ctx *gin.Context) {
	var graph *net.RawGraph
	netID := ctx.Param("netid")
	// This will infer what binder to use depending on the content-type header.
	if ctx.Bind(&graph) == nil {
		err := graph.Store(netID)
		if err != nil {
			log.Printf("Network Store err: %s\n", err.Error())
		}

	} else {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("malformed data\n"))
	}
}

// StorePlanHandler stores the bound json floorplan
// curl -H "Content-Type: application/json" --data @body.json http://localhost:8181/api/v1/floorplan/:planid"
func StorePlanHandler(ctx *gin.Context) {
	var polys []net.Polygon
	planID := ctx.Param("planid")
	// This will infer what binder to use, and unmarshal, depending on the content-type header
	if ctx.Bind(&polys) == nil {
		plan := &net.Floorplan{Polygons: polys}
		err := plan.Store(planID)
		if err != nil {
			log.Printf("Plan Store err: %s\n", err.Error())
		}
	} else {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("Plan Store err: malformed data\n"))
	}
}

// RenderPlanHandler renders the bound json floorplan
func RenderPlanHandler(ctx *gin.Context) {
	planID := ctx.Param("planid")
	var plan *net.Floorplan

	plan, err := plan.Load(planID)

	if err != nil {
		log.Printf("Render plan err: %s\n", err.Error())
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Render plan err.\n"))
		return
	}
	ctx.Writer.Header().Set("Content-Type", "image/svg+xml")

	// Render the plan directly to the http-writer
	plan.Render(ctx.Writer)
}

var wsupgrader = websocket.Upgrader{}

// wsHandler upgrades the gin connection, reads the incoming message
// (expected to be a planID) renders the floorplan as SVG to a string
// buffer then writes that string to the socket.
func wsHandler(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request
	w.Header().Set("Content-Type", "image/svg+xml")

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// TODO: validate the input msg
		var plan *net.Floorplan
		plan, err = plan.Load(string(msg))
		if err != nil {
			fmt.Printf("Failed to render plan: %s\n", string(msg))
			break
		}
		// Render the floorplan to a local string buffer
		buf := new(bytes.Buffer)
		plan.Render(buf)
		conn.WriteMessage(t, []byte(buf.String()))
	}
}
func main() {

	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/health", HealthHandler)
	//router.GET("/api/v1/tread/dir/:dir/duration/:dur", TreadHandler)
	router.GET("/api/v1/tread/dir/:dir", TreadHandler)
	router.GET("/api/v1/pod/dir/:dir/func/:func", ServoHandler)
	router.GET("/api/v1/network/:netid", RenderNetworkHandler)
	router.POST("/api/v1/network/:netid", StoreNetworkHandler)
	router.GET("/api/v1/floorplan/:planid", RenderPlanHandler)
	router.POST("/api/v1/floorplan/:planid", StorePlanHandler)
	router.GET("/ws", wsHandler)
	router.LoadHTMLGlob("./html/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/fp", func(c *gin.Context) {
		c.HTML(http.StatusOK, "fp.html", nil)
	})
	router.Static("/html", "./html") // to serve local js and css files

	robot, err := adabot.NewRobot()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	bot = robot

	port := ":8181"
	fmt.Printf("Listening on %s...\n", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to listen on port(%s): %s", port, err.Error()))
	}
}
