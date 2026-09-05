package domain

import "time"

type CronTaskResult struct {
	TaskName       string        `json:"task_name"`
	ProcessedCount int           `json:"processed_count"`
	SuccessCount   int           `json:"success_count"`
	ErrorCount     int           `json:"error_count"`
	Duration       time.Duration `json:"duration"`
	Errors         []string      `json:"errors,omitempty"`
}
