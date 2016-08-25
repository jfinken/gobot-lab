# gobot lab

Experiments in Go Powered Robotics

## Requirements

 * [gobot.io](https://gobot.io/), particularly [this fork](https://github.com/jfinken/gobot) which includes Adafruit Motor HAT support for the Raspberry Pi.

## Sample compilation for ARMv7:

    GOARM=7 GOARCH=arm GOOS=linux go build -v main.go
