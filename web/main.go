package main

import (
	"bufio"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"github.com/dchest/authcookie"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	sess "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/web/handlers"
	"github.com/zale144/instagram-bot/web/model"
	"github.com/zale144/instagram-bot/web/service"
	"strings"
)

var outCh = make(chan string)

func main() {
	go handlers.RegisterService()

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = "4040"
	}
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "4041"
	}
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		go runCommand("hostname", "-i")
		hostname = <-outCh
	}
	if !strings.Contains(hostname, "http") {
		hostname = "http://" + hostname
	}
	model.WebURL = hostname + ":" + webPort
	model.ApiURL = hostname + ":" + apiPort
	fmt.Println("WebURL: ", model.WebURL)
	fmt.Println("APIURL", model.ApiURL)

	e := echo.New()

	t := &wTemplate{
		templates: template.Must(template.ParseGlob("public/templates/*.html")),
	}
	e.Static("/static", "public/static")
	e.Renderer = t

	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// ***************** public ***************************
	e.GET("/login", func(c echo.Context) error {
		return c.File("public/static/html/login.html")
	})
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/admin/home")
	})
	e.POST("/login", service.AccountService{}.Login)
	e.GET("/logout", service.AccountService{}.Logout)
	e.GET("/calling-card/:user", service.UserService{}.CallingCard)      // for sharing
	e.GET("/user-info/:account/:user", service.UserService{}.GetProfile) // for html2image

	// ***************** private ***************************
	a := e.Group("/admin")
	a.Use(authMiddleware)
	a.GET("/home", func(c echo.Context) error {
		data := map[string]interface{}{
			"WebURL":   model.WebURL,
			"ApiURL":   model.ApiURL,
		}
		return c.Render(http.StatusOK, "home", data)
	})
	e.Logger.Fatal(e.Start(":" + webPort))
}

// authMiddleware is used to check if user is logged in
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(model.CookieName)
		if err == nil {
			login := authcookie.Login(cookie.Value, []byte(model.SECRET))
			if login == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			client := sess.NewSessionService("session", handlers.Srv.Client())
			rsp, err := client.Get(context.Background(), &sess.SessionRequest{
				Account: login,
			})
			if err != nil || rsp.Error != "" {
				service.AccountService{}.Logout(c)
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			c.Request().Header.Set(model.HEADER_AUTH_USER_ID, login)
			return next(c)
		}
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

type wTemplate struct {
	templates *template.Template
}

func (t *wTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// runCommand will get the container hostname
func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
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
