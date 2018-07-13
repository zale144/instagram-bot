package main

import (
	"instagram-bot/sessions/handlers"
	proto "instagram-bot/sessions/proto"
	"instagram-bot/sessions/session"
	"log"

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
