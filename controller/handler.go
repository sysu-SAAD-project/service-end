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

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sysu-saad-project/service-end/models/entities"
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

// ShowActApplysListHandler parse userOpenId and return activityList for specified user
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
		fmt.Println("Token is empty")
		w.WriteHeader(401)
		return
	}

	// Check token and return status code and openId
	// status code: 0 -> check error; 1 -> timeout; 2 -> ok
	tokenStatusCode, userOpenId = CheckToken(token)
	if tokenStatusCode != 2 {
		// user token string error or timeout, need login in again
		fmt.Println("Token Error or Timeout")
		w.WriteHeader(401)
		return
	}

	// Check whether user exist and return status code
	// status code: false -> not exist; true -> exist
	userStatusCode = dbservice.IsUserExist(userOpenId)
	if userStatusCode == false {
		// user not exist, need login in again
		fmt.Println("Please Login Again")
		w.WriteHeader(401)
		return
	}

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
		activityList := dbservice.GetActivityListByUserId(intPageNum-1, userOpenId)

		// Change each element to the format that we need
		infoArr := make([]Activity_StudentIdIntroduction, 0)
		for i := 0; i < len(activityList); i++ {
			tmp := Activity_StudentIdIntroduction{
				ID:        activityList[i].ID,
				Name:      activityList[i].Name,
				StartTime: activityList[i].StartTime.UnixNano() / int64(time.Millisecond),
				EndTime:   activityList[i].EndTime.UnixNano() / int64(time.Millisecond),
				Campus:    activityList[i].Campus,
				Type:      activityList[i].Type,
				Poster:    activityList[i].Poster,
				Location:  activityList[i].Location,
				StudentId: activityList[i].StudentId,
			}
			tmp.Poster = GetPoster(tmp.Poster, tmp.Type)
			infoArr = append(infoArr, tmp)
		}
		returnList := Activity_StudentIdList{
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
	sactId := mux.Vars(r)["actId"]
	var actId int
	if len(sactId) <= 0 {
		w.WriteHeader(400)
		return
	} else {
		actId, err = strconv.Atoi(sactId)
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

	// Check studentId validation
	var studentIdStatus bool = false
	studentIdStatus, _ = regexp.MatchString("^[1-9][0-9]{7}$", studentId)
	if studentIdStatus == false {
		w.WriteHeader(400)
		return
	}

	// Check phone validation
	var phoneStatus bool = false
	phoneStatus, _ = regexp.MatchString(`^(1[3|4|5|7|8][0-9]\d{8})$`, phone)
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
	ok := dbservice.SaveActApplyInDB(actId, userOpenId, userName, studentId, phone, school)
	if !ok {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}
}

// TokenHandler generate one effective for 300 days token
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	// expire in two weeks
	var exp = time.Hour * 24 * 300
	var hmacSampleSecret = []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "oXRoe0c7KDoAVGKOTYks_kaV2iQA",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write([]byte(tokenString))
	w.WriteHeader(200)
}

// UploadDiscussionHandler post discussion and deposite into DB
func UploadDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	// Check Authorization validation
	var token, userOpenId string = "", ""
	var tokenStatusCode int = -1
	var userStatusCode bool = false
	// var err error

	if len(r.Header.Get("Authorization")) > 0 {
		token = r.Header.Get("Authorization")
	}

	if token == "" {
		// user doesn't login in
		w.WriteHeader(401)
		fmt.Println("User does not login in")
		return
	}

	// Check token and return status code and openId
	// status code: 0 -> check error; 1 -> timeout; 2 -> ok
	tokenStatusCode, userOpenId = CheckToken(token)
	if tokenStatusCode != 2 {
		// user token string error or timeout, need login in again
		w.WriteHeader(401)
		fmt.Println("Need login in again")
		return
	}

	// Check whether user exist and return status code
	// status code: false -> not exist; true -> exist
	userStatusCode = dbservice.IsUserExist(userOpenId)
	if userStatusCode == false {
		// user not exist, need login in again
		w.WriteHeader(401)
		fmt.Println("User not exist, need login in again")
		return
	}

	// Parse req body
	var reqBody map[string]interface{}
	tmpBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(tmpBody, &reqBody)
	var mtype int = int(reqBody["type"].(float64))
	var content string = reqBody["content"].(string)

	// check form
	var typeStatus bool = false
	if mtype == 2 || mtype == 4 || mtype == 8 ||
		mtype == 6 || mtype == 10 || mtype == 12 {
		typeStatus = true
	}
	if typeStatus == false {
		w.WriteHeader(400)
		fmt.Println("typeStatus is false")
		return
	}

	var contentStatus bool = false
	if len(content) < 240 && len(content) > 0 {
		contentStatus = true
	}
	if contentStatus == false {
		fmt.Println("contentStatus is false")
		w.WriteHeader(400)
		return
	}

	currentTime := time.Now()
	discussionExist := dbservice.IsDiscussionExist(userOpenId, mtype, content, &currentTime)
	if discussionExist == true {
		w.WriteHeader(400)
		fmt.Println("discussionExist")
		return
	}

	// Everyting is ok
	ok := dbservice.SaveDiscussionInDB(userOpenId, mtype, content, &currentTime)
	if !ok {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}
}

// UploadCommentHandler post discussion and deposite into DB
func UploadCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Check Authorization validation
	var token, userOpenId string = "", ""
	var tokenStatusCode int = -1
	var userStatusCode bool = false
	// var err error

	if len(r.Header.Get("Authorization")) > 0 {
		token = r.Header.Get("Authorization")
	}

	if token == "" {
		// user doesn't login in
		w.WriteHeader(401)
		fmt.Println("User does not login in")
		return
	}

	// Check token and return status code and openId
	// status code: 0 -> check error; 1 -> timeout; 2 -> ok
	tokenStatusCode, userOpenId = CheckToken(token)
	if tokenStatusCode != 2 {
		// user token string error or timeout, need login in again
		w.WriteHeader(401)
		fmt.Println("Need login in again")
		return
	}

	// Check whether user exist and return status code
	// status code: false -> not exist; true -> exist
	userStatusCode = dbservice.IsUserExist(userOpenId)
	if userStatusCode == false {
		// user not exist, need login in again
		w.WriteHeader(401)
		fmt.Println("User not exist, need login in again")
		return
	}

	// Parse req body
	var reqBody map[string]interface{}
	tmpBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(tmpBody, &reqBody)
	var content string = reqBody["content"].(string)
	var precusor int = int(reqBody["precusor"].(float64))

	var contentStatus bool = false
	if len(content) < 240 && len(content) > 0 {
		contentStatus = true
	}
	if contentStatus == false {
		fmt.Println("contentStatus is false")
		w.WriteHeader(400)
		return
	}

	currentTime := time.Now()

	precusorExist := dbservice.IsPrecusorExist(precusor)
	if precusorExist == false {
		w.WriteHeader(400)
		fmt.Println("precusor do not exist")
		return
	}

	commentExist := dbservice.IsCommentExist(userOpenId, content, &currentTime, precusor)
	if commentExist == true {
		w.WriteHeader(400)
		fmt.Println("commentExist")
		return
	}

	// Everyting is ok
	ok := dbservice.SaveCommentInDB(userOpenId, content, &currentTime, precusor)
	if !ok {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}
}

func ListDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber, disType string
	if len(r.Form["type"]) <= 0 {
		w.WriteHeader(400)
		return
	}
	disType = r.Form["type"][0]
	if len(r.Form["page"]) > 0 {
		pageNumber = r.Form["page"][0]
	} else {
		pageNumber = "1"
	}
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	intType, err := strconv.Atoi(disType)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intPageNum > 0 && intType >= 2 {
		// Judge which type is required
		typeChoosed := []bool{false, false, false}
		typeChoosed[0] = (intType>>3)&1 == 1
		typeChoosed[1] = (intType>>2)&1 == 1
		typeChoosed[2] = (intType>>1)&1 == 1
		// Get required activity
		iterate := dbservice.GetDiscussionIterate()
		if iterate == nil {
			w.WriteHeader(500)
			return
		}
		defer iterate.Close()
		discus := new(entities.DiscussionInfo)
		discussList := make([]entities.DiscussionInfo, 0)
		// Record current number
		cnt := 0
		// Judge every item
		for iterate.Next() && len(discussList) < 10 {
			cnt++
			if cnt < (intPageNum-1)*10 {
				continue
			}
			err := iterate.Scan(discus)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			var i uint
			for i = 0; i < 3; i++ {
				if typeChoosed[i] && typeChoosed[i] == ((discus.Type>>(3-i))&1 == 1) {
					break
				}
			}
			if i < 3 {
				discussList = append(discussList, *discus)
			}
		}
		// Return value
		if len(discussList) <= 0 {
			w.WriteHeader(204)
			return
		}
		content := make([]DiscussInfo, 0)
		for _, v := range discussList {
			tmp := DiscussInfo{v.DisId, v.UserId, v.Type, v.Content, v.Time.UnixNano() / int64(time.Millisecond)}
			content = append(content, tmp)
		}
		ret, err := json.Marshal(DiscussList{content})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Write(ret)
	}
}

func ListCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber, precusor string
	if len(r.Form["precusor"]) <= 0 {
		w.WriteHeader(400)
		return
	}
	precusor = r.Form["precusor"][0]
	if len(r.Form["page"]) > 0 {
		pageNumber = r.Form["page"][0]
	} else {
		pageNumber = "1"
	}
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	intPrecusor, err := strconv.Atoi(precusor)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	if intPageNum > 0 && intPrecusor > 0 {
		commentList := dbservice.GetCommentsList(intPageNum-1, intPrecusor)
		if len(commentList) <= 0 {
			w.WriteHeader(204)
			return
		}
		content := make([]CommentInfo, 0)
		for _, v := range commentList {
			tmp := CommentInfo{v.Cid, v.UserId, v.Content, v.Time.UnixNano() / int64(time.Millisecond), v.Precusor}
			content = append(content, tmp)
		}
		ret, err := json.Marshal(CommentList{content})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Write(ret)
	}
}
