package main

import (
	"fmt"
	"github.com/colececil/automatic-soil-monitor/internal/bluetooth_broadcast"
	"github.com/colececil/automatic-soil-monitor/internal/moisture_data"
	"machine"
	"strconv"
	"time"
)

// The min and max moisture levels are set at build time using values in a .env file. See the readme for more
// information.
var minMoistureLevelString string
var maxMoistureLevelString string

var sensors [2]machine.ADC
var led machine.Pin
var ledPowerState bool
var minMoistureLevel uint16
var maxMoistureLevel uint16
var moistureData *moisture_data.MoistureData
var bluetoothBroadcast *bluetooth_broadcast.BluetoothBroadcast

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
	minAsUint64, err := strconv.ParseUint(minMoistureLevelString, 10, 16)
	if err != nil {
		err = fmt.Errorf("failed to parse min moisture level: %w", err)
		logErrorAndRestart(err)
	}
	minMoistureLevel = uint16(minAsUint64)

	maxAsUint64, err := strconv.ParseUint(maxMoistureLevelString, 10, 16)
	if err != nil {
		err = fmt.Errorf("failed to parse max moisture level: %w", err)
		logErrorAndRestart(err)
	}
	maxMoistureLevel = uint16(maxAsUint64)

	machine.InitADC()
	sensors[0] = machine.ADC{Pin: machine.PA02}
	sensors[0].Configure(machine.ADCConfig{})
	sensors[1] = machine.ADC{Pin: machine.PB02}
	sensors[1].Configure(machine.ADCConfig{})

	led = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Set(ledPowerState)

	moistureData = moisture_data.New(
		2,
		[]uint16{minMoistureLevel, minMoistureLevel},
		[]uint16{maxMoistureLevel, maxMoistureLevel},
	)
	bluetoothBroadcast, err = bluetooth_broadcast.New(moistureData)
	if err != nil {
		logErrorAndRestart(err)
	}
}

// readMoistureLevels reads and reports the moisture levels from the sensors.
func readMoistureLevels() {
	readMoistureLevel(0)
	readMoistureLevel(1)
	err := bluetoothBroadcast.SendAdvertisement()
	if err != nil {
		logErrorAndRestart(err)
	}
}

// readMoistureLevel reads and reports the moisture level from the given sensor.
func readMoistureLevel(sensorIndex int) {
	sensor := sensors[sensorIndex]
	name := "Sensor " + strconv.Itoa(sensorIndex+1)
	reading := sensor.Get()
	moistureData.UpdateReading(sensorIndex, reading)
	fmt.Printf("%s: %2d%% (%d)\n", name, moistureData.LatestReadingAsPercentage(sensorIndex), reading)
}

// toggleLed toggles the state of the LED.
func toggleLed() {
	led.Set(ledPowerState)
	ledPowerState = !ledPowerState
}

// logErrorAndRestart logs the given error and restarts the device.
func logErrorAndRestart(err error) {
	fmt.Println("Error:", err)
	fmt.Println("Restarting in 5 seconds...")
	time.Sleep(5 * time.Second)
	machine.CPUReset()
}
