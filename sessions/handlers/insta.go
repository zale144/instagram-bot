package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/zale144/instagram-bot/sessions/model"
	proto "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/sessions/service"
)

type Insta struct{}

// Message handles the send message request
func (m *Insta) Message(ctx context.Context, req *proto.MessageRequest, rsp *proto.MessageResponse) error {
	fmt.Printf("got: Sender: %s, Title: %s, Ricipient: %s, Text: %s", req.Sender, req.Title, req.Recipient, req.Text)

	s, err := service.GetSession(&model.Account{Username: req.Sender})
	if err != nil {
		log.Println(err)
		return err
	}
	userByName, err := s.GetUserByName(req.Recipient)
	if err != nil {
		log.Println(err)
		return err
	}
	response, err := s.SendDirectMessage(fmt.Sprintf("%v", userByName.ID), req.Text, req.Title)
	if err != nil {
		log.Println(err)
	}
	rsp.Response = response
	return err
}

// FollowedUsers handles the request to get all followed Instagram users
func (m *Insta) FollowedUsers(ctx context.Context, req *proto.UserReq, rsp *proto.Users) error {
	s, err := service.GetSession(&model.Account{
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

// UserInfo handles the request to retrieve basic Instagram user info
func (m *Insta) UserInfo(ctx context.Context, req *proto.UserReq, rsp *proto.UserResp) error {
	s, err := service.GetSession(&model.Account{
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

// UsersByHashtag handles the request to batch process
// all Instagram users associated with the provided hashtag
func (m *Insta) UsersByHashtag(ctx context.Context, req *proto.UserReq, rsp *proto.Users) error {
	s, err := service.GetSession(&model.Account{
		Username: req.Account,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	users := []*proto.User{}
	for _, v := range s.GetUsersByHashtag(req.Hashtag, int(req.Limit)) {
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

// Follow handles the request to follow the Instagram user ith provided username
func (m *Insta) Follow(ctx context.Context, req *proto.UserReq, rsp *proto.UserResp) error {
	s, err := service.GetSession(&model.Account{
		Username: req.Account,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	profile, err := s.Follow(req.Username)
	if err != nil {
		log.Println(err)
		return err
	}
	rsp.User = &profile
	return nil
}
