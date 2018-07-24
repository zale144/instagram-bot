package model

import "github.com/micro/go-micro"

const (
	CookieName          = "InstagramBot.Challenge"
	SECRET              = "$P$Bd2WdVjaRR/De58OX2qVu3XA6aiPaf."
	HEADER_AUTH_USER_ID = "Auth-User-Id"
)

var (
	ApiURL string
	Service micro.Service
)
