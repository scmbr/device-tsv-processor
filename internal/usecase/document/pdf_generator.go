package document

import "github.com/scmbr/device-tsv-processor/internal/domain"

type PDFGenerator interface {
	Generate(unitGUID string, messages []*domain.DeviceMessage) error
}
