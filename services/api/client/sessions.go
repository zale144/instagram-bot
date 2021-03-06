package client

import (
	"context"
	sess "github.com/zale144/instagram-bot/services/sessions/proto"
	"github.com/zale144/instagram-bot/services/api/model"
)

type Session struct{}

// FollowedUsers requests all followed users, from the microservice 'session'
func (s Session) FollowedUsers(account string) ([]*sess.User, error) {
	sClient := sess.NewInstaService("session", model.Service.Client())
	rsp, err := sClient.FollowedUsers(context.Background(), &sess.UserReq{
		Account: account,
	})
	if err != nil {
		return nil, err
	}
	return rsp.Users, nil
}

// UserInfo requests user info, from the microservice 'session'
func (s Session) UserInfo(account, username string) (*sess.User, error) {
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

// Message sends a message request to the microservice 'session'
func (s Session) Message(req *sess.MessageRequest) (string, error) {
	sClient := sess.NewInstaService("session", model.Service.Client())
	sRsp, err := sClient.Message(context.TODO(), req)
	if err != nil {
		return "", err
	}
	return sRsp.Response, nil
}

// UsersByHashtag sends a request to the microservice 'session'
// to process all users associated with the provided hashtag
func (s Session) UsersByHashtag(req *sess.UserReq) ([]*sess.User, error) {
	sClient := sess.NewInstaService("session", model.Service.Client())
	sRsp, err := sClient.UsersByHashtag(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	return sRsp.Users, nil
}

// Follow sends a request to the microservice 'session', to follow the specific user
func (s Session) Follow(account, username string) (*sess.User, error) {
	sClient := sess.NewInstaService("session", model.Service.Client())
	sRsp, err := sClient.Follow(context.TODO(), &sess.UserReq{
		Account:  account,
		Username: username,
	})
	if err != nil {
		return nil, err
	}
	return sRsp.User, nil
}