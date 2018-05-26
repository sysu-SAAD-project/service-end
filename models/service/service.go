package service

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"

	"github.com/sysu-saad-project/service-end/models/entities"
)

// GetActivityList return wanted activity list with given page number
func GetActivityList(pageNum int) []entities.ActivityInfo {
	activityList := make([]entities.ActivityInfo, 0)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	entities.Engine.Desc("id").Where("activity.verified = 1").Find(&activityList)
	from := pageNum * 10
	if from >= len(activityList) {
		return []entities.ActivityInfo{}
	}
	if from+10 > len(activityList) {
		return activityList[from:]
	}
	return activityList[from : from+10]
}

// GetActivityListByUserId return wanted activity list with given page number and userOpenId
func GetActivityListByUserId(pageNum int, userOpenId string) []entities.Activity_StudentIdInfo {
	activityList := make([]entities.Activity_StudentIdInfo, 0)
	actApplyList := GetActApplyListByUserId(userOpenId)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	for i := 0; i < len(actApplyList); i++ {
		ok, tmp := GetActivityInfo(actApplyList[i].ActId)
		if !ok {
			continue
		}
		activityList = append(activityList, entities.Activity_StudentIdInfo{
			ID:        tmp.ID,
			Name:      tmp.Name,
			StartTime: tmp.StartTime,
			EndTime:   tmp.EndTime,
			Campus:    tmp.Campus,
			Type:      tmp.Type,
			Poster:    tmp.Poster,
			Location:  tmp.Location,
			StudentId: actApplyList[i].StudentId,
		})
	}
	from := pageNum * 10
	if from >= len(activityList) {
		return []entities.Activity_StudentIdInfo{}
	}
	if from+10 > len(activityList) {
		return activityList[from:]
	}
	return activityList[from : from+10]
}

// GetActivityInfo return wanted activity detail information which is given by id
func GetActivityInfo(id int) (bool, entities.ActivityInfo) {
	var activity entities.ActivityInfo

	ok, _ := entities.Engine.ID(id).Where("activity.verified = 1").Get(&activity)
	return ok, activity
}

// Check whether user with openId exists --- fix user_id into userid
func IsUserExist(openId string) bool {
	has, _ := entities.Engine.Table("user").Where("user_id = ?", openId).Exist(&entities.UserInfo{})
	return has
}

// Check whether activity with actId exists
func IsActExist(actId int) bool {
	has, _ := entities.Engine.Table("activity").Where("id = ?", actId).Exist(&entities.ActivityInfo{})
	return has
}

// Check whether record with actId and userId exists
func IsRecordExist(actId int, studentId string) bool {
	has, _ := entities.Engine.Table("actapply").Where("actid = ? and studentid = ?", actId, studentId).Exist(&entities.ActApplyInfo{})
	return has
}

// IsDiscussionExist check whether discussion exists
func IsDiscussionExist(userOpenId string, mtype int, content string, mTime *time.Time) bool {
	has, _ := entities.Engine.Table("discussion").Where("userid = ? and type = ? and content = ? and time = ?", userOpenId, mtype, content, mTime.Format("2006-01-02 15:04:05")).Exist(&entities.DiscussionInfo{})
	return has
}

// IsPrecusorExist check whether precusor exists
func IsPrecusorExist(precusor int) bool {
	fmt.Println(precusor)
	has, _ := entities.Engine.Table("discussion").Where("disid = ?", precusor).Exist(&entities.DiscussionInfo{})
	return has
}

// IsCommentExist check whether comment exists
func IsCommentExist(userOpenId string, content string, mTime *time.Time, precusor int) bool {
	has, _ := entities.Engine.Table("comment").Where("userid = ? and content = ? and time = ? and precusor = ?", userOpenId, content, mTime.Format("2006-01-02 15:04:05"), precusor).Exist(&entities.CommentInfo{})
	return has
}

// Save user with openId in db
func SaveUserInDB(openId string) {
	user := entities.UserInfo{openId, "", "", ""}
	entities.Engine.InsertOne(&user)
	return
}

// Save actapply with info...(ActApplyInfo) indb
func SaveActApplyInDB(actId int, userId string, userName string, studentId string, phone string, school string) bool {
	actApply := entities.ActApplyInfo{actId, userId, userName, studentId, phone, school}
	_, err := entities.Engine.InsertOne(&actApply)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

// Save discussion with info(DiscussionInfo) in db
func SaveDiscussionInDB(userId string, mtype int, content string, mTime *time.Time) bool {
	discussion := entities.DiscussionInfo{UserId: userId, Type: mtype, Content: content, Time: mTime}
	_, err := entities.Engine.InsertOne(&discussion)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

// Save comment with info(CommentInfo) in db
func SaveCommentInDB(userId string, content string, mTime *time.Time, precusor int) bool {
	comment := entities.CommentInfo{UserId: userId, Content: content, Time: mTime, Precusor: precusor}
	_, err := entities.Engine.InsertOne(&comment)
	if err != nil {
		fmt.Println(err)
	}
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
	entities.Engine.Where("userid = ?", openId).Find(&actApplyList)
	return actApplyList
}

// GetDiscussionList return required discussion list
func GetDiscussionList(pageNum, tp int) []entities.DiscussionInfo {
	discussionList := make([]entities.DiscussionInfo, 0)
	entities.Engine.Where("type=?", tp).Desc("time").Find(&discussionList)
	from := pageNum * 10
	if from >= len(discussionList) {
		return []entities.DiscussionInfo{}
	}
	if from+10 > len(discussionList) {
		return discussionList[from:]
	}
	return discussionList[from : from+10]
}

// GetCommentsList gets needed comments list according to the precusor
func GetCommentsList(pageNum, precusor int) []entities.CommentInfo {
	CommentsInfo := make([]entities.CommentInfo, 0)
	entities.Engine.Where("precusor=?", precusor).Find(&CommentsInfo)
	from := pageNum * 10
	if from >= len(CommentsInfo) {
		return []entities.CommentInfo{}
	}
	if from+10 > len(CommentsInfo) {
		return CommentsInfo[from:]
	}
	return CommentsInfo[from : from+10]
}

// GetDiscussionIterate gets the iterate of all the db
func GetDiscussionIterate() *xorm.Rows {
	discussInfo := new(entities.DiscussionInfo)
	iter, err := entities.Engine.Rows(discussInfo)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return iter
}
