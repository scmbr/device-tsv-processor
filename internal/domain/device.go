package domain

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type Device struct {
	ID          int64
	GUID        string
	InvID       string
	MQTT        string
	ProcessedAt *time.Time
	Status      string
	CreatedAt   time.Time
}

func NewDevice(guid string, invID string, mqtt string, status string) (*Device, error) {
	device := &Device{
		GUID:      guid,
		InvID:     invID,
		MQTT:      mqtt,
		Status:    status,
		CreatedAt: time.Now(),
	}

	if err := device.Validate(); err != nil {
		return nil, err
	}
	return device, nil
}

func (d *Device) Validate() error {
	const op = "device.entity.validate"
	fields := map[string]string{}

	if d.GUID == "" {
		fields["guid"] = "is required"
	}
	if d.InvID == "" {
		fields["inv_id"] = "is required"
	}
	if d.Status == "" {
		fields["status"] = "is required"
	}

	if len(fields) > 0 {
		return errs.E(
			errs.KindInvalid,
			"DEVICE_INVALID",
			op,
			"invalid device",
			fields,
			nil,
		)
	}
	return nil
}
