package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/zale144/instagram-bot/api/client"
	"github.com/zale144/instagram-bot/api/model"
	"github.com/zale144/instagram-bot/api/storage"
	htmlToImage "github.com/zale144/instagram-bot/htmlToimage/proto"
	sess "github.com/zale144/instagram-bot/sessions/proto"
	"strings"
	"time"
)

type UserService struct{}

// get all followed users
func (ur UserService) GetAllFollowed(c echo.Context) error {
	// get user from token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	username := claims.Name

	// get all followed users
	followedUsers, err := client.Session{}.FollowedUsers(username)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, followedUsers)
}

var formValNames = []string{"title", "width", "height", "crop-h", "crop-w", "hostname"}

// process the profile into a link and send it to the user's Instagram profile
func (ur UserService) ProcessUser(c echo.Context) error {
	// get user from token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	account := claims.Name

	params := map[string]string{
		"account":  account,
		"username": c.Param("user"),
	}
	for n := range formValNames {
		if c.FormValue(formValNames[n]) == "" {
			err := errors.New("missing parameter '" + formValNames[n] + "'")
			c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
			return err
		}
		params[formValNames[n]] = c.FormValue(formValNames[n])
	}
	msg, err := ur.Process(params)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.String(http.StatusOK, msg)
}

// process the all profiles found by specified hashtag, limited by provided parameter
func (ur UserService) ProcessUsersByHashtag(c echo.Context) error {
	// get user from token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	account := claims.Name
	if c.Param("hashtag") == "" {
		err := errors.New("missing parameter 'hashtag'")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	limitStr := c.FormValue("limit")
	if limitStr == "" {
		err := errors.New("missing parameter 'limit'")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println(err)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	params := map[string]string{
		"account":  account,
		"hashtag":  c.Param("hashtag"),
		"username": "",
	}
	for n := range formValNames {
		if c.FormValue(formValNames[n]) == "" {
			err := errors.New("missing parameter '" + formValNames[n] + "'")
			c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
			return err
		}
		params[formValNames[n]] = c.FormValue(formValNames[n])
	}

	ongoingJob, err := storage.JobStorage{}.GetOngoingByHashTag(params["hashtag"])
	if err != nil {
		log.Println(err)
	}
	if ongoingJob != nil {
		err := errors.New("a job with the same hashtag is already running")
		log.Println(err)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	job := model.Job{HashTagName: params["hashtag"]}
	storage.JobStorage{}.Insert(&job)

	go func() {
		users := client.Session{}
		processedUsers := []model.UserDetails{}
		h := s.GetInsta().NewHashtag(params["hashtag"])
		for h.Next() {
			for i := range h.Sections {
				for _, i := range h.Sections[i].LayoutContent.Medias {
					if len(i.Item.Images.Versions) != 0 {
						processedUser, err := storage.ProcessedUserStorage{}.GetByUsername(i.Item.User.Username)
						if err != nil {
							log.Println(err)
							err = nil
						}
						if processedUser != nil {
							err := errors.New("user already processed")
							log.Println(err)
							err = nil
							continue
						}
						account = strings.Split(account, "@")[0]
						if i.Item.User.Username == account {
							err := errors.New("should not process the job issuer")
							log.Println(err)
							err = nil
							continue
						}
						fmt.Printf("processing[ url: %s  - user: %s\n", i.Item.Images.Versions[0].URL, i.Item.User.Username)
						details := model.UserDetails{
							i.Item.User.Username,
							i.Item.User.FullName,
							i.Item.User.Biography,
							i.Item.User.FollowerCount,
							i.Item.User.ProfilePicURL,
							"",
						}
						params["username"] = i.Item.User.Username

						user := model.ProcessedUser{Username: params["username"], Job: job, JobID: job.ID, ProcessedAt: time.Now().Unix()}

						_, err = ur.Process(params)
						if err != nil {
							log.Println(err)
							user.Successful = false
						} else {
							user.Successful = true
						}
						storage.ProcessedUserStorage{}.Insert(user)
						fmt.Printf("processed: %s\n", i.Item.User.Username)
						processedUsers = append(processedUsers, details)
						if len(processedUsers) == limit {
							goto end
						}
					}
				}
			}
		}
	end:
		fmt.Println("DONE")
		err := storage.JobStorage{}.NewJobUpdater(job.ID).FinishedAt(time.Now().Unix()).Update(nil)
		if err != nil {
			log.Println(err)
			err = errors.New("error updating job")
			c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		}
		return
	}()
	return c.String(http.StatusOK, "Job created!")
}

// process the profile into a link and send it to the user's Instagram profile
func (ur UserService) Process(params map[string]string) (string, error) {
	hostname := params["hostname"]
	url := hostname + "/user-info/" + params["account"] + "/" + params["username"]

	fmt.Println("user info at url: ", url)
	fmt.Println(params["username"])
	options := htmlToImage.ImageRequest{
		Input:  url,
		Format: "jpg",
		Name:   params["username"],
	}

	width, err := strconv.Atoi(params["width"])
	if err == nil {
		options.Width = int32(width)
	}
	height, err := strconv.Atoi(params["height"])
	if err == nil {
		options.Height = int32(height)
	}
	crop_h, err := strconv.Atoi(params["crop-h"])
	if err == nil {
		options.CropH = int32(crop_h)
	}
	crop_w, err := strconv.Atoi(params["crop-w"])
	if err == nil {
		options.CropW = int32(crop_w)
	}
	options.CropX = int32((width - crop_w) / 2)
	options.CropY = int32((height - crop_h) / 2)
	imgResp, err := client.HtmlToImage{}.Process(options)
	if err != nil {
		log.Println(err)
		return "", err
	}
	path := "files/images/profiles/" + options.Name + ".jpg"
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = f.Write(imgResp.Image)
	if err != nil {
		log.Println(err)
		return "", err
	}
	message := model.AppURL + "/calling-card/" + params["username"]
	mReq := &sess.MessageRequest{
		Sender:    params["account"],
		Recipient: params["username"],
		Title:     params["title"],
		Text:      message,
	}
	mRsp, err := client.Session{}.Message(mReq)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	fmt.Println(mRsp)

	return message, nil
}

// get processed users
func (ur UserService) GetProcessed(c echo.Context) error {
	page := c.Param("page")
	if page == "" {
		err := errors.New("no page passed")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		log.Println(err)
		return err
	}
	users, err := storage.ProcessedUserStorage{}.GetAll(uint(pageNum))
	if err != nil {
		err := errors.New("cannot get processed users")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// get processed users by job id
func (ur UserService) GetProcessedByJob(c echo.Context) error {
	jobIDStr := c.Param("jobID")
	if jobIDStr == "" {
		err := errors.New("no job id passed")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		log.Println(err)
		return err
	}
	page := c.Param("page")
	if page == "" {
		err := errors.New("no page passed")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		log.Println(err)
		return err
	}
	users, err := storage.ProcessedUserStorage{}.GetByJob(uint(jobID), uint(pageNum))
	if err != nil {
		err := errors.New("cannot get processed users")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, users)
}

/* // search profile by username
func (ur UserService) Search(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	username := claims.Name
	// get the session struct from the cache
	s, err := session.GetSession(&model.Account{Username: username})
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	// if there is no cached session, user should login again
	if s == nil {
		err := errors.New("session is expired")
		AccountService{}.Logout(c)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	profile, err := s.GetUserByName(c.Param("user"))
	if err != nil {
		err := errors.New("cannot get profile")
		AccountService{}.Logout(c)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, profile)
} */

/* // follow profile by username
func (ur UserService) Follow(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	username := claims.Name
	// get the session struct from the cache
	s, err := session.GetSession(&model.Account{Username: username})
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	// if there is no cached session, user should login again
	if s == nil {
		err := errors.New("session is expired")
		AccountService{}.Logout(c)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	profile, err := s.GetUserByName(c.Param("user"))
	if err != nil {
		err := errors.New("cannot get profile")
		AccountService{}.Logout(c)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	profile.Follow()
	return c.JSON(http.StatusOK, profile)
} */
