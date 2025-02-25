package domain

import "time"

type Journey struct {
	zone      int
	direction string
	timestamp time.Time
}

// NewJourney creates a new instance of Journey
func NewJourney(zone int, direction string, timestamp time.Time) *Journey {
	return &Journey{
		zone:      zone,
		direction: direction,
		timestamp: timestamp,
	}
}

// GetZone returns the zone of the journey
func (j *Journey) GetZone() int {
	return j.zone
}

// GetDirection returns the direction of the journey
func (j *Journey) GetDirection() string {
	return j.direction
}

// GetTimestamp returns the timestamp of the journey
func (j *Journey) GetTimestamp() time.Time {
	return j.timestamp
}
