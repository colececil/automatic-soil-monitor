package moisture_data

import (
	"fmt"
	"math"
)

// MoistureData represents the moisture sensor data.
type MoistureData struct {
	latestReadings        []uint16
	minReadings           []*uint16
	maxReadings           []*uint16
	sensorDryCalibrations []uint16
	sensorWetCalibrations []uint16
}

// New creates a new MoistureData instance. numSensors should be set to the number of moisture sensors connected to the
// device. sensorDryCalibrations and sensorWetCalibrations should be set to the expected dry and wet values for each
// sensor - these are used for calculating moisture percentages.
func New(numSensors int, sensorDryCalibrations []uint16, sensorWetCalibrations []uint16) *MoistureData {
	return &MoistureData{
		latestReadings:        make([]uint16, numSensors),
		minReadings:           make([]*uint16, numSensors),
		maxReadings:           make([]*uint16, numSensors),
		sensorDryCalibrations: sensorDryCalibrations,
		sensorWetCalibrations: sensorWetCalibrations,
	}
}

// NumSensors returns the number of sensors.
func (m *MoistureData) NumSensors() int {
	return len(m.latestReadings)
}

// UpdateReading updates the reading for the given sensor.
func (m *MoistureData) UpdateReading(sensorIndex int, reading uint16) {
	m.latestReadings[sensorIndex] = reading

	if m.minReadings[sensorIndex] == nil || reading < *m.minReadings[sensorIndex] {
		m.minReadings[sensorIndex] = &reading
	}

	if m.maxReadings[sensorIndex] == nil || reading > *m.maxReadings[sensorIndex] {
		m.maxReadings[sensorIndex] = &reading
	}

	fmt.Printf("Reading updated for Sensor %d:\n", sensorIndex+1)
	fmt.Printf("  - Current reading: %3d%% (%d)\n", m.readingAsPercentage(sensorIndex, int(reading)), reading)
	fmt.Printf("  - Min reading:     %3d%% (%d)\n",
		m.readingAsPercentage(sensorIndex, int(*m.minReadings[sensorIndex])), *m.minReadings[sensorIndex])
	fmt.Printf("  - Max reading:     %3d%% (%d)\n",
		m.readingAsPercentage(sensorIndex, int(*m.maxReadings[sensorIndex])), *m.maxReadings[sensorIndex])
}

// LatestReading returns the latest reading for the given sensor.
func (m *MoistureData) LatestReading(sensorIndex int) uint16 {
	return m.latestReadings[sensorIndex]
}

// LatestReadingAsPercentage returns the latest reading for the given sensor as a percentage.
func (m *MoistureData) LatestReadingAsPercentage(sensorIndex int) uint8 {
	reading := m.LatestReading(sensorIndex)
	return m.readingAsPercentage(sensorIndex, int(reading))
}

// MinReading returns the minimum reading for the given sensor.
func (m *MoistureData) MinReading(sensorIndex int) uint16 {
	return *m.minReadings[sensorIndex]
}

// MaxReading returns the maximum reading for the given sensor.
func (m *MoistureData) MaxReading(sensorIndex int) uint16 {
	return *m.maxReadings[sensorIndex]
}

// readingAsPercentage returns the given reading as a percentage, for the given sensor.
func (m *MoistureData) readingAsPercentage(sensorIndex int, reading int) uint8 {
	dryCalibration := int(m.sensorDryCalibrations[sensorIndex])
	wetCalibration := int(m.sensorWetCalibrations[sensorIndex])

	var percentage float64
	if wetCalibration < dryCalibration {
		percentage = (float64(reading-dryCalibration) / float64(wetCalibration-dryCalibration)) * 100
	} else {
		percentage = (float64(wetCalibration-reading) / float64(wetCalibration-dryCalibration)) * 100
	}

	percentage = math.Max(percentage, 0)
	percentage = math.Min(percentage, 100)
	return uint8(math.Round(percentage))
}
