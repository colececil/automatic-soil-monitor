.DEFAULT_GOAL := flash
.PHONY: flash build

include .env

TINYGO_ARGS = \
	-ldflags \
		"-X main.broadcastIntervalSetting=${BROADCAST_INTERVAL} \
		-X main.sensorPinsSetting=${SENSOR_PINS} \
		-X main.sensorDryCalibrationsSetting=${SENSOR_DRY_CALIBRATIONS} \
		-X main.sensorWetCalibrationsSetting=${SENSOR_WET_CALIBRATIONS}" \
	-target ${MICROCONTROLLER_TYPE} \
	-size full \
	./cmd/automatic_soil_monitor

flash:
	@echo "Building and flashing program..."
	tinygo flash -monitor $(TINYGO_ARGS)

build:
	@echo "Building program..."
	tinygo build $(TINYGO_ARGS)