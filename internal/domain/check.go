package domain

import "time"

type CheckResult struct {
	SiteID    int64
	OK        bool
	HTTPCode  int
	LatencyMS int
	Error     error
	CheckedAT time.Time
}
