package domain

import "time"

type ProcessedFile struct {
	Filename    string
	ProcessedAt time.Time
	Status      string
}
