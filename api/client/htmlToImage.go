package client

import (
	"context"
	"log"

	htmlToImage "github.com/zale144/instagram-bot/htmlToimage/proto"
	"github.com/micro/go-micro/client"
)

type HtmlToImage struct{}

// Process sends a request to 'htmltoimage' microservice to process a single user
func (h HtmlToImage) Process(options htmlToImage.ImageRequest) (*htmlToImage.ImageResponse, error) {
	hClient := htmlToImage.NewHtmlToImageService("htmltoimage", client.DefaultClient)
	hResp, err := hClient.Process(context.TODO(), &options)
	if err != nil {
		log.Println(err)
		return &htmlToImage.ImageResponse{}, err
	}
	return hResp.Recv()
}
