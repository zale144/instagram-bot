package client

import (
	"github.com/zale144/instagram-bot/services/api/proto"
	"github.com/zale144/instagram-bot/services/web/model"
	"context"
)

type Api struct {}

// Login calls the 'api' microservice and creates a new JWT token for the REST API
func (a Api) Login(username string) (string, error) {
	aClient := api.NewLoginService("api", model.Service.Client())
	aRsp, err := aClient.Login(context.TODO(), &api.LoginReq{
		Username: username,
	})
	if err != nil {
		return "", err
	}
	return aRsp.Token, nil
}