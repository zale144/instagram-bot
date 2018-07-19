package handlers

import (
	"context"
	"log"

	"github.com/zale144/instagram-bot/sessions/model"
	proto "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/sessions/service"
)

// Session implements the proto service Session
type Session struct{}

// Get handles the get session request
func (m *Session) Get(ctx context.Context, req *proto.SessionRequest, rsp *proto.SessionResponse) error {
	_, err := service.GetSession(&model.Account{
		Username: req.Account,
		Password: req.Password,
	})
	if err != nil {
		log.Println(err)
		rsp.Error = err.Error()
	}
	return err
}

// Remove handles the request to remove a session from cache
func (m *Session) Remove(ctx context.Context, req *proto.SessionRequest, rsp *proto.SessionResponse) error {
	service.Remove(req.Account)
	rsp.Error = ""
	return nil
}