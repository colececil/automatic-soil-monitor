# automatic-soil-monitor

This is a microcontroller project designed to monitor and report the soil moisture level for one or more potted plants. It broadcasts the data via BLE at a set interval using the [BTHome protocol](https://bthome.io/), which makes it easy to track using a home automation platform like [Home Assistant](https://www.home-assistant.io/).

## Parts and assembly

![Photo of two capacitive moisture sensors soldered to an Arduino Nano 33 IoT. The moisture sensors are stuck into the soil of a potted cactus.](assets/automatic-soil-monitor.jpg)

I used an [Arduino Nano 33 IOT](https://store-usa.arduino.cc/products/arduino-nano-33-iot) and [these capacitive moisture sensors](https://www.amazon.com/gp/product/B07SYBSHGX/), but this project will work with different (but functionally equivalent) parts. You can use [any microcontroller compatible with TinyGo](https://tinygo.org/docs/reference/microcontrollers/), as long as it has a BLE chip supported by [Go Bluetooth](https://github.com/tinygo-org/bluetooth) (see the readme there for details). In addition, any moisture sensor that gives an analog output should work. The project will work with as many moisture sensors as your hardware allows, as long as the BLE data doesn't exceed the byte limit of the BLE spec.

Each capacitive moisture sensor has three wires: a power wire ("VCC"), a ground wire ("GND"), and a signal wire ("AOUT"). The following connections need to be made:

- Each moisture sensor's "VCC" wire needs to be connected to a power output pin on the microcontroller that matches the sensor's power specifications. (All the moisture sensors can share the same power output pin.)
- Each moisture sensor's "GND" wire needs to be connected to the microcontroller's ground pin. (All the moisture sensors can share the same ground pin.)
- Each moisture sensor's "AOUT" wire needs to be connected to a separate analog input pin on the microcontroller.

## How to run the project

This project is built using the [PlatformIO](https://platformio.org/) framework. To run the project, you'll first need to install either the [PlatformIO CLI](https://docs.platformio.org/en/latest/core/installation/index.html) or a [PlatformIO IDE extension](https://docs.platformio.org/en/latest/integration/ide/pioide.html). Once PlatformIO is set up and your microcontroller is plugged in, you can upload and run the project using the [CLI run command](https://docs.platformio.org/en/latest/core/userguide/cmd_run.html), or by doing so through your IDE.

### Prerequisites

In order to build the project and flash it to your microcontroller, you'll need the following:

- [Go v1.22.3+](https://go.dev/)
- [TinyGo v0.34+](https://tinygo.org/): Follow the install guide specific to your operating system [here](https://tinygo.org/getting-started/install/), and be sure to note the instructions specific to the microcontroller you're using. (You should also be able to find setup instructions specific to your microcontroller [here](https://tinygo.org/docs/reference/microcontrollers/).)
- Depending on your microcontroller's type of BLE chip, you may need to do some additional setup. See the [Go Bluetooth](https://github.com/tinygo-org/bluetooth) readme for details.

### Configuration

Before building, the project requires some configuration values to be set. You'll need to create a `.env` file at the project root, and then set the following environment variables:

- `MICROCONTROLLER_TYPE`: The type of microcontroller you're using. See [this page](https://tinygo.org/docs/reference/microcontrollers/machine/) for a list of valid microcontroller types. You should be able to find the right value to use for your microcontroller [here](https://tinygo.org/docs/reference/microcontrollers/). Once you navigate to the page for your microcontroller, do a search for `-target` to find the example command for flashing to the microcontroller. The value you'll need is the one directly following the `-target` flag.
- `BROADCAST_INTERVAL`: The interval at which the project will broadcast updated data via BLE. This should be a duration in the format used by Go's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration) function.
- `SENSOR_PINS`: The analog input pins on the microcontroller that will be used to read the moisture sensor values. This should be a comma-separated list of the pin numbers that your moisture sensors are connected to on your microcontroller. Note that the actual numbers that TinyGo uses are needed, not the pin names. You should be able to find these under your microcontroller's type [here](https://tinygo.org/docs/reference/microcontrollers/machine/).
- `SENSOR_DRY_CALIBRATIONS`: The "dry" calibration values for each moisture sensor. This should be a comma-separated list of integers, in the same order as the pins in `SENSOR_PINS`. See the [Calibration](#calibration) section below for more information.
- `SENSOR_WET_CALIBRATIONS`: The "wet" calibration values for each moisture sensor. This should be a comma-separated list of integers, in the same order as the pins in `SENSOR_PINS`. See the [Calibration](#calibration) section below for more information.

Here is an example of what a `.env` file for this project might look like:

```
MICROCONTROLLER_TYPE=arduino-nano33
BROADCAST_INTERVAL=1h
SENSOR_PINS=2,34
SENSOR_DRY_CALIBRATIONS=56656,56128
SENSOR_WET_CALIBRATIONS=31872,31232
```

### Building and Flashing

### Calibration