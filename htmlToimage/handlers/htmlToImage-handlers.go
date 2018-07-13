package handlers

import (
	"context"
	"log"

	"github.com/zale144/instagram-bot/htmlToimage/proto"

	proto "github.com/zale144/instagram-bot/htmlToimage/proto"

	htmlToimage "github.com/zale144/instagram-bot/htmlToimage/process"
)

type HtmlToImage struct{}

func (h HtmlToImage) Process(ctx context.Context, req *proto.ImageRequest, rsp proto.HtmlToImage_ProcessStream) error {
	res, err := htmlToimage.GenerateImage(req)
	if err != nil {
		log.Println(err)
		return err
	}
	err = rsp.Send(&htmlToimage_proto.ImageResponse{Image: res})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
