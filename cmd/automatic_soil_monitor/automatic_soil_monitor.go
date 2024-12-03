package main

import (
	"fmt"
	"machine"
	"strconv"
	"time"
	"tinygo.org/x/bluetooth"
)

// The min and max moisture levels are set at build time using values in a .env file. See the readme for more
// information.
var minMoistureLevelString string
var maxMoistureLevelString string

var btHomeUuid = bluetooth.New16BitUUID(0xFCD2)

const deviceInformation = uint8(0x40)
const moistureObjectId = uint8(0x2F)

var sensors [2]machine.ADC
var moisturePercentages [2]uint8
var led machine.Pin
var ledPowerState bool
var bleAdapter *bluetooth.Adapter
var bleAdvertisement *bluetooth.Advertisement
var bleServiceData []bluetooth.ServiceDataElement
var btHomeData []byte
var minMoistureLevel int
var maxMoistureLevel int

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
	var err error
	minMoistureLevel, err = strconv.Atoi(minMoistureLevelString)
	if err != nil {
		fmt.Println("Failed to parse min moisture level:", err)
		restart()
	}
	maxMoistureLevel, err = strconv.Atoi(maxMoistureLevelString)
	if err != nil {
		fmt.Println("Failed to parse max moisture level:", err)
		restart()
	}

	machine.InitADC()
	sensors[0] = machine.ADC{Pin: machine.PA02}
	sensors[0].Configure(machine.ADCConfig{})
	sensors[1] = machine.ADC{Pin: machine.PB02}
	sensors[1].Configure(machine.ADCConfig{})

	led = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Set(ledPowerState)

	bleAdapter = bluetooth.DefaultAdapter
	err = bleAdapter.Enable()
	if err != nil {
		fmt.Println("Failed to enable BLE adapter:", err)
		restart()
	}

	btHomeData = []byte{
		deviceInformation,
		moistureObjectId,
		moisturePercentages[0],
		moistureObjectId,
		moisturePercentages[1],
	}

	bleAdvertisement = bleAdapter.DefaultAdvertisement()
	bleServiceData = []bluetooth.ServiceDataElement{
		{
			UUID: btHomeUuid,
			Data: btHomeData,
		},
	}
	err = bleAdvertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName:   "soil-monitor",
		ServiceData: bleServiceData,
	})
	if err != nil {
		fmt.Println("Failed to configure BLE advertisement:", err)
		restart()
	}

	err = bleAdvertisement.Start()
	if err != nil {
		fmt.Println("Failed to start BLE advertisement:", err)
		restart()
	}
}

// updateServiceData updates the service data with the current moisture levels.
func updateServiceData() {
	btHomeData[2] = moisturePercentages[0]
	btHomeData[4] = moisturePercentages[1]
	err := bleAdvertisement.SetServiceData(bleServiceData)
	if err != nil {
		fmt.Println("Failed to update BLE service data:", err)
		restart()
	}

	fmt.Print("Service data updated with data: \"")
	for i, dataByte := range bleServiceData[0].Data {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%02X", dataByte)
	}
	fmt.Println("\"")
}

// readMoistureLevels reads and reports the moisture levels from the sensors.
func readMoistureLevels() {
	readMoistureLevel(0)
	readMoistureLevel(1)
	updateServiceData()
}

// readMoistureLevel reads and reports the moisture level from the given sensor with the given name.
func readMoistureLevel(sensorIndex uint8) {
	sensor := sensors[sensorIndex]
	name := "Sensor " + strconv.Itoa(int(sensorIndex+1))
	reading := sensor.Get()
	moisturePercentages[sensorIndex] = calculatePercentage(reading)
	fmt.Printf("%s: %2d%% (%d)\n", name, moisturePercentages[sensorIndex], reading)
}

// calculatePercentage calculates the percentage of the given value between the min and max moisture levels.
func calculatePercentage(value uint16) uint8 {
	if maxMoistureLevel < minMoistureLevel {
		return uint8((float64(int(value)-minMoistureLevel) / float64(maxMoistureLevel-minMoistureLevel)) * 100)
	}
	return uint8((float64(maxMoistureLevel-int(value)) / float64(maxMoistureLevel-minMoistureLevel)) * 100)
}

// toggleLed toggles the state of the LED.
func toggleLed() {
	led.Set(ledPowerState)
	ledPowerState = !ledPowerState
}

// restart restarts the device.
func restart() {
	fmt.Println("Restarting in 5 seconds...")
	time.Sleep(5 * time.Second)
	machine.CPUReset()
}
