package main

import (
	"log"

	"github.com/zale144/instagram-bot/sessions/handlers"
	proto "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/sessions/service"

	micro "github.com/micro/go-micro"
	"os"
	"github.com/zale144/instagram-bot/sessions/model"
)

func main() {

	model.RpcURI = os.Getenv("RPC_URI")
	if model.RpcURI == "" {
		model.RpcURI = "http://localhost:4000"
	}
	// start the Sessions cache management
	go service.Sessions()

	srv := micro.NewService(
		micro.Name("session"),
		micro.Version("latest"),
	)
	srv.Init()

	proto.RegisterSessionHandler(srv.Server(), &handlers.Session{})
	proto.RegisterInstaHandler(srv.Server(), &handlers.Insta{})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
