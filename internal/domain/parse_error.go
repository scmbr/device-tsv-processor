package domain

import (
	"time"

	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type ParseError struct {
	Filename   string
	LineNumber int
	Message    string
	CreatedAt  time.Time
}

func NewParseError(filename string, lineNumber int, message string) (*ParseError, error) {
	pe := &ParseError{
		Filename:   filename,
		LineNumber: lineNumber,
		Message:    message,
		CreatedAt:  time.Now(),
	}

	if err := pe.Validate(); err != nil {
		return nil, err
	}
	return pe, nil
}

func (p *ParseError) Validate() error {
	const op = "parse_error.entity.validate"

	fields := map[string]string{}

	if p.Filename == "" {
		fields["filename"] = "is required"
	}
	if p.LineNumber <= 0 {
		fields["line_number"] = "must be positive"
	}
	if p.Message == "" {
		fields["message"] = "is required"
	}

	if len(fields) > 0 {
		return errs.E(
			errs.KindInvalid,
			"PARSE_ERROR_INVALID",
			op,
			"invalid parse error",
			fields,
			nil,
		)
	}

	return nil
}
