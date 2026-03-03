package domain

import "time"

type Device struct {
	GUID        string
	Invid       string
	MQTT        string
	ProcessedAt time.Time
	Status      string
}
