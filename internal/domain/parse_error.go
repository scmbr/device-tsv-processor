package domain

import "time"

type ParseError struct {
	Filename   string
	LineNumber int
	Message    string
	CreatedAt  time.Time
}
