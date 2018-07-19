package client

import (
	"context"
	"log"

	"github.com/zale144/instagram-bot/api/handlers"
	sess "github.com/zale144/instagram-bot/sessions/proto"
)

type Session struct{}

func (s Session) FollowedUsers(account string) ([]*sess.User, error) {
	sClient := sess.NewInstaService("session", handlers.Srv.Client())
	rsp, err := sClient.FollowedUsers(context.Background(), &sess.UserReq{
		Account: account,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rsp.Users, nil
}

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

func (s Session) Message(req *sess.MessageRequest) (string, error) {
	sClient := sess.NewInstaService("session", handlers.Srv.Client())
	sRsp, err := sClient.Message(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return sRsp.Response, nil
}

func (s Session) UsersByHashtag(req *sess.UserReq) ([]*sess.User, error) {
	sClient := sess.NewInstaService("session", handlers.Srv.Client())
	sRsp, err := sClient.UsersByHashtag(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sRsp.Users, nil
}