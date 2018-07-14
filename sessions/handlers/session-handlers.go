package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/zale144/instagram-bot/sessions/model"
	proto "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/sessions/session"
)

type Session struct{}

func (m *Session) Message(ctx context.Context, req *proto.MessageRequest, rsp *proto.MessageResponse) error {
	fmt.Printf("got: Sender: %s, Title: %s, Ricipient: %s, Text: %s", req.Sender, req.Title, req.Recipient, req.Text)

	s, err := session.GetSession(&model.Account{Username: req.Sender})
	if err != nil {
		log.Println(err)
		rsp.Error = err.Error()
		return err
	}
	userByName, err := s.GetUserByName(req.Recipient)
	if err != nil {
		log.Println(err)
		rsp.Error = err.Error()
		return err
	}
	response, err := s.SendDirectMessage(fmt.Sprintf("%v", userByName.ID), req.Text, req.Title)
	if err != nil {
		log.Println(err)
		rsp.Error = err.Error()
	}
	rsp.Response = response
	return err
}

func (m *Session) Get(ctx context.Context, req *proto.SessionRequest, rsp *proto.SessionResponse) error {
	_, err := session.GetSession(&model.Account{
		Username: req.Account,
		Password: req.Password,
	})
	if err != nil {
		log.Println(err)
		rsp.Error = err.Error()
	}
	return err
}

func (m *Session) FollowedUsers(ctx context.Context, req *proto.SessionRequest, rsp *proto.Users) error {
	s, err := session.GetSession(&model.Account{
		Username: req.Account,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	users := []*proto.User{}
	for _, v := range s.GetAllFollowedUsers() {
		user := &proto.User{
			Username:       v.Username,
			FullName:       v.FullName,
			Description:    v.Description,
			FollowerCount:  int64(v.FollowerCount),
			ProfilePicUrl:  v.ProfilePicUrl,
			FeaturedPicUrl: v.FeaturedPicUrl,
		}
		users = append(users, user)
	}
	rsp.Users = users
	return err
}

func (m *Session) UserInfo(ctx context.Context, req *proto.UserReq, rsp *proto.UserResp) error {
	s, err := session.GetSession(&model.Account{
		Username: req.Account,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	v, err := s.GetProfileInfo(req.Username)
	user := &proto.User{
		Username:       v.Username,
		FullName:       v.FullName,
		Description:    v.Description,
		FollowerCount:  int64(v.FollowerCount),
		ProfilePicUrl:  v.ProfilePicUrl,
		FeaturedPicUrl: v.FeaturedPicUrl,
	}
	rsp.User = user
	return err
}

func (m *Session) Remove(ctx context.Context, req *proto.SessionRequest, rsp *proto.SessionResponse) error {
	session.Remove(req.Account)
	rsp.Error = ""
	return nil
}
