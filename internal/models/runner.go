package models

import (
	"time"

	"reqium/internal/enums"
)

type RunResult struct {
	RequestName string
	Method      string
	URL         string
	StatusCode  int
	Duration    time.Duration
	Status      enums.RunnerStatus
	Error       string
}
