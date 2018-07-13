package resource

import (
	"fmt"
	"instagram-bot/web/clients"
	"instagram-bot/web/model"
	"log"
	"net/http"
	"strings"
	"time"

	session "github.com/zale144/instagram-bot/sessions/proto"

	"github.com/dchest/authcookie"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type AccountResource struct{}

// jwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// method for handling login requests
func (ar AccountResource) Login(c echo.Context) error {

	username, password, ok := c.Request().BasicAuth()

	// No Authentication header
	if ok != true {
		err := fmt.Errorf("bad auth credentials")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	if username == "" || password == "" {
		return echo.ErrUnauthorized
	}
	client := *clients.SessClient
	rsp, err := client.Get(*clients.SessCtx, &session.SessionRequest{
		Account:  username,
		Password: password,
	})
	if err != nil || rsp.Error != "" {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	fmt.Printf("logged in as: %s\n", username)

	claims := &JwtCustomClaims{
		username,
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
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	log.Println(t)
	// get the session cookie
	cookie := &http.Cookie{
		Name:  model.CookieName,
		Value: authcookie.NewSinceNow(username, 24*time.Hour, []byte(model.SECRET)),
		Path:  "/",
	}
	fmt.Printf("got cookie: %v\n", cookie)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// method for handling logout requests
func (ar AccountResource) Logout(c echo.Context) error {
	// expire the cookie
	cookie := &http.Cookie{
		Name:    model.CookieName,
		Expires: time.Now(),
		Path:    "/",
	}
	c.SetCookie(cookie)

	user, err := GetUsernameFromCookie(&c)
	if err == nil {
		client := *clients.SessClient
		rsp, err := client.Remove(*clients.SessCtx, &session.SessionRequest{
			Account: user,
		})
		if err != nil || rsp.Error != "" {
			c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
			return err
		}
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

// get the username from the cookie
func GetUsernameFromCookie(cp *echo.Context) (string, error) {
	c := *cp
	headers := c.Request().Header
	cookieStr := headers.Get("cookie")
	if cookieStr == "" {
		err := fmt.Errorf("empty cookie")
		log.Println(err.Error())
	}
	value := strings.Replace(cookieStr, model.CookieName+"=", "", -1)
	email := authcookie.Login(value, []byte(model.SECRET))
	if email == "" {
		err := fmt.Errorf("no user authenticated")
		log.Println(err.Error())
		return "", err
	}
	return email, nil
}
