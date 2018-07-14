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

	srv := micro.NewService(
		micro.Name("instagram.bot.session"),
		micro.Version("latest"),
	)

	srv.Init()

	proto.RegisterSessionHandler(srv.Server(), &handlers.Service{})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
