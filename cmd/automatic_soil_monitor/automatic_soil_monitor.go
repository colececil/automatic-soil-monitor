package main

import (
	"machine"
	"time"
	"tinygo.org/x/bluetooth"
)

var sensor1 machine.ADC
var sensor2 machine.ADC
var led machine.Pin
var ledPowerState bool
var bleAdapter *bluetooth.Adapter

func main() {
	initialize()
	for {
		toggleLed()
		readMoistureLevels()
		time.Sleep(time.Second)
	}
}

// initialize initializes the necessary components.
func initialize() {
	machine.InitADC()
	sensor1 = machine.ADC{Pin: machine.PA02}
	sensor1.Configure(machine.ADCConfig{})
	sensor2 = machine.ADC{Pin: machine.PB02}
	sensor2.Configure(machine.ADCConfig{})

	led = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Set(ledPowerState)

	bleAdapter = bluetooth.DefaultAdapter
	err := bleAdapter.Enable()
	if err != nil {
		println("Failed to enable BLE adapter:", err)
		restart()
	}

	bleAdvertisement := bleAdapter.DefaultAdvertisement()
	err = bleAdvertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName: "automatic-soil-monitor",
		Interval:  bluetooth.NewDuration(5 * time.Second),
	})
	if err != nil {
		println("Failed to configure BLE advertisement:", err)
		restart()
	}

	err = bleAdvertisement.Start()
	if err != nil {
		println("Failed to start BLE advertisement:", err)
		restart()
	}
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

// restart restarts the device.
func restart() {
	println("Restarting in 5 seconds...")
	time.Sleep(5 * time.Second)
	machine.CPUReset()
}
