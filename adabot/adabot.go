package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/chzyer/readline"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
func adafruitDCMotorRunner(a *i2c.AdafruitMotorHatDriver, dcMotor int) (err error) {

	//log.Printf("%s\tRun Loop...\n", time.Now().String())
	// set the speed:
	var speed int32 = 255 // 255 = full speed!
	if err = a.SetDCMotorSpeed(dcMotor, speed); err != nil {
		return
	}
	// run FORWARD
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitForward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(2000 * time.Millisecond)
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitRelease); err != nil {
		return
	}
	// run BACKWARD
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitBackward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(2000 * time.Millisecond)
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitRelease); err != nil {
		return
	}
	return
}
func main() {
	gbot := gobot.NewGobot()
	r := raspi.NewRaspiAdaptor("raspi")
	adaFruit := i2c.NewAdafruitMotorHatDriver(r, "adafruit")

	robot := gobot.NewRobot("adaFruitBot",
		[]gobot.Connection{r},
		[]gobot.Device{adaFruit},
		nil, //work,
	)

	gbot.AddRobot(robot)

	/*
	* 2016-08-24 jfinken: commenting out the work func and the Start loop
	* and using readline CLI.
	*
	* I reject your reality and substitue my own...
	 */
	robot.Start() // required for init
	dcMotor := 3  // 0-based

	pi := "\xCE\xA0"
	fmt.Printf("Come to the dork side we have %s\n", pi)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "gobot> ",
		HistoryFile: getHomeDir() + "/.gobot_history",
	})

	if err != nil {
		panic(err)
	}
	defer rl.Close()
	// start cli loop...
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		// detect EOF before parsing
		if line == "" {
			continue
		}
		// TODO: switch on WASD
		adafruitDCMotorRunner(adaFruit, dcMotor)
	}
}
