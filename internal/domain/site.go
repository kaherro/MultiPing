package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Site struct {
	ID          int64     `json:"id"`
	URL         string    `json:"url"`
	IntervalSec int       `json:"interval_sec"`
	TimeoutMS   int       `json:"timeout_ms"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
}
