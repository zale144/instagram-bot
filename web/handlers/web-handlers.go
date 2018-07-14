package handlers

import (
	"context"
	"log"

	micro "github.com/micro/go-micro"
	proto "github.com/zale144/instagram-bot/web/proto"
)

var Srv micro.Service

type service struct{}

func (w *service) Job(ctx context.Context, req *proto.JobReq, rsp *proto.JobResp) error {
	return nil
}

func (w *service) User(ctx context.Context, req *proto.UserReq, rsp *proto.UserResp) error {
	return nil
}

func RegisterService() {
	Srv = micro.NewService(
		micro.Name("instagram.bot.web"),
		micro.Version("latest"),
	)
	Srv.Init()
	proto.RegisterWebHandler(Srv.Server(), new(service))

	if err := Srv.Run(); err != nil {
		log.Fatal(err)
	}
}
