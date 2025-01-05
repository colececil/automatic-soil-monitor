package bluetooth_broadcast

import (
	"fmt"
	"github.com/colececil/automatic-soil-monitor/internal/moisture_data"
	"github.com/hybridgroup/go-bthome"
	"time"
	"tinygo.org/x/bluetooth"
)

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

	btHomeData, err := b.getBtHomeData()
	if err != nil {
		return fmt.Errorf("failed to construct BTHome service data: %w", err)
	}

	err = b.advertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName:   "soil-monitor",
		ServiceData: []bluetooth.ServiceDataElement{btHomeData},
		Interval:    bluetooth.NewDuration(b.broadcastInterval),
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

// getBtHomeData constructs BTHome service data using the latest data from the MoistureData instance.
func (b *BluetoothBroadcast) getBtHomeData() (bluetooth.ServiceDataElement, error) {
	payload := &bthome.Payload{}
	for i := 0; i < b.moistureData.NumSensors(); i++ {
		err := payload.AddData(
			bthome.DataValue{
				Type:  bthome.Moisture8,
				Value: []byte{b.moistureData.LatestReadingAsPercentage(i)},
			},
		)
		if err != nil {
			return bluetooth.ServiceDataElement{}, err
		}
	}
	return payload.ServiceData(), nil
}
