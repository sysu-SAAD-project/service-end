package controller

import (
	"github.com/sysu-saad-project/service-end/models/entities"
)

// ActivityList defines the return format
type ActivityList struct {
	Content []entities.ActivityInfo `json:"content"`
}

// ErrorMessage defines error format
type ErrorMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"msg"`
}
