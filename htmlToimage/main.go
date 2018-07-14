package main

import (
	"log"

	"github.com/zale144/instagram-bot/htmlToimage/handlers"

	proto "github.com/zale144/instagram-bot/htmlToimage/proto"

	micro "github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Name("instagram.bot.htmltoimage"),
		micro.Version("latest"),
	)
	service.Init()
	proto.RegisterHtmlToImageHandler(service.Server(), new(handlers.HtmlToImage))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
