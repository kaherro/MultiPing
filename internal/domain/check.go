package domain

import "time"

type Checker struct {
	SiteID    int64
	OK        bool
	HTTPCode  int
	LatencyMS int
	Error     error
	CheckedAT time.Time
}
