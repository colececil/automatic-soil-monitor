package bluetooth_broadcast

import (
	"fmt"
	"github.com/colececil/automatic-soil-monitor/internal/moisture_data"
	"tinygo.org/x/bluetooth"
)

var btHomeServiceUuid = bluetooth.New16BitUUID(0xFCD2)

const btHomeDeviceInformation = uint8(0x40)
const btHomeMoistureObjectId = uint8(0x2F)

// BluetoothBroadcast represents a bluetooth adapter that broadcasts data from an associated MoistureData instance,
// using the BTHome protocol.
type BluetoothBroadcast struct {
	adapter       *bluetooth.Adapter
	advertisement *bluetooth.Advertisement
	serviceData   []bluetooth.ServiceDataElement
	btHomeData    []byte
	moistureData  *moisture_data.MoistureData
}

// New creates a new BluetoothBroadcast instance. moistureData should be set to a MoistureData instance that contains
// the data to broadcast.
func New(moistureData *moisture_data.MoistureData) (*BluetoothBroadcast, error) {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return nil, fmt.Errorf("failed to enable BLE adapter: %w", err)
	}

	btHomeData := make([]byte, 1+5*moistureData.NumSensors())
	setBtHomeData(btHomeData, moistureData)

	serviceData := []bluetooth.ServiceDataElement{
		{
			UUID: btHomeServiceUuid,
			Data: btHomeData,
		},
	}

	advertisement := adapter.DefaultAdvertisement()
	err = advertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName:   "soil-monitor",
		ServiceData: serviceData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to configure BLE advertisement: %w", err)
	}

	err = advertisement.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start BLE advertisement: %w", err)
	}

	return &BluetoothBroadcast{
		adapter:       adapter,
		advertisement: advertisement,
		serviceData:   serviceData,
		btHomeData:    btHomeData,
		moistureData:  moistureData,
	}, nil
}

// SendAdvertisement sends a BLE advertisement containing the latest data from the MoistureData instance.
func (b *BluetoothBroadcast) SendAdvertisement() error {
	setBtHomeData(b.btHomeData, b.moistureData)
	err := b.advertisement.SetServiceData(b.serviceData)
	if err != nil {
		return fmt.Errorf("failed to update BLE service data: %w", err)
	}
	fmt.Println("BLE advertisement sent.")
	return nil
}

// setBtHomeData updates the btHomeData with the latest data from the MoistureData instance.
func setBtHomeData(btHomeData []byte, moistureData *moisture_data.MoistureData) {
	btHomeData[0] = btHomeDeviceInformation
	for i := 0; i < moistureData.NumSensors(); i++ {
		pos := i * 5
		btHomeData[pos+1] = btHomeMoistureObjectId
		btHomeData[pos+2] = moistureData.LatestReadingAsPercentage(i)
	}
}
