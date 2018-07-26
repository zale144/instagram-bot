package client

import (
	sess "github.com/zale144/instagram-bot/services/sessions/proto"
	"log"
	"context"
	"github.com/zale144/instagram-bot/services/web/model"
)

type Session struct{}

// Get calls the session microservice and fetches a new session
func (s Session) Get(account, password string) (string, error) {
	sClient := sess.NewSessionService("session", model.Service.Client())
	rsp, err := sClient.Get(context.TODO(), &sess.SessionRequest{
		Account:  account,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return rsp.Account, nil
}

// Remove calls the microservice 'session' and removes the user session
func (s Session) Remove(account string) (string, error) {
	sClient := sess.NewSessionService("session", model.Service.Client())
	rsp, err := sClient.Remove(context.TODO(), &sess.SessionRequest{
		Account: account,
	})
	if err != nil {
		return "", err
	}
	return rsp.Account, nil
}

type Insta struct {}

// UserInfo requests user info, from the microservice 'session'
func (i Insta) UserInfo(account, username string) (*sess.User, error) {
	sClient := sess.NewInstaService("session", model.Service.Client())
	rsp, err := sClient.UserInfo(context.Background(), &sess.UserReq{
		Account:  account,
		Username: username,
	})
	if err != nil {
		return nil, err
	}
	return rsp.User, nil
}