package handlers

import (
	"context"
	"log"
	proto "github.com/zale144/instagram-bot/services/htmlToimage/proto"
	"github.com/zale144/instagram-bot/services/htmlToimage/service"
)

// HtmlToImage implements the proto HtmlToImage service
type HtmlToImage struct{}

// Process will handle process html to image request
func (h HtmlToImage) Process(ctx context.Context, req *proto.ImageRequest, rsp proto.HtmlToImage_ProcessStream) error {
	res, err := service.GenerateImage(req)
	if err != nil {
		log.Println(err)
		return err
	}
	// return the processed image as a stream
	err = rsp.Send(&proto.ImageResponse{Image: res})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
