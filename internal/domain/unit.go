package domain

import "time"

type Unit struct {
	GUID        string
	Invid       string
	MQTT        string
	ProcessedAt time.Time
	Status      string
}
