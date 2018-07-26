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
	"github.com/zale144/instagram-bot/services/api/client"
	"github.com/zale144/instagram-bot/services/api/model"
	"github.com/zale144/instagram-bot/services/api/storage"
	htmlToImage "github.com/zale144/instagram-bot/services/htmlToimage/proto"
	sess "github.com/zale144/instagram-bot/services/sessions/proto"
	"strings"
	"time"
)

type UserService struct{}

var formValNames = []string{"title", "width", "height", "crop-h", "crop-w"}

// GetAllFollowed fetches all followed users
func (ur UserService) GetAllFollowed(c echo.Context) error {
	// get user from token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	username := claims.Name

	// get all followed users
	followedUsers, err := client.Session{}.FollowedUsers(username)
	if err != nil {
		log.Println(err)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, followedUsers)
}

// ProcessUser handles requests to process a single user, by provided username
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

// ProcessUsersByHashtag processes all profiles found by specified hashtag, limited by provided parameter
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
		// no return here
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
		processedUsers := []model.UserDetails{}
		usrReq := &sess.UserReq{
			Account: account,
			Hashtag: params["hashtag"],
			Limit: int64(limit),
		}
		users, err := client.Session{}.UsersByHashtag(usrReq)
		if err != nil {
			log.Println(err)
			return
		}
		for _, u := range users {
			processedUser, err := storage.ProcessedUserStorage{}.GetByUsername(u.Username)
			if err != nil {
				log.Println(err)
				err = nil
				// no return here
			}
			if processedUser != nil {
				log.Println("user already processed")
				continue
			}
			account = strings.Split(account, "@")[0]
			if u.Username == account {
				log.Println("should not process the job issuer")
				continue
			}
			details := model.UserDetails{
				Username: u.Username,
				FullName: u.FullName,
				Description: u.Description,
				FollowerCount: int(u.FollowerCount),
				ProfilePicUrl: u.ProfilePicUrl,
			}
			params["username"] = u.Username

			user := model.ProcessedUser{Username: params["username"], Job: job, JobID: job.ID, ProcessedAt: time.Now().Unix()}

			_, err = ur.Process(params)
			if err != nil {
				log.Println(err)
				err = nil
				// no return here
			}
			user.Successful = err == nil
			err = storage.ProcessedUserStorage{}.Insert(user)
			if err != nil {
				log.Println(err)
				// no return here
			}
			processedUsers = append(processedUsers, details)
			if len(processedUsers) == limit {
				goto end
			}
		}
	end:
		err = storage.JobStorage{}.NewJobUpdater(job.ID).FinishedAt(time.Now().Unix()).Update(nil)
		if err != nil {
			log.Println(err)
			// no return here
		}
		return
	}()
	return c.String(http.StatusOK, "Job created!")
}

// Process processes the profile into a link and send it to the user's Instagram profile
func (ur UserService) Process(params map[string]string) (string, error) {
	url := "/user-info/" + params["account"] + "/" + params["username"]

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
	cropH, err := strconv.Atoi(params["crop-h"])
	if err == nil {
		options.CropH = int32(cropH)
	}
	cropW, err := strconv.Atoi(params["crop-w"])
	if err == nil {
		options.CropW = int32(cropW)
	}
	options.CropX = int32((width - cropW) / 2)
	options.CropY = int32((height - cropH) / 2)
	imgResp, err := client.HtmlToImage{}.Process(options)
	if err != nil {
		log.Println(err)
		return "", err
	}
	path := model.ImagesPath + options.Name + ".jpg"
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
	message := model.WebURL + "/calling-card/" + params["username"]
	mReq := &sess.MessageRequest{
		Sender:    params["account"],
		Recipient: params["username"],
		Title:     params["title"],
		Text:      message,
	}
	_, err = client.Session{}.Message(mReq)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	return message, nil
}

// GetProcessed will get all processed users, paginated
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

// GetProcessedByJob fetches processed users by job id
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

// Search searches an Instagram profile by username
func (ur UserService) Search(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	account := claims.Name
	username := c.Param("user")
	if username == "" {
		err := errors.New("username must be provided")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	profile, err := client.Session{}.UserInfo(account, username)
	if err != nil {
		err := errors.New("cannot get user info")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, profile)
}

// Follow follows an Instagram profile by provided username
func (ur UserService) Follow(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	account := claims.Name
	username := c.Param("user")
	if username == "" {
		err := errors.New("username must be provided")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	profile, err := client.Session{}.Follow(account, username)
	if err != nil {
		err := errors.New("error following user")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, profile)
}
