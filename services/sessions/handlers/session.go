package handlers

import (
	"context"
	"log"
	"github.com/zale144/instagram-bot/services/sessions/model"
	proto "github.com/zale144/instagram-bot/services/sessions/proto"
	"github.com/zale144/instagram-bot/services/sessions/service"
)

// Session implements the proto service Session
type Session struct{}

// Get handles the get session request
func (m *Session) Get(ctx context.Context, req *proto.SessionRequest, rsp *proto.SessionResponse) error {
	s, err := service.GetSession(&model.Account{
		Username: req.Account,
		Password: req.Password,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	rsp.Account = s.GetInsta().Account.Username
	return nil
}

// Remove handles the request to remove a session from cache
func (m *Session) Remove(ctx context.Context, req *proto.SessionRequest, rsp *proto.SessionResponse) error {
	service.Remove(req.Account)
	return nil
}
