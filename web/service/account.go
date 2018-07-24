package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/authcookie"
	"github.com/labstack/echo"
	"github.com/zale144/instagram-bot/api/proto"
	sess "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/zale144/instagram-bot/web/model"
	"github.com/micro/go-micro/client"
)

type AccountService struct{}

// Login handles login requests
func (ar AccountService) Login(c echo.Context) error {

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
	sClient := sess.NewSessionService("session", client.DefaultClient)
	rsp, err := sClient.Get(context.TODO(), &sess.SessionRequest{
		Account:  username,
		Password: password,
	})
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	fmt.Printf("logged in as: %s\n", rsp.Account)

	aClient := api.NewLoginService("api", client.DefaultClient)
	aRsp, err := aClient.Login(context.TODO(), &api.LoginReq{
		Username: username,
	})
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	log.Println(aRsp.Token)
	// get the session cookie
	cookie := &http.Cookie{
		Name:  model.CookieName,
		Value: authcookie.NewSinceNow(username, 24*time.Hour, []byte(model.SECRET)),
		Path:  "/",
	}
	fmt.Printf("got cookie: %v\n", cookie)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": aRsp.Token,
	})
}

// Logout handles logout requests
func (ar AccountService) Logout(c echo.Context) error {
	// expire the cookie
	cookie := &http.Cookie{
		Name:    model.CookieName,
		Expires: time.Now(),
		Path:    "/",
	}
	c.SetCookie(cookie)

	user, err := GetUsernameFromCookie(&c)
	if err == nil {
		clnt := sess.NewSessionService("session", client.DefaultClient)
		rsp, err := clnt.Remove(context.TODO(), &sess.SessionRequest{
			Account: user,
		})
		if err != nil {
			c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
			return err
		}
		fmt.Println(rsp.Account)
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

// GetUsernameFromCookie gets the username from the cookie
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
