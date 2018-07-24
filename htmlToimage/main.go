package main

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/zale144/instagram-bot/htmlToimage/handlers"
	proto "github.com/zale144/instagram-bot/htmlToimage/proto"
	"github.com/zale144/instagram-bot/htmlToimage/service"
	k8s "github.com/micro/kubernetes/go/micro"
	"os"
)

func main() {

	service.WebURI = os.Getenv("WEB_LOCAL")

	srv := k8s.NewService(
		micro.Name("htmltoimage"),
	)
	srv.Init()
	proto.RegisterHtmlToImageHandler(srv.Server(), new(handlers.HtmlToImage))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
