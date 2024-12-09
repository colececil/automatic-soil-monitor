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

// These settings are initialized at build time using values in a .env file. See the readme for more information.
var broadcastIntervalSetting string
var sensorPinsSetting string
var sensorDryCalibrationsSetting string
var sensorWetCalibrationsSetting string

var broadcastInterval time.Duration
var sensors []machine.ADC
var led machine.Pin
var ledPowerState bool
var moistureData *moisture_data.MoistureData
var bluetoothBroadcast *bluetooth_broadcast.BluetoothBroadcast

func main() {
	initialize()
	for {
		toggleLed()
		err := broadcastCurrentMoistureLevels()
		if err != nil {
			logErrorAndRestart(err)
		}
		time.Sleep(broadcastInterval)
	}
}

// initialize initializes the variables and components necessary for the program to run.
func initialize() {
	var err error
	broadcastInterval, err = time.ParseDuration(broadcastIntervalSetting)
	if err != nil {
		err = fmt.Errorf("failed to parse broadcast duration: %w", err)
		logErrorAndRestart(err)
	}

	err = initializeSensors()
	if err != nil {
		logErrorAndRestart(err)
	}

	sensorDryCalibrations, err := getSensorCalibrations(sensorDryCalibrationsSetting)
	if err != nil {
		logErrorAndRestart(err)
	}
	sensorWetCalibrations, err := getSensorCalibrations(sensorWetCalibrationsSetting)
	if err != nil {
		logErrorAndRestart(err)
	}
	moistureData = moisture_data.New(
		len(sensors),
		sensorDryCalibrations,
		sensorWetCalibrations,
	)

	bluetoothBroadcast, err = bluetooth_broadcast.New(moistureData)
	if err != nil {
		logErrorAndRestart(err)
	}

	initializeLed()

	fmt.Printf("Initialization complete with %d moisture sensors.\n", len(sensors))
}

// initializeSensors initializes the pins for the moisture sensors, using the pins identified in the .env file. It
// returns an error if the string from the .env file does not match the expected format.
func initializeSensors() error {
	machine.InitADC()
	sensorPinStrings := strings.Split(sensorPinsSetting, ",")
	sensors = make([]machine.ADC, len(sensorPinStrings))
	for i, pinString := range sensorPinStrings {
		pinAsUint64, err := strconv.ParseUint(pinString, 10, 8)
		if err != nil {
			err = fmt.Errorf("failed to parse pin number for Sensor %d: %w", i+1, err)
			return err
		}
		sensors[i] = machine.ADC{Pin: machine.Pin(pinAsUint64)}
		sensors[i].Configure(machine.ADCConfig{})
	}
	return nil
}

// getSensorCalibrations parses the given string of comma-separated numbers and returns an equivalent slice of uint16
// values. It returns an error if the given string does not match the expected format, or if the number of calibrations
// does not match the number of sensors in the sensor slice.
func getSensorCalibrations(calibrationsString string) ([]uint16, error) {
	calibrationStrings := strings.Split(calibrationsString, ",")
	calibrations := make([]uint16, len(calibrationStrings))
	if len(calibrationStrings) != len(sensors) {
		err := fmt.Errorf("number of calibrations does not match number of sensors")
		return nil, err
	}
	for i, calibrationString := range calibrationStrings {
		calibration, err := strconv.ParseUint(calibrationString, 10, 16)
		if err != nil {
			err = fmt.Errorf("failed to parse calibration number for Sensor %d: %w", i+1, err)
			return nil, err
		}
		calibrations[i] = uint16(calibration)
	}
	return calibrations, nil
}

// broadcastCurrentMoistureLevels reads the current moisture levels from the sensors and broadcasts them via BLE. It
// returns an error if there is an issue broadcasting the data.
func broadcastCurrentMoistureLevels() error {
	fmt.Println("\nGetting updated data from moisture sensors...")
	for i := range sensors {
		sensor := sensors[i]
		reading := sensor.Get()
		moistureData.UpdateReading(i, reading)
	}
	return bluetoothBroadcast.SendAdvertisement()
}

// initializeLed initializes the LED pin.
func initializeLed() {
	led = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Set(ledPowerState)
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
