package moisture_data

import "math"

type MoistureData struct {
	latestReadings    []uint16
	minReadings       []*uint16
	maxReadings       []*uint16
	minMoistureLevels []uint16
	maxMoistureLevels []uint16
}

// New creates a new MoistureData instance. numSensors should be set to the number of moisture sensors connected to the
// device. minMoistureLevels and maxMoistureLevels should be set to the minimum and maximum moisture levels for each
// sensor, and are used for calculating moisture percentages.
func New(numSensors int, minMoistureLevels []uint16, maxMoistureLevels []uint16) *MoistureData {
	return &MoistureData{
		latestReadings:    make([]uint16, numSensors),
		minReadings:       make([]*uint16, numSensors),
		maxReadings:       make([]*uint16, numSensors),
		minMoistureLevels: minMoistureLevels,
		maxMoistureLevels: maxMoistureLevels,
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
}

// LatestReading returns the latest reading for the given sensor.
func (m *MoistureData) LatestReading(sensorIndex int) uint16 {
	return m.latestReadings[sensorIndex]
}

// LatestReadingAsPercentage returns the latest reading for the given sensor as a percentage.
func (m *MoistureData) LatestReadingAsPercentage(sensorIndex int) uint8 {
	minLevel := int(m.minMoistureLevels[sensorIndex])
	maxLevel := int(m.maxMoistureLevels[sensorIndex])
	reading := int(m.LatestReading(sensorIndex))

	var percentage float64
	if maxLevel < minLevel {
		percentage = (float64(reading-minLevel) / float64(maxLevel-minLevel)) * 100
	} else {
		percentage = (float64(maxLevel-reading) / float64(maxLevel-minLevel)) * 100
	}

	percentage = math.Max(percentage, 0)
	percentage = math.Min(percentage, 100)
	return uint8(math.Round(percentage))
}

// MinReading returns the minimum reading for the given sensor.
func (m *MoistureData) MinReading(sensorIndex int) uint16 {
	return *m.minReadings[sensorIndex]
}

// MaxReading returns the maximum reading for the given sensor.
func (m *MoistureData) MaxReading(sensorIndex int) uint16 {
	return *m.maxReadings[sensorIndex]
}
