package parser

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
)

type TSVParser struct{}

func NewTSVParser() *TSVParser {
	return &TSVParser{}
}

func (p *TSVParser) Parse(
	ctx context.Context,
	path string,
) ([]*domain.DeviceMessage, []*domain.ParseError, error) {
	const op = "tsv_parser.parse"
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, errs.Wrap(op, err)
	}
	defer file.Close()

	var (
		messages []*domain.DeviceMessage
		errors   []*domain.ParseError
	)

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {

		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}

		lineNumber++
		line := scanner.Text()

		if lineNumber == 1 {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 15 {
			if pe, err := domain.NewParseError(path, lineNumber, "invalid column count"); err == nil {
				errors = append(errors, pe)
			} else {
				return nil, nil, fmt.Errorf("create parse error: %w", err)
			}
			continue
		}

		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}

		level := 0
		if l, err := strconv.Atoi(fields[8]); err == nil {
			level = l
		}

		bit := 0
		if b, err := strconv.Atoi(fields[13]); err == nil {
			bit = b
		}

		invertBit := false
		lower := strings.ToLower(fields[14])
		if lower == "1" || lower == "true" {
			invertBit = true
		}

		msg, err := domain.NewDeviceMessage(
			fields[3],  // UnitGUID
			fields[2],  // InvID
			fields[1],  // MQTT
			fields[4],  // MsgID
			fields[5],  // Text
			fields[6],  // Context
			fields[7],  // Class
			level,      // Level
			fields[9],  // Area
			fields[10], // Addr
			fields[11], // Block
			fields[12], // Type
			bit,        // Bit
			invertBit,  // InvertBit

		)
		if err != nil {
			if pe, e := domain.NewParseError(path, lineNumber, err.Error()); e == nil {
				errors = append(errors, pe)
			} else {
				return nil, nil, fmt.Errorf("create parse error: %w", e)
			}
			continue
		}

		messages = append(messages, msg)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scan file: %w", err)
	}

	return messages, errors, nil
}
