package client

import (
	"context"
	"log"

	sess "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/web/handlers"
)

type Session struct{}

func (s Session) UserInfo(account, username string) (*sess.User, error) {
	sClient := sess.NewInstaService("session", handlers.Srv.Client())
	rsp, err := sClient.UserInfo(context.Background(), &sess.UserReq{
		Account:  account,
		Username: username,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rsp.User, nil
}
