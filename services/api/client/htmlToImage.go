package client

import (
	"context"
	htmlToImage "github.com/zale144/instagram-bot/services/htmlToimage/proto"
	"github.com/zale144/instagram-bot/services/api/model"
)

type HtmlToImage struct{}

// Process sends a request to 'htmltoimage' microservice to process a single user
func (h HtmlToImage) Process(options htmlToImage.ImageRequest) (*htmlToImage.ImageResponse, error) {
	hClient := htmlToImage.NewHtmlToImageService("htmltoimage", model.Service.Client())
	hResp, err := hClient.Process(context.TODO(), &options)
	if err != nil {
		return nil, err
	}
	return hResp.Recv()
}
