package main

import (
	"machine"
	"time"
)

var sensor1 machine.ADC
var sensor2 machine.ADC
var led machine.Pin
var ledPowerState bool

func main() {
	initialize()
	for {
		toggleLed()
		readMoistureLevels()
		time.Sleep(time.Second)
	}
}

// initialize initializes the sensors and the LED.
func initialize() {
	machine.InitADC()
	sensor1 = machine.ADC{Pin: machine.PA02}
	sensor1.Configure(machine.ADCConfig{})
	sensor2 = machine.ADC{Pin: machine.PB02}
	sensor2.Configure(machine.ADCConfig{})

	led = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Set(ledPowerState)
}

// readMoistureLevels reads and reports the moisture levels from the sensors.
func readMoistureLevels() {
	readMoistureLevel(sensor1, "Sensor 1")
	readMoistureLevel(sensor2, "Sensor 2")
}

// readMoistureLevel reads and reports the moisture level from the given sensor with the given name.
func readMoistureLevel(input machine.ADC, name string) {
	reading := input.Get()
	println(name, ": ", reading)
}

// toggleLed toggles the state of the LED.
func toggleLed() {
	led.Set(ledPowerState)
	ledPowerState = !ledPowerState
}
