package client

import (
	"context"
	"log"

	"github.com/zale144/instagram-bot/api/handlers"
	htmlToImage "github.com/zale144/instagram-bot/htmlToimage/proto"
)

type HtmlToImage struct{}

func (h HtmlToImage) Process(options htmlToImage.ImageRequest) (*htmlToImage.ImageResponse, error) {
	hClient := htmlToImage.NewHtmlToImageService("htmltoimage", handlers.Srv.Client())
	hResp, err := hClient.Process(context.TODO(), &options)
	if err != nil {
		log.Println(err)
		return &htmlToImage.ImageResponse{}, err
	}
	return hResp.Recv()
}
