package domain

import "time"

type JobStatus string

const (
    JobStatusPending   JobStatus = "pending"
    JobStatusAssigned  JobStatus = "assigned"
    JobStatusCompleted JobStatus = "completed"
    JobStatusFailed    JobStatus = "failed"
)

type Job struct {
    ID         int64
    SiteID     int64
    AgentID    int64
    Status     JobStatus
    CreatedAt  time.Time
    AssignedAt *time.Time
    DeadlineAt *time.Time
    FinishedAt *time.Time
}