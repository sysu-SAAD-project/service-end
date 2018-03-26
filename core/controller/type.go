package controller

import (
	"github.com/sysu-saad-project/service-end/core/models/entities"
)

// SortRule defines the rule of sorting
type SortRule struct {
	Direction    string `json:"direction"`
	Property     string `json:"property"`
	IgnoreCase   bool   `json:"ignoreCase"`
	NullHandling string `json:"nullHandling"`
	Ascending    bool   `json:"ascending"`
	Descending   bool   `json:"descending"`
}

// ActivityList defines the return format
type ActivityList struct {
	Content          []entities.ActivityInfo `json:"content"`
	Last             bool                    `json:"last"`
	TotalPages       int                     `json:"totalPages"`
	TotalElements    int                     `json:"totalElements"`
	NumberOfElements int                     `json:"numberOfElements"`
	Sort             []SortRule              `json:"sort"`
	First            bool                    `json:"first"`
	Size             int                     `json:"size"`
	Number           int                     `json:"number"`
}
