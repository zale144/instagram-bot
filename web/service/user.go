package service

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/zale144/instagram-bot/web/client"
	"github.com/zale144/instagram-bot/web/model"
)

type UserService struct{}

// get profile info
func (ur UserService) GetProfile(c echo.Context) error {
	account := c.Param("account")
	username := c.Param("user")

	details, err := client.Session{}.UserInfo(account, username)
	if err != nil {
		err := errors.New("cannot get user info")
		// TODO use sessions for this
		//AccountService{}.Logout(c)
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
