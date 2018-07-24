package handlers

import (
	"context"
	"log"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/zale144/instagram-bot/services/api/model"
	proto "github.com/zale144/instagram-bot/services/api/proto"
)

type Api struct{}

func (w *Api) Job(ctx context.Context, req *proto.JobReq, rsp *proto.JobResp) error {
	// TODO
	return nil
}

func (w *Api) User(ctx context.Context, req *proto.UserReq, rsp *proto.UserResp) error {
	// TODO
	return nil
}

type LoginService struct{}

// Login handles a login request for the api service
func (l *LoginService) Login(ctx context.Context, req *proto.LoginReq, rsp *proto.LoginResp) error {
	log.Println("GOT API LOGIN REQUEST")
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

/*// RegisterService registers the 'api' microservice
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
}*/
