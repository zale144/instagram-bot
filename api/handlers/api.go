package handlers

import (
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	micro "github.com/micro/go-micro"
	"github.com/zale144/instagram-bot/api/model"
	proto "github.com/zale144/instagram-bot/api/proto"
)

var Srv micro.Service

type Api struct{}

func (w *Api) Job(ctx context.Context, req *proto.JobReq, rsp *proto.JobResp) error {
	return nil
}

func (w *Api) User(ctx context.Context, req *proto.UserReq, rsp *proto.UserResp) error {
	return nil
}

type LoginService struct{}

func (l *LoginService) Login(ctx context.Context, req *proto.LoginReq, rsp *proto.LoginResp) error {
	claims := &model.JwtCustomClaims{
		req.Username,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(model.SECRET))
	if err != nil {
		log.Println(err)
		return err
	}
	rsp.Token = t
	return nil
}

func RegisterService() {
	Srv = micro.NewService(
		micro.Name("api"),
		micro.Version("latest"),
	)
	Srv.Init()
	proto.RegisterApiHandler(Srv.Server(), new(Api))
	proto.RegisterLoginServiceHandler(Srv.Server(), new(LoginService))

	if err := Srv.Run(); err != nil {
		log.Fatal(err)
	}
}