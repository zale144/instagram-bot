package model

import "time"

type Account struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
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
