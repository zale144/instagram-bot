package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/dchest/authcookie"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zale144/instagram-bot/services/web/model"
	"fmt"
	k8s "github.com/micro/kubernetes/go/micro"
	"github.com/zale144/instagram-bot/services/web/service"
	"github.com/micro/go-micro"
	"github.com/zale144/instagram-bot/services/web/client"
)

func main() {

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

	go reqService()

	e.Logger.Fatal(e.Start(":8081"))
}

// reqService registers the 'web' microservice
func reqService()  {
	model.Service = k8s.NewService(
		micro.Name("web"),
		micro.Version("latest"),
	)
	model.Service.Init()

	if err := model.Service.Run(); err != nil {
		log.Fatal(err)
	}
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
			acc, err := client.Session{}.Get(login, "")
			if err != nil {
				service.AccountService{}.Logout(c)
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			fmt.Println(acc)
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

