package main

import (
	"log"

	micro "github.com/micro/go-micro"
	"github.com/zale144/instagram-bot/htmlToimage/handlers"
	proto "github.com/zale144/instagram-bot/htmlToimage/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("htmltoimage"),
		micro.Version("latest"),
	)
	service.Init()
	proto.RegisterHtmlToImageHandler(service.Server(), new(handlers.HtmlToImage))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
