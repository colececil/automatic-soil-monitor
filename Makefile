.DEFAULT_GOAL := flash
.PHONY: flash build

include .env

TINYGO_ARGS = \
	-ldflags \
		"-X main.minMoistureLevelString=${MIN_MOISTURE_LEVEL} -X main.maxMoistureLevelString=${MAX_MOISTURE_LEVEL}" \
	-target ${MICROCONTROLLER_TYPE} \
	-size full \
	./cmd/automatic_soil_monitor

flash:
	@echo "Building and flashing program..."
	tinygo flash -monitor $(TINYGO_ARGS)

build:
	@echo "Building program..."
	tinygo build $(TINYGO_ARGS)