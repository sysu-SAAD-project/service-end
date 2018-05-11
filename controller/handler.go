package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	dbservice "github.com/sysu-saad-project/service-end/models/service"
)

// ShowActivitiesListHandler get required page number and return detailed activity list
func ShowActivitiesListHandler(w http.ResponseWriter, r *http.Request) {
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber string
	if len(r.Form["pageNum"]) > 0 {
		pageNumber = r.Form["pageNum"][0]
	} else {
		pageNumber = "1"
	}
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intPageNum > 0 {
		// Get activity list
		activityList := dbservice.GetActivityList(intPageNum - 1)

		// Change each element to the format that we need
		infoArr := make([]ActivityIntroduction, 0)
		for i := 0; i < len(activityList); i++ {
			tmp := ActivityIntroduction{
				ID:        activityList[i].ID,
				Name:      activityList[i].Name,
				StartTime: activityList[i].StartTime.UnixNano() / int64(time.Millisecond),
				EndTime:   activityList[i].EndTime.UnixNano() / int64(time.Millisecond),
				Campus:    activityList[i].Campus,
				Type:      activityList[i].Type,
				Poster:    activityList[i].Poster,
				Location:  activityList[i].Location,
			}
			tmp.Poster = GetPoster(tmp.Poster, tmp.Type)
			infoArr = append(infoArr, tmp)
		}
		returnList := ActivityList{
			Content: infoArr,
		}

		// Transfer it to json
		stringList, err := json.Marshal(returnList)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(500)
			return
		}
		if len(activityList) <= 0 {
			w.WriteHeader(204)
		} else {
			w.Write(stringList)
		}
	} else {
		w.WriteHeader(400)
	}
}

// ShowActivityDetailHandler return required activity details with given activity id
func ShowActivityDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intID > 0 {
		ok, activityInfo := dbservice.GetActivityInfo(intID)
		if ok {
			// Convert to ms
			retMsg := ActivityInfo{
				ID:              activityInfo.ID,
				Name:            activityInfo.Name,
				StartTime:       activityInfo.StartTime.UnixNano() / int64(time.Millisecond),
				EndTime:         activityInfo.EndTime.UnixNano() / int64(time.Millisecond),
				Campus:          activityInfo.Campus,
				Location:        activityInfo.Location,
				EnrollCondition: activityInfo.EnrollCondition,
				Sponsor:         activityInfo.Sponsor,
				Type:            activityInfo.Type,
				PubStartTime:    activityInfo.PubStartTime.UnixNano() / int64(time.Millisecond),
				PubEndTime:      activityInfo.PubEndTime.UnixNano() / int64(time.Millisecond),
				Detail:          activityInfo.Detail,
				Reward:          activityInfo.Reward,
				Introduction:    activityInfo.Introduction,
				Requirement:     activityInfo.Requirement,
				Poster:          activityInfo.Poster,
				Qrcode:          activityInfo.Qrcode,
				Email:           activityInfo.Email,
				Verified:        activityInfo.Verified,
			}
			retMsg.Poster = GetPoster(retMsg.Poster, retMsg.Type)
			stringInfo, err := json.Marshal(retMsg)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				w.WriteHeader(500)
				return
			}
			w.Write(stringInfo)
		} else {
			w.WriteHeader(204)
		}
	} else {
		w.WriteHeader(400)
	}
}

// UserLoginHandler return token string with given user code
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse parameters
	var reqBody map[string]interface{}
	tmpBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(tmpBody, &reqBody)

	var code string = reqBody["code"].(string)
	var token, jwt, openId, tokenOpenId string = "", "", "", ""
	var tokenStatusCode int = -1
	var userStatusCode bool = false
	var err error

	if len(r.Header.Get("Authorization")) > 0 {
		token = r.Header.Get("Authorization")
	}

	// Condition: token exists
	if token != "" {
		// Check token and return status code and params
		// status code: 0 -> check error; 1 -> timeout; 2 -> ok
		tokenStatusCode, tokenOpenId = CheckToken(token)

		// Check whether user exist and return status code
		// status code: false -> not exist; true -> exist
		userStatusCode = dbservice.IsUserExist(tokenOpenId)

		if tokenStatusCode == 2 && userStatusCode == true {
			jwt = token
		} else if tokenStatusCode == 0 {
			// token check error
			w.WriteHeader(401)
			return
		}
	}

	// Condition: token not exists or user not exists while token exists
	// Use HTTP Request get openid from Wechat server
	if token == "" || userStatusCode == false {
		openId, err = GetUserOpenId(code)
		// For test
		// openId, _ = GetUserOpenId(code)
		// openId = "OPENID"
		// For test

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(400)
			return
		}
		// token ok but user not exists, maybe mistake delete
		if openId == tokenOpenId && tokenStatusCode == 2 {
			dbservice.SaveUserInDB(openId)
			jwt = token
		}

		// Check whether user exist, if user don't exist then save user openid in db
		if !dbservice.IsUserExist(openId) {
			dbservice.SaveUserInDB(openId)
		}
	}

	// Condition: token timeout or not exists
	// Generate jwt with openid(sub), issuance time(iat) and expiration time(exp)
	if jwt == "" {
		jwt, err = GenerateJWT(openId)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(400)
			return
		}
	}

	tmpToken := TokenInfo{jwt}
	stringInfo, err := json.Marshal(tmpToken)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	w.Write(stringInfo)
}

// ShowActApplysListHandler parse userOpenId and return detailed activity apply list
func ShowActApplysListHandler(w http.ResponseWriter, r *http.Request) {
	var token, userOpenId string = "", ""
	var tokenStatusCode int = -1
	var userStatusCode bool = false
	var err error

	if len(r.Header.Get("Authorization")) > 0 {
		token = r.Header.Get("Authorization")
	}

	if token == "" {
		// user doesn't login in
		w.WriteHeader(401)
		return
	}

	// Check token and return status code and openId
	// status code: 0 -> check error; 1 -> timeout; 2 -> ok
	tokenStatusCode, userOpenId = CheckToken(token)
	if tokenStatusCode != 2 {
		// user token string error or timeout, need login in again
		w.WriteHeader(401)
		return
	}

	// Check whether user exist and return status code
	// status code: false -> not exist; true -> exist
	userStatusCode = dbservice.IsUserExist(userOpenId)
	if userStatusCode == false {
		// user not exist, need login in again
		w.WriteHeader(401)
		return
	}

	// Get activity apply list
	actApplyList := dbservice.GetActApplyListByUserId(userOpenId)

	// Change each element to the format that we need
	infoArr := make([]ActApplyInfo, 0)
	for i := 0; i < len(actApplyList); i++ {
		tmp := ActApplyInfo{
			Actid:     actApplyList[i].Actid,
			UserName:  actApplyList[i].UserName,
			StudentId: actApplyList[i].StudentId,
			Phone:     actApplyList[i].Phone,
			School:    actApplyList[i].School,
		}
		infoArr = append(infoArr, tmp)
	}
	returnList := ActApplyList{
		Content: infoArr,
	}

	// Transfer it to json
	stringList, err := json.Marshal(returnList)
	if err != nil {
		// error
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(500)
		return
	}
	if len(actApplyList) <= 0 {
		// user doesn't sign up any activity
		w.WriteHeader(204)
	} else {
		// ok
		w.Write(stringList)
	}
}

// UploadActApplyHandler post participant's info and deposite into DB
func UploadActApplyHandler(w http.ResponseWriter, r *http.Request) {
	// Check Authorization validation
	var token, userOpenId string = "", ""
	var tokenStatusCode int = -1
	var userStatusCode bool = false
	var err error

	if len(r.Header.Get("Authorization")) > 0 {
		token = r.Header.Get("Authorization")
	}

	if token == "" {
		// user doesn't login in
		w.WriteHeader(401)
		return
	}

	// Check token and return status code and openId
	// status code: 0 -> check error; 1 -> timeout; 2 -> ok
	tokenStatusCode, userOpenId = CheckToken(token)
	if tokenStatusCode != 2 {
		// user token string error or timeout, need login in again
		w.WriteHeader(401)
		return
	}

	// Check whether user exist and return status code
	// status code: false -> not exist; true -> exist
	userStatusCode = dbservice.IsUserExist(userOpenId)
	if userStatusCode == false {
		// user not exist, need login in again
		w.WriteHeader(401)
		return
	}

	// Parse req form
	r.ParseForm()
	var actId int
	if len(r.Form["actId"]) <= 0 {
		w.WriteHeader(400)
		return
	} else {
		actId, err = strconv.Atoi(r.Form["actId"][0])
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(400)
		}
	}

	// Parse req body
	var reqBody map[string]interface{}
	tmpBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(tmpBody, &reqBody)
	var userName string = reqBody["username"].(string)
	var studentId string = reqBody["studentid"].(string)
	var phone string = reqBody["phone"].(string)
	var school string = reqBody["school"].(string)

	// Check activity exists
	var actExists bool = false
	actExists = dbservice.IsActExist(actId)
	if actExists == false {
		w.WriteHeader(400)
		return
	}

	// Check userId validation
	var studentIdStatus bool = false
	studentIdStatus, _ = regexp.MatchString("^[1-9][0-9]{7}$", studentId)
	if studentIdStatus == false {
		w.WriteHeader(400)
		return
	}

	// Check phone validation
	var phoneStatus bool = false
	phoneStatus, _ = regexp.MatchString(`^(1[3|4|5|8][0-9]\d{8})$`, phone)
	if phoneStatus == false {
		w.WriteHeader(400)
		return
	}

	// Check user repeated registration
	var recordExists bool = false
	recordExists = dbservice.IsRecordExist(actId, studentId)
	if recordExists == true {
		w.WriteHeader(400)
		return
	}

	// Everything is ok
	ok := dbservice.SaveActApplyInDB(actId, userId, userName, studentId, phone, school)
	if !ok {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}
}
