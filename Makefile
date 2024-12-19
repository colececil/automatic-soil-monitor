.DEFAULT_GOAL := flash
.PHONY: flash build

include .env

TINYGO_ARGS = \
	-ldflags \
		"-X main.sensorPins=${SENSOR_PINS} \
		-X main.sensorDryCalibrations=${SENSOR_DRY_CALIBRATIONS} \
		-X main.sensorWetCalibrations=${SENSOR_WET_CALIBRATIONS}" \
	-target ${MICROCONTROLLER_TYPE} \
	-size full \
	./cmd/automatic_soil_monitor

flash:
	@echo "Building and flashing program..."
	tinygo flash -monitor $(TINYGO_ARGS)

build:
	@echo "Building program..."
	tinygo build $(TINYGO_ARGS)