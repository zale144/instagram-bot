package session

import (
	"errors"
	"fmt"
	"instagram-bot/sessions/model"
	proto "instagram-bot/sessions/proto"
	"log"
	"net/url"
	"sort"
	"time"

	"github.com/zale144/goinsta"
)

type Session struct {
	insta *goinsta.Instagram
}

// create a new sessions
func NewSession(account *model.Account) (*Session, error) {
	s := &Session{
		insta: goinsta.New(account.Username, account.Password),
	}
	err := s.insta.Login()
	if err != nil || s.insta.Account == nil {
		err = errors.New("Bad credentials or permission needed from Instagram")
		log.Println(err)
		return nil, err
	}
	err = s.insta.Export(account.Username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// and import again
	s.insta, err = goinsta.Import(account.Username)
	if s.insta.Account == nil {
		err := errors.New(fmt.Sprintf("cannot import goinsta config with name: %s", account.Username))
		log.Println(err)
		return nil, err
	}
	// save to cache
	SaveSession(s, account.Username)
	return s, nil
}

// get all followed users
func (s *Session) GetAllFollowedUsers() []proto.User {
	account := s.insta.Account
	followingUsers := []proto.User{}
	following := account.Following()

	for following.Next() {
		for i := range following.Users {
			details := proto.User{
				Username:      following.Users[i].Username,
				FullName:      following.Users[i].FullName,
				Description:   following.Users[i].Biography,
				FollowerCount: int64(following.Users[i].FollowerCount),
				ProfilePicUrl: following.Users[i].ProfilePicURL,
			}
			followingUsers = append(followingUsers, details)
		}
	}
	return followingUsers
}

// get users by hashtag
func (s *Session) GetUsersByHashtag(hashtag string, limit int) []proto.User {
	users := []proto.User{}

	h := s.insta.NewHashtag(hashtag)
	for h.Next() {
		for i := range h.Sections {
			for _, i := range h.Sections[i].LayoutContent.Medias {
				if len(i.Item.Images.Versions) != 0 {
					fmt.Printf("url: %s  - user: %s\n", i.Item.Images.Versions[0].URL, i.Item.User.Username)
					if Contains(users, i.Item.User.Username) { // or in database
						continue
					}
					details := proto.User{
						Username:      i.Item.User.Username,
						FullName:      i.Item.User.FullName,
						Description:   i.Item.User.Biography,
						FollowerCount: int64(i.Item.User.FollowerCount),
						ProfilePicUrl: i.Item.User.ProfilePicURL,
					}
					users = append(users, details)
				}
			}
		}
		if len(users) == limit {
			break
		}
	}
	return users
}

// Checks if the user is in the slice
func Contains(users []proto.User, username string) bool {
	for _, currentUser := range users {
		if currentUser.Username == username {
			return true
		}
	}
	return false
}

// get the user struct from instagram API
func (s *Session) GetUserByName(name string) (*goinsta.User, error) {
	user, err := s.insta.Profiles.ByName(name)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// get the largest picture from an image item
func getLargestCandidate(candidates []goinsta.Candidate) goinsta.Candidate {
	m := candidates[0]
	for _, c := range candidates {
		if c.Width*c.Height > m.Width*m.Height {
			m = c
		}
	}
	return m
}

// get the recent user media feed
func GetRecentUserMedias(user *goinsta.User) ([]*model.Media, error) {
	media := user.Feed()
	fmt.Println(media.Error())

	var images []*model.Media
	for media.Next() {
		for _, item := range media.Items {
			candidates := item.Images.Versions
			if len(candidates) == 0 {
				continue
			}
			m := getLargestCandidate(item.Images.Versions)
			// remove token from url
			mediaURL, err := cleanURL(m.URL)
			if err != nil {
				return nil, err
			}
			images = append(images, &model.Media{
				ID:          item.ID,
				URL:         mediaURL,
				IsLandscape: item.OriginalWidth > item.OriginalHeight,
				IsPicOfUser: item.PhotoOfYou,
				UserID:      user.ID,
				Username:    user.FullName,
				LikeCount:   item.Likes,
				PostedAt:    time.Unix(int64(item.Caption.CreatedAt), 0),
			})
		}
	}
	return images, nil
}

// check if user is followed
func (s *Session) isFollowed(name string) bool {
	following := s.GetAllFollowedUsers()
	for _, usr := range following {
		if usr.Username == name {
			return true
		}
	}
	return false
}

// compose the profile info
func (s *Session) GetProfileInfo(name string) (proto.User, error) {
	var userDetails proto.User
	// get the user by name
	user, err := s.GetUserByName(name)
	if err != nil {
		log.Println(err)
		return userDetails, err
	}
	// if user's account is private
	if user.IsPrivate {
		// check if user is followed
		if !s.isFollowed(name) {
			err = errors.New("account is private and not followed by you")
			return userDetails, err
		}
	}
	// get the recent user images
	images, err := GetRecentUserMedias(user)
	if err != nil {
		log.Println(err)
		return userDetails, err
	}
	// get the image with most likes
	img, err := GetImageWithMostLikes(images)
	if err != nil {
		log.Println(err)
		return userDetails, err
	}
	// set fields values
	userDetails.Username = user.Username
	userDetails.FullName = user.FullName
	userDetails.Description = user.Biography
	userDetails.FeaturedPicUrl = img.URL
	userDetails.ProfilePicUrl = user.ProfilePicURL
	userDetails.FollowerCount = int64(user.FollowerCount)

	return userDetails, nil
}

// send direct message
func (s *Session) SendDirectMessage(id, message, title string) (string, error) {
	response, err := s.insta.DirectMessage(id, message, title)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println(response)
	return response.Status, nil
}

// logout from instagram
func (s *Session) Logout() {
	s.insta.Logout()
}

// get the image with most likes
func GetImageWithMostLikes(images []*model.Media) (*model.Media, error) {
	mostLiked := &model.Media{URL: "https://akm-img-a-in.tosshub.com/indiatoday/images/story/201603/photostory_647_032216025305.jpg"}
	if len(images) > 0 {
		sort.SliceStable(images, func(i, j int) bool {
			return images[i].LikeCount > images[j].LikeCount
		})
		mostLiked = images[0]

		for _, img := range images {
			if img.IsLandscape {
				mostLiked = img
				if img.IsPicOfUser {
					break
				}
			}
		}
	}
	return mostLiked, nil
}

// remove the token from the url
func cleanURL(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	return u.String(), nil
}

func (s *Session) GetInsta() *goinsta.Instagram {
	return s.insta
}
