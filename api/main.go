package main

import (
	"flag"
	"log"

	"github.com/zale144/instagram-bot/api/model"
	"github.com/zale144/instagram-bot/api/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zale144/instagram-bot/api/handlers"
	"os"
)

var (
	dbInfo = flag.String("db-info", "postgres://test:test@db/insta_db?sslmode=disable", "database connection string")
	pImages   = flag.String("pImages", "files/images/profiles", "path to profile images folder")
)

func main() {
	flag.Parse()
	go handlers.RegisterService()

	model.WebURL = os.Getenv("WEB_HOST")

	model.DBInfo = *dbInfo
	err := model.InitDB()
	if err != nil {
		log.Fatalf("cannot initialize db: %v", err)
		return
	}

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Static("/profile-images", *pImages)

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

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "4041"
	}
	e.Logger.Fatal(e.Start(":" + apiPort))
}
