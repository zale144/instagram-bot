package main

import (
	"log"

	"github.com/zale144/instagram-bot/sessions/handlers"
	proto "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/sessions/session"

	micro "github.com/micro/go-micro"
)

func main() {
	go session.Sessions()

	service := micro.NewService(
		micro.Name("instagram-bot.session"),
		micro.Version("latest"),
	)
	service.Init()
	proto.RegisterSessionHandler(service.Server(), new(handlers.Session))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
