package domain

import "time"

type Agent struct {
    ID         int64
    Name       string
    TokenHash  string
    Enabled    bool
    LastSeenAt *time.Time
    CreatedAt  time.Time
}