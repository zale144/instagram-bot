package clients

import (
	"context"

	session "github.com/zale144/instagram-bot/sessions/proto"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
)

var (
	SessClient  *session.SessionService
	SessCtx     *context.Context
	SessService *micro.Service
)

func RegisterSessionClient() {
	clientService := micro.NewService()
	clientService.Init()

	client := session.NewSessionService("instagram-bot.session", clientService.Client())

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"Purpose": "hashtag",
	})
	SessService, SessClient, SessCtx = &clientService, &client, &ctx
}
