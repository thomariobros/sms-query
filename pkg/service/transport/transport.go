package transport

import (
	"time"
)

const (
	StatusUnknownLine      = "UNKNOWN_LINE"
	StatusUnknownStop      = "UNKNOWN_STOP"
	StatusUnknownDirection = "UNKNOWN_DIRECTION"
	StatusNoService        = "NO_SERVICE"
	StatusOk               = "OK"
	StatusError            = "ERROR"
)

type TransportService interface {
	GetNextDepartureTime(lineType, lineID, stopName, direction string, max int) (*NextDeparture, error)
}

type NextDeparture struct {
	Status    string
	Now       time.Time
	StopTimes []DateTimeOrMoreThan1Hour
}

type DateTimeOrMoreThan1Hour struct {
	DateTime      time.Time
	MoreThan1Hour bool
}
