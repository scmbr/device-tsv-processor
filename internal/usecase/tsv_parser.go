package usecase

import "github.com/scmbr/device-tsv-processor/internal/domain"

type TSVParser interface {
	Parse(path string) ([]domain.DeviceMessage, []domain.ParseError, error)
}
