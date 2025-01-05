package bluetooth_broadcast

import (
	"fmt"
	"github.com/colececil/automatic-soil-monitor/internal/moisture_data"
	"time"
	"tinygo.org/x/bluetooth"
)

var btHomeServiceUuid = bluetooth.New16BitUUID(0xFCD2)

const btHomeDeviceInformation = uint8(0x40)
const btHomeMoistureObjectId = uint8(0x2F)

// BluetoothBroadcast represents a bluetooth adapter that broadcasts data from an associated MoistureData instance,
// using the BTHome protocol.
type BluetoothBroadcast struct {
	advertisement     *bluetooth.Advertisement
	moistureData      *moisture_data.MoistureData
	broadcastInterval time.Duration
	isRunning         bool
}

// New creates a new BluetoothBroadcast instance. moistureData should be set to a MoistureData instance that contains
// the data to broadcast.
func New(moistureData *moisture_data.MoistureData, broadcastInterval time.Duration) (*BluetoothBroadcast, error) {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return nil, fmt.Errorf("failed to enable BLE adapter: %w", err)
	}

	advertisement := adapter.DefaultAdvertisement()

	return &BluetoothBroadcast{
		advertisement:     advertisement,
		moistureData:      moistureData,
		broadcastInterval: broadcastInterval,
	}, nil
}

// AdvertiseLatestData starts advertising the latest data from the MoistureData instance. If it's already advertising
// when this is called, it will be stopped and restarted.
func (b *BluetoothBroadcast) AdvertiseLatestData() error {
	var err error

	if b.isRunning {
		err = b.advertisement.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop BLE advertisement: %w", err)
		}
	}

	err = b.advertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName: "soil-monitor",
		ServiceData: []bluetooth.ServiceDataElement{
			{
				UUID: btHomeServiceUuid,
				Data: b.getBtHomeData(),
			},
		},
		Interval: bluetooth.NewDuration(b.broadcastInterval),
	})
	if err != nil {
		return fmt.Errorf("failed to configure BLE advertisement: %w", err)
	}

	err = b.advertisement.Start()
	if err != nil {
		return fmt.Errorf("failed to start BLE advertisement: %w", err)
	}

	fmt.Println("Started BLE advertisement using latest data.")
	b.isRunning = true
	return nil
}

// getBtHomeData constructs BTHome data using the latest data from the MoistureData instance.
func (b *BluetoothBroadcast) getBtHomeData() []byte {
	btHomeData := make([]byte, 1+5*b.moistureData.NumSensors())
	btHomeData[0] = btHomeDeviceInformation
	for i := 0; i < b.moistureData.NumSensors(); i++ {
		pos := i * 2
		btHomeData[pos+1] = btHomeMoistureObjectId
		btHomeData[pos+2] = b.moistureData.LatestReadingAsPercentage(i)
	}
	return btHomeData
}
