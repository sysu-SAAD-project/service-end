package entities

import "time"

// ActivityInfo store activity information
type ActivityInfo struct {
	ID              int    `xorm:"pk autoincr 'id'"`
	Name            string `xorm:"varchar(30) notnull"`
	StartTime       *time.Time
	EndTime         *time.Time
	Campus          int
	Location        string `xorm:"varchar(100)"`
	EnrollCondition string `xorm:"varchar(50)"`
	Sponsor         string `xorm:"varchar(50)"`
	Type            int
	PubStartTime    *time.Time
	PubEndTime      *time.Time
	Detail          string `xorm:"varchar(150)" `
	Reward          string `xorm:"varchar(30)"`
	Introduction    string `xorm:"varchar(50)"`
	Requirement     string `xorm:"varchar(50)"`
	Poster          string `xorm:"varchar(64)"`
	Qrcode          string `xorm:"varchar(64)"`
	Email           string `xorm:"varchar(255)"`
	Verified        int
}

type UserInfo struct {
	UserId   string `xorm:"varchar(64) pk"`
	UserName string `xorm:"varchar(64)"`
	Email    string `xorm:"varchar(100)"`
	Phone    string `xorm:"varchar(20)"`
}

type ActApplyInfo struct {
	Actid    int    `xorm:"int notnull pk"`
	UserId   string `xorm:"varchar(64) notnull pk"`
	UserName string `xorm:"varchar(64)"`
	Email    string `xorm:"varchar(100)"`
	Phone    string `xorm:"varchar(20)"`
	School   string `xorm:"varchar(100)"`
}

// TableName defines table name
func (u ActivityInfo) TableName() string {
	return "activity"
}

func (u UserInfo) TableName() string {
	return "user"
}

func (u ActApplyInfo) TableName() string {
	return "actApply"
}
