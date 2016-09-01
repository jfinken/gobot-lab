package adabot

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

// Robot defines a type abstracting the unexported driver type.
type Robot struct {
	adafruit *i2c.AdafruitMotorHatDriver
}

// NewRobot constructs and initializes an unexported driver object.
func NewRobot() *Robot {

	gbot := gobot.NewGobot()
	r := raspi.NewRaspiAdaptor("raspi")
	adaFruit := i2c.NewAdafruitMotorHatDriver(r, "adafruit")

	robot := gobot.NewRobot("adaFruitBot",
		[]gobot.Connection{r},
		[]gobot.Device{adaFruit},
		nil, // nil work func
	)

	gbot.AddRobot(robot)

	// effectively init
	robot.Start()
	return &Robot{adafruit: adaFruit}
}

// Left runs both DC-Motors in opposite directions for the given amount of time in seconds.
func (bot *Robot) Left(sec int) (err error) {
	motorPort := 0
	motorStarboard := 1
	var speed int32 = 255 // 255 = full speed!
	if err = bot.adafruit.SetDCMotorSpeed(motorPort, speed); err != nil {
		return
	}
	if err = bot.adafruit.SetDCMotorSpeed(motorStarboard, speed); err != nil {
		return
	}
	//--------------------------
	// BUG: direction is flipped
	//--------------------------
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitForward); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitBackward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(time.Duration(sec) * time.Second)
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitRelease); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitRelease); err != nil {
		return
	}
	return
}

// Right runs both DC-Motors in opposite directions for the given amount of time in seconds.
func (bot *Robot) Right(sec int) (err error) {
	motorPort := 0
	motorStarboard := 1
	var speed int32 = 255 // 255 = full speed!
	if err = bot.adafruit.SetDCMotorSpeed(motorPort, speed); err != nil {
		return
	}
	if err = bot.adafruit.SetDCMotorSpeed(motorStarboard, speed); err != nil {
		return
	}
	// BUG: direction
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitBackward); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitForward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(time.Duration(sec) * time.Second)
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitRelease); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitRelease); err != nil {
		return
	}
	return
}

// Backward runs both DC-Motors backward for the given amount of time in seconds.
func (bot *Robot) Backward(sec int) (err error) {
	motorPort := 0
	motorStarboard := 1
	var speed int32 = 255 // 255 = full speed!
	if err = bot.adafruit.SetDCMotorSpeed(motorPort, speed); err != nil {
		return
	}
	if err = bot.adafruit.SetDCMotorSpeed(motorStarboard, speed); err != nil {
		return
	}
	// run
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitBackward); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitBackward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(time.Duration(sec) * time.Second)
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitRelease); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitRelease); err != nil {
		return
	}
	return
}

// Forward runs both DC-Motors forward for the given amount of time in seconds.
func (bot *Robot) Forward(sec int) (err error) {
	motorPort := 0
	motorStarboard := 1
	var speed int32 = 255 // 255 = full speed!
	if err = bot.adafruit.SetDCMotorSpeed(motorPort, speed); err != nil {
		return
	}
	if err = bot.adafruit.SetDCMotorSpeed(motorStarboard, speed); err != nil {
		return
	}
	// run FORWARD
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitForward); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitForward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(time.Duration(sec) * time.Second)
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitRelease); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitRelease); err != nil {
		return
	}
	return
}

// DCMotorRunner is simply a test runner for the given motor
func (bot *Robot) DCMotorRunner(dcMotor int) (err error) {

	//log.Printf("%s\tRun Loop...\n", time.Now().String())
	// set the speed:
	var speed int32 = 255 // 255 = full speed!
	if err = bot.adafruit.SetDCMotorSpeed(dcMotor, speed); err != nil {
		return
	}
	// run FORWARD
	if err = bot.adafruit.RunDCMotor(dcMotor, i2c.AdafruitForward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(2000 * time.Millisecond)
	if err = bot.adafruit.RunDCMotor(dcMotor, i2c.AdafruitRelease); err != nil {
		return
	}
	// run BACKWARD
	if err = bot.adafruit.RunDCMotor(dcMotor, i2c.AdafruitBackward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(2000 * time.Millisecond)
	if err = bot.adafruit.RunDCMotor(dcMotor, i2c.AdafruitRelease); err != nil {
		return
	}
	return
}
