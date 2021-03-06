package adabot

import (
	"log"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	// Min pulse length out of 4096
	servoMin = 150
	// Max pulse length out of 4096
	servoMax = 700
	// Limiting the max this servo can rotate (in deg)
	maxDegree = 180
	// Number of degrees to increase per call
	degIncrease       = 10
	yawDeg            = 90
	pitchDeg          = 90
	yawChannel   byte = 1
	pitchChannel byte = 2
)

func degree2pulse(deg int) int32 {
	pulse := servoMin
	pulse += ((servoMax - servoMin) / maxDegree) * deg
	return int32(pulse)
}

// Robot defines a type abstracting the unexported driver type.
type Robot struct {
	adafruit *i2c.AdafruitMotorHatDriver
}

// NewRobot constructs and initializes an unexported driver object.
func NewRobot() (*Robot, error) {

	// Now in gobot.io 1.0: Metal Gobot
	// 	when you want to use the individual Gobot packages yourself to have the
	//  greatest control, or to more easily integrate Gobot functionality into
	//  your existing Golang programs.
	r := raspi.NewAdaptor()
	r.Connect()
	adaFruit := i2c.NewAdafruitMotorHatDriver(r)
	adaFruit.Start()

	/*
		// Custom init for attached servo hat and motors
		// Changing from the default 0x40 address because this configuration involves
		// a Servo HAT stacked on top of a DC/Stepper Motor HAT on top of the Pi.
		stackedHatAddr := 0x41

		// update the I2C address state
		adaFruit.SetServoHatAddress(stackedHatAddr)

		freq := 60.0
		if err := adaFruit.SetServoMotorFreq(freq); err != nil {
			return nil, err
		}
		// start in the middle of the 180-deg range in both yaw and pitch
		pulse := degree2pulse(yawDeg)
		if err := adaFruit.SetServoMotorPulse(yawChannel, 0, pulse); err != nil {
			return nil, err
		}
		pulse = degree2pulse(pitchDeg)
		if err := adaFruit.SetServoMotorPulse(pitchChannel, 0, pulse); err != nil {
			return nil, err
		}
	*/
	return &Robot{adafruit: adaFruit}, nil
}

// Stop releases both DC-Motors.  Stop that shizzle
func (bot *Robot) Stop() (err error) {
	// NOTE: these equate to the HAT "port", see Google docs wiring diagram
	motorPort := 3
	motorStarboard := 2
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitRelease); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitRelease); err != nil {
		return
	}
	return
}

// Left runs both DC-Motors in opposite directions
func (bot *Robot) Left(sec int) (err error) {
	// NOTE: these equate to the HAT "port", see Google docs wiring diagram
	motorPort := 3
	motorStarboard := 2
	var speed int32 = 255 // 255 = full speed!
	if err = bot.adafruit.SetDCMotorSpeed(motorPort, speed); err != nil {
		return
	}
	if err = bot.adafruit.SetDCMotorSpeed(motorStarboard, speed); err != nil {
		return
	}
	//------------------------------------
	// BUG: direction or wiring is flipped
	//------------------------------------
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitForward); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitBackward); err != nil {
		return
	}
	return
}

// Right runs both DC-Motors in opposite directions
func (bot *Robot) Right(sec int) (err error) {
	// NOTE: these equate to the HAT "port", see Google docs wiring diagram
	motorPort := 3
	motorStarboard := 2
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
	return
}

// Backward runs both DC-Motors backward
func (bot *Robot) Backward(sec int) (err error) {
	// NOTE: these equate to the HAT "port", see Google docs wiring diagram
	motorPort := 3
	motorStarboard := 2
	var speed int32 = 255 // 255 = full speed!
	if err = bot.adafruit.SetDCMotorSpeed(motorPort, speed); err != nil {
		return
	}
	if err = bot.adafruit.SetDCMotorSpeed(motorStarboard, speed); err != nil {
		return
	}
	// run
	if err = bot.adafruit.RunDCMotor(motorPort, i2c.AdafruitForward); err != nil {
		return
	}
	if err = bot.adafruit.RunDCMotor(motorStarboard, i2c.AdafruitForward); err != nil {
		return
	}
	return
}

// Forward runs both DC-Motors forward.
func (bot *Robot) Forward(sec int) (err error) {
	// NOTE: these equate to the HAT "port", see Google docs wiring diagram
	motorPort := 3
	motorStarboard := 2
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
	return
}

// Pitch will rotate the vertical oriented servo up/down based on the sign of dir.
func (bot *Robot) Pitch(dir int) (err error) {
	var pulse int32
	if dir > 0 {
		pitchDeg -= degIncrease
		pulse = degree2pulse(pitchDeg)
	} else {
		pitchDeg += degIncrease
		pulse = degree2pulse(pitchDeg)
	}
	if err = bot.adafruit.SetServoMotorPulse(pitchChannel, 0, pulse); err != nil {
		log.Printf(err.Error())
		return
	}
	return
}

// Yaw will rotate the horizontal oriented servo left/right based on the sign of dir.
func (bot *Robot) Yaw(dir int) (err error) {

	var pulse int32
	if dir <= 0 {
		// DEC
		yawDeg -= degIncrease
		pulse = degree2pulse(yawDeg)
	} else {
		// INCR
		yawDeg += degIncrease
		pulse = degree2pulse(yawDeg)
	}
	if err = bot.adafruit.SetServoMotorPulse(yawChannel, 0, pulse); err != nil {
		log.Printf(err.Error())
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
