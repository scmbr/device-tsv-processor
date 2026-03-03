package domain

import (
	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type DeviceMessage struct {
	ID        int64
	GUID      string
	InvID     string
	MsgID     string
	Text      string
	Context   string
	Class     string
	Level     int
	Area      string
	Addr      string
	Block     string
	Type      string
	Bit       int
	InvertBit bool
}

func NewDeviceMessage(
	guid string,
	invID string,
	msgID string,
	text string,
	context string,
	class string,
	level int,
	area string,
	addr string,
	block string,
	typ string,
	bit int,
	invertBit bool,
) (*DeviceMessage, error) {
	msg := &DeviceMessage{
		GUID:      guid,
		InvID:     invID,
		MsgID:     msgID,
		Text:      text,
		Context:   context,
		Class:     class,
		Level:     level,
		Area:      area,
		Addr:      addr,
		Block:     block,
		Type:      typ,
		Bit:       bit,
		InvertBit: invertBit,
	}

	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

func (m *DeviceMessage) Validate() error {
	const op = "device_message.entity.validate"

	fields := map[string]string{}

	if m.GUID == "" {
		fields["guid"] = "is required"
	}
	if m.InvID == "" {
		fields["device_guid"] = "is required"
	}
	if m.MsgID == "" {
		fields["msg_id"] = "is required"
	}
	if m.Level < 0 {
		fields["level"] = "must be non-negative"
	}
	if m.Bit < 0 || m.Bit > 31 {
		fields["bit"] = "must be between 0 and 31"
	}

	if len(fields) > 0 {
		return errs.E(
			errs.KindInvalid,
			"DEVICE_MESSAGE_INVALID",
			op,
			"invalid device message",
			fields,
			nil,
		)
	}
	return nil
}
