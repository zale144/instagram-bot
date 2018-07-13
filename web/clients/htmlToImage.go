package clients

import (
	"context"

	htmlToImage "github.com/zale144/instagram-bot/htmlToimage/proto"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
)

var (
	HtmlClient  *htmlToImage.HtmlToImageService
	HtmlCtx     *context.Context
	HtmlService *micro.Service
)

func RegisterHtmlToImageClient() {
	clientService := micro.NewService()
	clientService.Init()

	client := htmlToImage.NewHtmlToImageService("instagram-bot.htmlToImage", clientService.Client())

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"Purpose": "hashtag",
	})
	HtmlService, HtmlClient, HtmlCtx = &clientService, &client, &ctx
}
