.DEFAULT_GOAL := flash
.PHONY: flash build check-env

include .env

TINYGO_ARGS = \
	-ldflags \
		"-X main.updateIntervalSetting=${UPDATE_INTERVAL} \
		-X main.sensorPinsSetting=${SENSOR_PINS} \
		-X main.sensorDryCalibrationsSetting=${SENSOR_DRY_CALIBRATIONS} \
		-X main.sensorWetCalibrationsSetting=${SENSOR_WET_CALIBRATIONS}" \
	-target ${MICROCONTROLLER_TYPE} \
	-size full \
	./cmd/automatic_soil_monitor

check-env:
ifndef UPDATE_INTERVAL
	@echo "The UPDATE_INTERVAL environment variable must be set."
	@exit 1
endif
ifndef SENSOR_PINS
	@echo "The SENSOR_PINS environment variable must be set."
	@exit 1
endif
ifndef SENSOR_DRY_CALIBRATIONS
	@echo "The SENSOR_DRY_CALIBRATIONS environment variable must be set."
	@exit 1
endif
ifndef SENSOR_WET_CALIBRATIONS
	@echo "The SENSOR_WET_CALIBRATIONS environment variable must be set."
	@exit 1
endif
ifndef MICROCONTROLLER_TYPE
	@echo "The MICROCONTROLLER_TYPE environment variable must be set."
	@exit 1
endif

flash: check-env
	@echo "Building and flashing program..."
	tinygo flash -monitor $(TINYGO_ARGS)

build: check-env
	@echo "Building program..."
	tinygo build $(TINYGO_ARGS)