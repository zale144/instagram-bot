package main

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/zale144/instagram-bot/htmlToimage/handlers"
	proto "github.com/zale144/instagram-bot/htmlToimage/proto"
	"github.com/zale144/instagram-bot/htmlToimage/service"
	"os"
)

func main() {

	service.WebURI = os.Getenv("WEB_LOCAL")

	srv := micro.NewService(
		micro.Name("htmltoimage"),
		micro.Version("latest"),
	)
	srv.Init()
	proto.RegisterHtmlToImageHandler(srv.Server(), new(handlers.HtmlToImage))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
