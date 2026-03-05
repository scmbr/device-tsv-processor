package pdf_document

import (
	"fmt"
	"os"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) Generate(path string, messages []*domain.DeviceMessage) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, m := range messages {
		_, err := f.WriteString(fmt.Sprintf("%+v\n", m))
		if err != nil {
			return err
		}
	}

	return nil
}
