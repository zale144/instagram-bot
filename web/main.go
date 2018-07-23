package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/dchest/authcookie"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	sess "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/web/handlers"
	"github.com/zale144/instagram-bot/web/model"
	"github.com/zale144/instagram-bot/web/service"
)

func main() {
	go handlers.RegisterService()

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
	e.GET("/user-info/:account/:user", service.UserService{}.GetProfile) // for htmlToimage

	// ***************** private ***************************
	a := e.Group("/admin")
	a.Use(authMiddleware)

	model.ApiURL = os.Getenv("API_HOST")

	a.GET("/home", func(c echo.Context) error {
		data := map[string]interface{}{
			"ApiURL":   model.ApiURL,
		}
		return c.Render(http.StatusOK, "home", data)
	})
	webPort := os.Getenv("WEB_PORT")
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

