package main

import (
	"log"

	"github.com/zale144/instagram-bot/sessions/handlers"
	proto "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/sessions/service"
	k8s "github.com/micro/kubernetes/go/micro"
	"github.com/micro/go-micro"
	"github.com/zale144/instagram-bot/sessions/model"
	"os"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	cli "github.com/micro/go-plugins/client/grpc"
	srv "github.com/micro/go-plugins/server/grpc"
)

func main() {

	model.RpcURI = os.Getenv("RPC_URI")

	// start the Sessions cache management
	go service.Sessions()

	serv := k8s.NewService(
		micro.Name("session"),
		micro.Version("latest"),
	)
	serv.Init()

	proto.RegisterSessionHandler(serv.Server(), &handlers.Session{})
	proto.RegisterInstaHandler(serv.Server(), &handlers.Insta{})

	client.DefaultClient = cli.NewClient()
	server.DefaultServer = srv.NewServer()

	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}

}
