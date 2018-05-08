package controller

// ActivityList defines the return format
type ActivityList struct {
	Content []ActivityIntroduction `json:"content"`
}

// ErrorMessage defines error format
type ErrorMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"msg"`
}

// ActivityIntroduction include required information in activity list page
type ActivityIntroduction struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Campus    int    `json:"campus"`
	Type      int    `json:"type"`
	Poster    string `json:"poster"`
	Location  string `json:"location"`
}

// ActivityInfo stores json format the front-end wanted
type ActivityInfo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	StartTime       int64  `json:"startTime"`
	EndTime         int64  `json:"endTime"`
	Campus          int    `json:"campus"`
	Location        string `json:"location"`
	EnrollCondition string `json:"enrollCondition"`
	Sponsor         string `json:"sponsor"`
	Type            int    `json:"type"`
	PubStartTime    int64  `json:"pubStartTime"`
	PubEndTime      int64  `json:"pubEndTime"`
	Detail          string `json:"detail"`
	Reward          string `json:"reward"`
	Introduction    string `json:"introduction"`
	Requirement     string `json:"requirement"`
	Poster          string `json:"poster"`
	Qrcode          string `json:"qrcode"`
	Email           string `json:"email"`
	Verified        int    `json:"verified"`
}

// TokenInfo stores json format the front-end wanted
type TokenInfo struct {
	Token string `json:"token"`
}

// ActApplyInfo stores json format the front-end wanted
type ActApplyInfo struct {
	Actid 	 int  `json:"actid"`
	UserName string `json:"username"`
	Email 	 string `json:"email"`
    Phone 	 string `json:"phone"`
    School 	 string `json:"school"`
}

// ActApplyList defines the return format
type ActApplyList struct {
	Content []ActApplyInfo `json:"content"`
}