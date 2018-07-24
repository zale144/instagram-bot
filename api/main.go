package main

import (
	"log"

	"github.com/zale144/instagram-bot/api/model"
	"github.com/zale144/instagram-bot/api/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	k8s "github.com/micro/kubernetes/go/micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	cli "github.com/micro/go-plugins/client/grpc"
	srv "github.com/micro/go-plugins/server/grpc"
	"github.com/zale144/instagram-bot/api/handlers"
	proto "github.com/zale144/instagram-bot/api/proto"
	"github.com/micro/go-micro"
)

/*var (
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
	dbConnString = fmt.Sprintf("postgres://%s:%s@db/%s?sslmode=disable",	dbUser, dbPass, dbName)
	dbInfo = flag.String("db-info", dbConnString, "database connection string")
	pImages   = flag.String("pImages", "files/images/profiles", "path to profile images folder")
)*/

func main() {
	//flag.Parse()

	model.WebURL = os.Getenv("WEB_HOST")

	/*model.DBInfo = *dbInfo
	err := model.InitDB()
	if err != nil {
		log.Fatalf("cannot initialize db: %v", err)
		return
	}*/

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Static("/profile-images", "files/images/profiles")

	api := e.Group("/api")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &model.JwtCustomClaims{},
		SigningKey: []byte(model.SECRET),
	}
	api.Use(middleware.JWTWithConfig(config))

	api.GET("/followed", service.UserService{}.GetAllFollowed)
	api.GET("/process/:user", service.UserService{}.ProcessUser)
	api.GET("/process-by-hashtag/:hashtag", service.UserService{}.ProcessUsersByHashtag)
	api.GET("/processed/:page", service.UserService{}.GetProcessed)
	api.GET("/processed-by-job/:jobID/:page", service.UserService{}.GetProcessedByJob)
	api.GET("/jobs", service.JobService{}.GetJobs)
	api.GET("/search/:user", service.UserService{}.Search)
	api.GET("/follow/:user", service.UserService{}.Follow)

	go regService()

	e.Logger.Fatal(e.Start(":8081"))
}

func regService()  {
	srvc := k8s.NewService(
		micro.Name("api"),
		micro.Version("latest"),
	)

	srvc.Init()
	proto.RegisterLoginServiceHandler(srvc.Server(), &handlers.LoginService{})
	proto.RegisterApiHandler(srvc.Server(), &handlers.Api{})

	client.DefaultClient = cli.NewClient()
	server.DefaultServer = srv.NewServer()

	if err := srvc.Run(); err != nil {
		log.Fatal(err)
	}
}
