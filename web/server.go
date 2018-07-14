package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/zale144/instagram-bot/web/model"
	"github.com/zale144/instagram-bot/web/resource"

	microclient "github.com/micro/go-micro/client"
	sess "github.com/zale144/instagram-bot/sessions/proto"

	"github.com/dchest/authcookie"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	outCh  = make(chan string)
	dbInfo = flag.String("db-info", "root:root@tcp(localhost:3306)/cabani_insta?parseTime=true", "database connection string")
)

func main() {
	model.Port = "4040"
	model.AppURL = "http://localhost:" + model.Port

	flag.Parse()

	model.DBInfo = *dbInfo
	err := model.InitDB()
	if err != nil {
		log.Fatalf("cannot initialize db: %v", err)
		return
	}
	if len(os.Args) > 1 {
		URL := os.Args[1]
		if !strings.Contains(URL, "http") {
			URL = "http://" + os.Args[1]
		}
		model.AppURL = URL
	} else {
		go setTunnelURL()

		out := <-outCh
		model.AppURL = strings.Split(out, ": ")[1]
	}

	fmt.Println("app url: ", model.AppURL)

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/templates/profile/*.html")),
	}

	e.Static("/static", "public/static")

	e.Renderer = t

	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.GET("/login", func(c echo.Context) error {
		return c.File("public/static/html/login.html")
	})

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/admin/home")
	})
	e.POST("/login", resource.AccountResource{}.Login)
	e.GET("/logout", resource.AccountResource{}.Logout)

	e.GET("/user-info/:account/:user", resource.UserResource{}.GetProfile) // for html2image
	e.GET("/calling-card/:user", resource.UserResource{}.CallingCard)      // for sharing

	a := e.Group("/admin")
	a.Use(authMiddleware)

	a.GET("/home", func(c echo.Context) error {
		return c.File("public/static/html/index.html")
	})

	api := e.Group("/api")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &resource.JwtCustomClaims{},
		SigningKey: []byte(model.SECRET),
	}
	api.Use(middleware.JWTWithConfig(config))

	api.GET("/followed", resource.UserResource{}.GetAllFollowed)
	api.GET("/process/:user", resource.UserResource{}.ProcessUser)
	//api.GET("/process-by-hashtag/:hashtag", resource.UserResource{}.ProcessUsersByHashtag)
	api.GET("/processed/:page", resource.UserResource{}.GetProcessed)
	api.GET("/processed-by-job/:jobID/:page", resource.UserResource{}.GetProcessedByJob)
	api.GET("/jobs", resource.JobResource{}.GetJobs)
	//api.GET("/search/:user", resource.UserResource{}.Search)
	//api.GET("/follow/:user", resource.UserResource{}.Follow)

	e.Logger.Fatal(e.Start(":" + model.Port))
}

// check if user is logged in
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(model.CookieName)
		if err == nil {
			login := authcookie.Login(cookie.Value, []byte(model.SECRET))
			if login == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			client := sess.NewSessionService("instagram.bot.session", microclient.DefaultClient)
			rsp, err := client.Get(context.Background(), &sess.SessionRequest{
				Account: login,
			})
			if err != nil || rsp.Error != "" {
				resource.AccountResource{}.Logout(c)
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			c.Request().Header.Set(model.HEADER_AUTH_USER_ID, login)
			return next(c)
		}
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func setTunnelURL() {
	cmd := exec.Command("lt", "--port", model.Port)
	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			outCh <- scanner.Text()
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return
	}
}
