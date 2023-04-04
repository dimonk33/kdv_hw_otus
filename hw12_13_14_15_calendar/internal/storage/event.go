package storage

import "time"

type Event struct {
	ID          int64
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	OwnUserId   int
}
