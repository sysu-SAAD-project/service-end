package service

import (
	"github.com/sysu-saad-project/service-end/models/entities"
)

// GetActivityList return wanted activity list with given page number
func GetActivityList(pageNum int) []entities.ActivityInfo {
	activityList := make([]entities.ActivityInfo, 0)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	entities.Engine.Desc("id").Limit(10, pageNum*10).Where("activity.verified = 1").Find(&activityList)
	return activityList
}

// GetActivityInfo return wanted activity detail information which is given by id
func GetActivityInfo(id int) (bool, entities.ActivityInfo) {
	var activity entities.ActivityInfo

	ok, _ := entities.Engine.ID(id).Where("activity.verified = 1").Get(&activity)
	return ok, activity
}

// Check whether user with openId exists
func IsUserExist(openId string) int {
	return 1
}

// Save user with openId in db
func SaveUserInDB(openId string) {
	
}