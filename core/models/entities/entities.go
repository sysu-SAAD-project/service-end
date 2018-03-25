package entities

import (
	"time"
)

// ActivityInfo store activity information
type ActivityInfo struct {
	ID              int        `xorm:"pk autoincr" json:"id"`
	Name            string     `xorm:"varchar(30) notnull" json:"name"`
	StartTime       *time.Time `json:"startTime"`
	EndTime         *time.Time `json:"endTime"`
	Campus          int8       `json:"campus"`
	Location        string     `xorm:"varchar(100)" json:"location"`
	EnrollCondition string     `xorm:"varchar(50)" json:"enrollCondition"`
	Sponsor         string     `xorm:"varchar(50)" json:"sponsor"`
	Type            int8       `json:"type"`
	PubStartTime    *time.Time `json:"pubStartTime"`
	PubEndTime      *time.Time `json:"pubEndTime"`
	Detail          string     `xorm:"varchar(150)" json:"detail"`
	Reward          string     `xorm:"varchar(30)" json:"reward"`
	Introduction    string     `xorm:"varchar(50)" json:"introduction"`
	Requirement     string     `xorm:"varchar(50)" json:"requirement"`
	Poster          string     `xorm:"varchar(64)" json:"poster"`
	QRcode          string     `xorm:"varchar(64)"  json:"qrcode"`
	Email           string     `xorm:"varchar(255)"  json:"email"`
	Verified        int8       `json:"verified"`
}
