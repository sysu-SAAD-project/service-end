package controller

import (
	"time"
)

// ActivityList defines the return format
type ActivityList struct {
	Content []ActivityIntroduction `json:"content"`
}

// ErrorMessage defines error format
type ErrorMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"msg"`
}

// ActivityIntroduction include required information in activity list page
type ActivityIntroduction struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
	Campus    int        `json:"campus"`
	Type      int        `json:"type"`
	Poster    string     `json:"poster"`
	Location  string     `json:"location"`
}
