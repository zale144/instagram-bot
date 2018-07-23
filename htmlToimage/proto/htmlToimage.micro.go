// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: htmlToimage/proto/htmlToimage.proto

/*
Package htmltoimage is a generated protocol buffer package.

It is generated from these files:
	htmlToimage/proto/htmlToimage.proto

It has these top-level messages:
	ImageRequest
	ImageResponse
*/
package htmltoimage

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for HtmlToImage service

type HtmlToImageService interface {
	Process(ctx context.Context, in *ImageRequest, opts ...client.CallOption) (HtmlToImage_ProcessService, error)
}

type htmlToImageService struct {
	c    client.Client
	name string
}

func NewHtmlToImageService(name string, c client.Client) HtmlToImageService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "htmltoimage"
	}
	return &htmlToImageService{
		c:    c,
		name: name,
	}
}

func (c *htmlToImageService) Process(ctx context.Context, in *ImageRequest, opts ...client.CallOption) (HtmlToImage_ProcessService, error) {
	req := c.c.NewRequest(c.name, "HtmlToImage.Process", &ImageRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &htmlToImageServiceProcess{stream}, nil
}

type HtmlToImage_ProcessService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*ImageResponse, error)
}

type htmlToImageServiceProcess struct {
	stream client.Stream
}

func (x *htmlToImageServiceProcess) Close() error {
	return x.stream.Close()
}

func (x *htmlToImageServiceProcess) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *htmlToImageServiceProcess) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *htmlToImageServiceProcess) Recv() (*ImageResponse, error) {
	m := new(ImageResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for HtmlToImage service

type HtmlToImageHandler interface {
	Process(context.Context, *ImageRequest, HtmlToImage_ProcessStream) error
}

func RegisterHtmlToImageHandler(s server.Server, hdlr HtmlToImageHandler, opts ...server.HandlerOption) {
	type htmlToImage interface {
		Process(ctx context.Context, stream server.Stream) error
	}
	type HtmlToImage struct {
		htmlToImage
	}
	h := &htmlToImageHandler{hdlr}
	s.Handle(s.NewHandler(&HtmlToImage{h}, opts...))
}

type htmlToImageHandler struct {
	HtmlToImageHandler
}

func (h *htmlToImageHandler) Process(ctx context.Context, stream server.Stream) error {
	m := new(ImageRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.HtmlToImageHandler.Process(ctx, m, &htmlToImageProcessStream{stream})
}

type HtmlToImage_ProcessStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*ImageResponse) error
}

type htmlToImageProcessStream struct {
	stream server.Stream
}

func (x *htmlToImageProcessStream) Close() error {
	return x.stream.Close()
}

func (x *htmlToImageProcessStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *htmlToImageProcessStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *htmlToImageProcessStream) Send(m *ImageResponse) error {
	return x.stream.Send(m)
}
