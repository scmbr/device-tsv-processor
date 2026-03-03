package usecase

import "github.com/scmbr/device-tsv-processor/internal/domain"

type DocumentGenerator interface {
	Generate(path string, messages []*domain.DeviceMessage) error
}
