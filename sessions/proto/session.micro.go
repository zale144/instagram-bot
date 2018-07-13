// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: session.proto

/*
Package sessions_proto is a generated protocol buffer package.

It is generated from these files:
	session.proto

It has these top-level messages:
	SessionRequest
	SessionResponse
	Users
	UserReq
	UserResp
	User
	MessageRequest
	MessageResponse
*/
package sessions_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Session service

type SessionService interface {
	Get(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*SessionResponse, error)
	Remove(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*SessionResponse, error)
	Message(ctx context.Context, in *MessageRequest, opts ...client.CallOption) (*MessageResponse, error)
	FollowedUsers(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*Users, error)
	UserInfo(ctx context.Context, in *UserReq, opts ...client.CallOption) (*UserResp, error)
}

type sessionService struct {
	c    client.Client
	name string
}

func NewSessionService(name string, c client.Client) SessionService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "sessions.proto"
	}
	return &sessionService{
		c:    c,
		name: name,
	}
}

func (c *sessionService) Get(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*SessionResponse, error) {
	req := c.c.NewRequest(c.name, "Session.Get", in)
	out := new(SessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionService) Remove(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*SessionResponse, error) {
	req := c.c.NewRequest(c.name, "Session.Remove", in)
	out := new(SessionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionService) Message(ctx context.Context, in *MessageRequest, opts ...client.CallOption) (*MessageResponse, error) {
	req := c.c.NewRequest(c.name, "Session.Message", in)
	out := new(MessageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionService) FollowedUsers(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*Users, error) {
	req := c.c.NewRequest(c.name, "Session.FollowedUsers", in)
	out := new(Users)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionService) UserInfo(ctx context.Context, in *UserReq, opts ...client.CallOption) (*UserResp, error) {
	req := c.c.NewRequest(c.name, "Session.UserInfo", in)
	out := new(UserResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Session service

type SessionHandler interface {
	Get(context.Context, *SessionRequest, *SessionResponse) error
	Remove(context.Context, *SessionRequest, *SessionResponse) error
	Message(context.Context, *MessageRequest, *MessageResponse) error
	FollowedUsers(context.Context, *SessionRequest, *Users) error
	UserInfo(context.Context, *UserReq, *UserResp) error
}

func RegisterSessionHandler(s server.Server, hdlr SessionHandler, opts ...server.HandlerOption) {
	type session interface {
		Get(ctx context.Context, in *SessionRequest, out *SessionResponse) error
		Remove(ctx context.Context, in *SessionRequest, out *SessionResponse) error
		Message(ctx context.Context, in *MessageRequest, out *MessageResponse) error
		FollowedUsers(ctx context.Context, in *SessionRequest, out *Users) error
		UserInfo(ctx context.Context, in *UserReq, out *UserResp) error
	}
	type Session struct {
		session
	}
	h := &sessionHandler{hdlr}
	s.Handle(s.NewHandler(&Session{h}, opts...))
}

type sessionHandler struct {
	SessionHandler
}

func (h *sessionHandler) Get(ctx context.Context, in *SessionRequest, out *SessionResponse) error {
	return h.SessionHandler.Get(ctx, in, out)
}

func (h *sessionHandler) Remove(ctx context.Context, in *SessionRequest, out *SessionResponse) error {
	return h.SessionHandler.Remove(ctx, in, out)
}

func (h *sessionHandler) Message(ctx context.Context, in *MessageRequest, out *MessageResponse) error {
	return h.SessionHandler.Message(ctx, in, out)
}

func (h *sessionHandler) FollowedUsers(ctx context.Context, in *SessionRequest, out *Users) error {
	return h.SessionHandler.FollowedUsers(ctx, in, out)
}

func (h *sessionHandler) UserInfo(ctx context.Context, in *UserReq, out *UserResp) error {
	return h.SessionHandler.UserInfo(ctx, in, out)
}
