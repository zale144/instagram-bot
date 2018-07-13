package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	CookieName          = "CabaniInstagram.Challenge"
	SECRET              = "$P$Bd2WdVjaRR/De58OX2qVu3XA6aiPaf."
	HEADER_AUTH_USER_ID = "Auth-User-Id"
)

var (
	Port   string
	AppURL string
)

type Account struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type UserBrief struct {
	ID            int64
	Name          string
	ProfilePicUrl string
}

type UserDetails struct {
	Username       string
	FullName       string
	Description    string
	FollowerCount  int
	ProfilePicUrl  string
	FeaturedPicUrl string
}

type ProcessedUser struct {
	gorm.Model
	Username    string `sql:"unique"`
	Job         Job    `gorm:"ForeignKey:ID;AssociationForeignKey:JobID"`
	JobID       uint
	ProcessedAt int64
	Successful  bool
}

type Job struct {
	gorm.Model
	HashTagName string
	FinishedAt  int64
}

type Media struct {
	ID          string
	URL         string
	IsLandscape bool
	IsPicOfUser bool
	UserID      int64
	Username    string
	LikeCount   int
	PostedAt    time.Time
}

func (m *Media) String() string {
	return fmt.Sprintf("Media: [%d likes] @%s %s",
		m.LikeCount, m.Username, m.URL,
	)
}
