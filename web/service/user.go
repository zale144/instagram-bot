package service

import (
	"errors"
	"net/http"
	sess "github.com/zale144/instagram-bot/sessions/proto"
	"github.com/labstack/echo"
	"github.com/zale144/instagram-bot/web/model"
	"context"
	"github.com/micro/go-micro/client"
)

type UserService struct{}

// GetProfile renders the page where basic profile info is shown
// it is consumed by the 'htmlToimage' microservice,
// to generate an image out of the html page
func (ur UserService) GetProfile(c echo.Context) error {
	account := c.Param("account")
	username := c.Param("user")

	sClient := sess.NewInstaService("session", client.DefaultClient)
	rsp, err := sClient.UserInfo(context.Background(), &sess.UserReq{
		Account:  account,
		Username: username,
	})
	details := rsp.User
	if err != nil {
		err := errors.New("cannot get user info")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	data := map[string]interface{}{
		"Username":       details.Username,
		"FullName":       details.FullName,
		"Description":    details.Description,
		"ProfilePicUrl":  details.ProfilePicUrl,
		"FeaturedPicUrl": details.FeaturedPicUrl,
		"FollowerCount":  details.FollowerCount,
	}
	return c.Render(http.StatusOK, "profile", data)
}

// CallingCard will serve the calling card HTML of the Instagram user
func (ur UserService) CallingCard(c echo.Context) error {

	data := map[string]interface{}{
		"Username": c.Param("user"),
		"ApiURL":   model.ApiURL,
	}
	return c.Render(http.StatusOK, "calling-card", data)
}
