package pdf_document

import (
	"fmt"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2/pkg/repository"

	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"

	"github.com/scmbr/device-tsv-processor/internal/domain"
)

type PDFGenerator struct {
	fontPath string
}

func NewPDFGenerator(fontPath string) *PDFGenerator {
	return &PDFGenerator{
		fontPath: fontPath,
	}
}

func (g *PDFGenerator) Generate(path string, messages []*domain.DeviceMessage) error {

	fontName := "custom"

	customFonts, err := repository.New().
		AddUTF8Font(fontName, fontstyle.Normal, g.fontPath).
		AddUTF8Font(fontName, fontstyle.Bold, g.fontPath).
		AddUTF8Font(fontName, fontstyle.Italic, g.fontPath).
		AddUTF8Font(fontName, fontstyle.BoldItalic, g.fontPath).
		Load()
	if err != nil {
		return err
	}

	cfg := config.NewBuilder().
		WithCustomFonts(customFonts).
		WithDefaultFont(&props.Font{
			Family: fontName,
		}).
		Build()

	m := maroto.New(cfg)
	m.AddRow(12,
		text.NewCol(12, "ОТЧЁТ ПО СООБЩЕНИЯМ УСТРОЙСТВА", props.Text{
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	m.AddRow(6,
		text.NewCol(6, fmt.Sprintf("Дата: %s", time.Now().Format(time.RFC3339)), props.Text{Size: 9}),
		text.NewCol(6, fmt.Sprintf("Всего: %d", len(messages)), props.Text{Align: align.Right, Size: 9}),
	)
	m.AddRow(4)
	grey := props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}

	header := m.AddRow(7,
		text.NewCol(2, "Device", headerStyle()),
		text.NewCol(1, "Lvl", headerStyle()),
		text.NewCol(2, "Class", headerStyle()),
		text.NewCol(4, "Message", headerStyle()),
		text.NewCol(2, "Created", headerStyle()),
	)

	header.WithStyle(&props.Cell{
		BackgroundColor: &grey,
	})
	for _, msg := range messages {
		m.AddRow(6,
			text.NewCol(2, fmt.Sprintf("%d", msg.DeviceID), body()),
			text.NewCol(1, fmt.Sprintf("%d", msg.Level), body()),
			text.NewCol(2, msg.Class, body()),
			text.NewCol(4, msg.Text, body()),
			text.NewCol(2, msg.CreatedAt.Format("2006-01-02 15:04"), body()),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return err
	}

	return document.Save(path)
}

func headerStyle() props.Text {
	return props.Text{
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  9,
	}
}

func body() props.Text {
	return props.Text{
		Size: 8,
	}
}
