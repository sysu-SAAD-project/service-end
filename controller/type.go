package controller

// ActivityList defines the return format
type ActivityList struct {
	Content []ActivityIntroduction `json:"content"`
}

type Activity_StudentIdList struct {
	Content []Activity_StudentIdIntroduction `json:"content"`
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

// ActivityIntroducation_StudentIdIntroduction
type Activity_StudentIdIntroduction struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Campus    int    `json:"campus"`
	Type      int    `json:"type"`
	Poster    string `json:"poster"`
	Location  string `json:"location"`
	StudentId string `json:"studentId"`
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
	ActId     int    `json:"actid"`
	UserId    string `json:"userid"`
	UserName  string `json:"username"`
	StudentId string `json:"studentid"`
	Phone     string `json:"phone"`
	School    string `json:"school"`
}

// ActApplyList defines the return format
type ActApplyList struct {
	Content []ActApplyInfo `json:"content"`
}

// DiscussInfo contains info about a discuss
type DiscussInfo struct {
	ID       int    `json:"disid"`
	UserName string `json:"username"`
	Type     int    `json:"type"`
	Content  string `json:"content"`
	Time     int64  `json:"time"`
}

// DiscussList contains list of discussion information
type DiscussList struct {
	Content []DiscussInfo `json:"content"`
}

// CommentInfo contains comment information
type CommentInfo struct {
	ID       int    `json:"cid"`
	UserName string `json:"username"`
	Content  string `json:"content"`
	Time     int64  `json:"time"`
	Precusor int    `json:"precisor"`
}

// CommentList contains list of comments
type CommentList struct {
	Content []CommentInfo `json:"content"`
}

// Error Message Class
type appError struct {
	Code    int 	`json:"codeId"`
	Message string  `json:"message"`
	Error   error   `json:"error"`
}

// Log Error Message Class
type logError struct {
	Api     string 	`json:"api"`
	Message string  `json:"message"`
	Error   error   `json:"error"`
}