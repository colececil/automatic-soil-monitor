package main

import (
	"fmt"
	"github.com/colececil/automatic-soil-monitor/internal/bluetooth_broadcast"
	"github.com/colececil/automatic-soil-monitor/internal/moisture_data"
	"machine"
	"strconv"
	"strings"
	"time"
)

// The sensor dry and wet calibrations are set at build time using values in a .env file. See the readme for more
// information.
var sensorPins string
var sensorDryCalibrations string
var sensorWetCalibrations string

var sensors []machine.ADC
var led machine.Pin
var ledPowerState bool
var moistureData *moisture_data.MoistureData
var bluetoothBroadcast *bluetooth_broadcast.BluetoothBroadcast

func main() {
	initialize()
	for {
		toggleLed()
		readMoistureLevels()
		time.Sleep(time.Second) // Todo: Change duration to a more sensible value.
	}
}

// initialize initializes the necessary components.
func initialize() {
	machine.InitADC()
	sensorPinStrings := strings.Split(sensorPins, ",")
	sensors = make([]machine.ADC, len(sensorPinStrings))
	for i, pinString := range sensorPinStrings {
		pinAsUint64, err := strconv.ParseUint(pinString, 10, 8)
		if err != nil {
			err = fmt.Errorf("failed to parse pin number for Sensor %d: %w", i+1, err)
			logErrorAndRestart(err)
		}
		sensors[i] = machine.ADC{Pin: machine.Pin(pinAsUint64)}
		sensors[i].Configure(machine.ADCConfig{})
	}

	sensorDryCalibrationStrings := strings.Split(sensorDryCalibrations, ",")
	sensorDryCalibrations := make([]uint16, len(sensorDryCalibrationStrings))
	if len(sensorDryCalibrationStrings) != len(sensors) {
		err := fmt.Errorf("number of dry calibrations does not match number of sensors")
		logErrorAndRestart(err)
	}
	for i, calibrationString := range sensorDryCalibrationStrings {
		calibration, err := strconv.ParseUint(calibrationString, 10, 16)
		if err != nil {
			err = fmt.Errorf("failed to parse dry calibration number for Sensor %d: %w", i+1, err)
			logErrorAndRestart(err)
		}
		sensorDryCalibrations[i] = uint16(calibration)
	}

	sensorWetCalibrationStrings := strings.Split(sensorWetCalibrations, ",")
	sensorWetCalibrations := make([]uint16, len(sensorWetCalibrationStrings))
	if len(sensorWetCalibrationStrings) != len(sensors) {
		err := fmt.Errorf("number of wet calibrations does not match number of sensors")
		logErrorAndRestart(err)
	}
	for i, calibrationString := range sensorWetCalibrationStrings {
		calibration, err := strconv.ParseUint(calibrationString, 10, 16)
		if err != nil {
			err = fmt.Errorf("failed to parse wet calibration number for Sensor %d: %w", i+1, err)
			logErrorAndRestart(err)
		}
		sensorWetCalibrations[i] = uint16(calibration)
	}

	led = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Set(ledPowerState)

	moistureData = moisture_data.New(
		2,
		sensorDryCalibrations,
		sensorWetCalibrations,
	)
	var err error
	bluetoothBroadcast, err = bluetooth_broadcast.New(moistureData)
	if err != nil {
		logErrorAndRestart(err)
	}
}

// readMoistureLevels reads and reports the moisture levels from the sensors.
func readMoistureLevels() {
	for i := range sensors {
		readMoistureLevel(i)
	}
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
