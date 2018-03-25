package controller

import (
	"github.com/sysu-saad-project/service-end/core/models/entities"
)

type SortRule struct {
	direction string
	property string
	ignoreCase bool
	nullHandling string
	ascending bool
	descending bool
}

type ActivityList struct {
	content []entities.ActivityInfo
	last bool
	totalPages int
	totalElements int
	numberOfElements int
	sort []SortRule
	first bool
	size int
	number int
}