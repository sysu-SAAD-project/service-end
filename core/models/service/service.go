package service

import (
	"github.com/sysu-saad-project/service-end/core/models/entities"
)

// GetActivityList return wanted activity list with given page number
func GetActivityList(pageNum int) []entities.ActivityInfo {
	activityList := make([]entities.ActivityInfo, 0)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	entities.Engine.Find(&activityList)
	return activityList
}

// GetActivityInfo return wanted activity detail information which is given by id
func GetActivityInfo(id int) entities.ActivityInfo {
	var activity entities.ActivityInfo

	entities.Engine.ID(id).Get(&activity)
	return activity
}
