package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jfinken/gobot-lab/adabot"
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
//  curl host:8181/api/v1/tread/dir/forward/duration/5
//  curl host:8181/api/v1/tread/dir/left/duration/2
//  curl host:8181/api/v1/tread/dir/right/duration/2
//  curl host:8181/api/v1/tread/dir/backward/duration/5
func TreadHandler(ctx *gin.Context) {
	dir := ctx.Param("dir")
	duration := ctx.Param("dur")

	sec, err := strconv.Atoi(duration)
	if err != nil {
		errMsg := fmt.Sprintf("PARAM: %s", err.Error())
		ctx.String(http.StatusBadRequest, errMsg)
		return
	}
	// TODO: validate the directions
	switch dir {
	case "forward":
		bot.Forward(sec)
	case "backward":
		bot.Backward(sec)
	case "left":
		bot.Left(sec)
	case "right":
		bot.Right(sec)
	}
	ctx.String(http.StatusOK, fmt.Sprintf("dir: %s, duration: %s\n", dir, duration))
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
func main() {

	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/health", HealthHandler)
	router.GET("/api/v1/tread/dir/:dir/duration/:dur", TreadHandler)
	router.GET("/api/v1/pod/dir/:dir/func/:func", ServoHandler)
	router.LoadHTMLGlob("./html/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

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
