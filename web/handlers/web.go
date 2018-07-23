package handlers

import (
	"log"

	micro "github.com/micro/go-micro"
	proto "github.com/zale144/instagram-bot/web/proto"
)

var Srv micro.Service

type service struct{}

// RegisterService registers the 'web' microservice
func RegisterService() {
	Srv = micro.NewService(
		micro.Name("web"),
		micro.Version("latest"),
	)
	Srv.Init()
	proto.RegisterWebHandler(Srv.Server(), new(service))

	if err := Srv.Run(); err != nil {
		log.Fatal(err)
	}
}
