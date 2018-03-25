package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	dbservice "github.com/sysu-saad-project/service-end/core/models/service"
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
		panic(err)
	}

	// Get activity list and transfer it to json
	activityList := dbservice.GetActivityList(intPageNum - 1)
	returnList := ActivityList{
		Content:          activityList,
		Last:             true,
		TotalPages:       1,
		TotalElements:    len(activityList),
		NumberOfElements: len(activityList),
		First:            intPageNum == 1,
		Sort: []SortRule{SortRule{
			Direction:    "DESC",
			Property:     "id",
			IgnoreCase:   false,
			NullHandling: "NATIVE",
			Ascending:    false,
			Descending:   true,
		}},
	}
	stringList, err := json.Marshal(returnList)
	if err != nil {
		panic(err)
	}
	w.Write(stringList)
}

// ShowActivityDetailHandler return required activity details with given activity id
func ShowActivityDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	activityInfo := dbservice.GetActivityInfo(intID)
	stringInfo, err := json.Marshal(activityInfo)
	if err != nil {
		panic(err)
	}
	w.Write(stringInfo)
}
