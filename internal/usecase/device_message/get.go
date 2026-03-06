package device_message

import (
	"context"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/repository"
)

type GetDeviceMessages struct {
	messageRepo repository.DeviceMessageRepository
}

func NewGetDeviceMessages(messageRepo repository.DeviceMessageRepository) *GetDeviceMessages {
	return &GetDeviceMessages{messageRepo: messageRepo}
}

type GetDeviceMessagesInput struct {
	GUID   string
	Offset int
	Limit  int
}

func (uc *GetDeviceMessages) Execute(ctx context.Context, input GetDeviceMessagesInput) ([]*domain.DeviceMessage, int, error) {
	const op = "usecase.get_device_messages"
	var err error
	res, total, err := uc.messageRepo.GetByDeviceGUID(ctx, input.GUID, input.Offset, input.Limit)
	if err != nil {
		return nil, 0, errs.Wrap(op, err)
	}
	return res, total, nil
}
