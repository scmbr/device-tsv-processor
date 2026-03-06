package pdf_document

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) Generate(path string, messages []*domain.DeviceMessage) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 10)
	pdf.AddPage()

	for i, m := range messages {
		pdf.CellFormat(0, 6, formatMessage(i+1, m), "", 1, "", false, 0, "")
	}

	return pdf.OutputFileAndClose(path)
}

func formatMessage(index int, m *domain.DeviceMessage) string {
	return fmt.Sprintf(
		"%d) ID: %d | DeviceID: %d | MQTT: %s | InvID: %s | UnitGUID: %s\n"+
			"   MsgID: %s | Text: %s | Context: %s | Class: %s | Level: %d\n"+
			"   Area: %s | Addr: %s | Block: %s | Type: %s | Bit: %d | InvertBit: %v\n"+
			"   CreatedAt: %s\n",
		index,
		m.ID, m.DeviceID, m.MQTT, m.InvID, m.UnitGUID,
		m.MsgID, m.Text, m.Context, m.Class, m.Level,
		m.Area, m.Addr, m.Block, m.Type, m.Bit, m.InvertBit,
		m.CreatedAt.Format(time.RFC3339),
	)
}
