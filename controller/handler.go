package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	r.ParseForm()
	var code string = r.PostForm.Get("code")
	var token string = ""
	if len(r.Header.Get("Authorization")) > 0 {
		token = r.Header.Get("Authorization")
	}

	// If token exists, which means user openid exists and no need to request from Wechat
	if token != "" {
		// Check token and return status code and params
		// status code: 0 -> check error; 1 -> timeout; 2 -> ok
		[tokenStatusCode, tokenOpenId] = checkToken(token)
		
		// Check whether user exist and return status code
		// status code: 0 -> error; 1 -> ok
		userStatusCode = isUserExist(tokenOpenId)

		if checkStatusCode == 2 && userStatusCode == 1 {
			w.WriteHeader(500)
			w.Write(token)
			return
		} else if tokenStatusCode == 0 {
			// token error
			w.WriteHeader(401)
			return
		}
	}

	// Condition: If token not exists or user not exists while token exists
	// Use HTTP Request get openid from Wechat server
	if token == "" || userStatusCode == 0 {
		[sessionKey, openId] = getUserOpenId(code)
		// token ok but user not exists, maybe mistake delete
		if openId == tokenOpenId && tokenStatusCode == 2 {
			saveUserInDB(openId)
			w.WriteHeader(500)
			w.Write(token)
			return
		}

		// Check whether user exist, if user don't exist then save user openid in db
		if !isUserExist(openId) {
			saveUserInDB(openId)
		}
	}

	// When go to this step, only one condition: token timeout
	// Generate jwt with openid(sub), issuance time(iat) and expiration time(exp)
	jwt, err := generateJWT(openId)

	// Return jwt string
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(500)
	w.Write(jwt)
}
