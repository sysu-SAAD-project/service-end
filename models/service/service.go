package service

import (
	"fmt"

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

// Check whether user with openId exists --- fix user_id into userid
func IsUserExist(openId string) bool {
	has, _ := entities.Engine.Table("user").Where("userid = ?", openId).Exist()
	return has
}

// Check whether activity with actId exists
func IsActExist(actId int) bool {
	has, _ := entities.Engine.Table("activity").Where("id = ?", actId).Exist()
	return has
}

// Check whether record with actId and userId exists
func IsRecordExist(actId int, studentId string) bool {
	has, _ := entities.Engine.Table("actapply").Where("actid = ? and studentid = ?", actId, studentId).Exist()
	return has
}

// Save user with openId in db
func SaveUserInDB(openId string) {
	user := entities.UserInfo{openId, "", "", ""}
	entities.Engine.Table("user").InsertOne(&user)
	return
}

// Save actapply with info...(ActApplyInfo) indb
func SaveActApplyInDB(actId int, userId string, userName string, studentId string, phone string, school string) bool {
	actApply := entities.ActApplyInfo{actId, userId, userName, studentId, phone, school}
	_, err := entities.Engine.Table("actApply").InsertOne(&actApply)
	return err == nil
}

// CheckUserID check if the user exists in the db --- yubei's part but I change user_id into userid
func CheckUserID(userid string) bool {
	user := new(entities.UserInfo)
	count, err := entities.Engine.Where("userid = ?", userid).Count(user)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count == 1
}

// GetActApplyListByUserId return wanted activity apply list with given user openId
func GetActApplyListByUserId(openId string) []entities.ActApplyInfo {
	actApplyList := make([]entities.ActApplyInfo, 0)
	entities.Engine.Table("actApply").Where("userid = ?", openId).Find(&actApplyList)
	return actApplyList
}
