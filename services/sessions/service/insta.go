package service

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"time"
	"github.com/zale144/goinsta"
	"github.com/zale144/instagram-bot/services/sessions/client"
	"github.com/zale144/instagram-bot/services/sessions/model"
	proto "github.com/zale144/instagram-bot/services/sessions/proto"
)

// GetAllFollowedUsers will retrieve all followed Instagram users
func (s *Session) GetAllFollowedUsers() []proto.User {
	account := s.insta.Account
	followingUsers := []proto.User{}
	following := account.Following()

	for following.Next() {
		for i := range following.Users {
			details := ConvertUser(&following.Users[i])
			followingUsers = append(followingUsers, details)
		}
	}
	return followingUsers
}

// GetUsersByHashtag will get all Instagram users that have media with provided hashtag
func (s *Session) GetUsersByHashtag(hashtag string, limit int) []proto.User {
	users := []proto.User{}

	h := s.insta.NewHashtag(hashtag)
	for h.Next() {
		for i := range h.Sections {
			for _, i := range h.Sections[i].LayoutContent.Medias {
				if len(i.Item.Images.Versions) != 0 {
					fmt.Printf("url: %s  - user: %s\n", i.Item.Images.Versions[0].URL, i.Item.User.Username)
					if Contains(users, i.Item.User.Username) {
						continue
					}
					details := ConvertUser(&i.Item.User)
					users = append(users, details)
					if len(users) == limit {
						return users
					}
				}
			}
		}
	}
	return users
}

// Contains checks if the user is in the slice
func Contains(users []proto.User, username string) bool {
	for _, currentUser := range users {
		if currentUser.Username == username {
			return true
		}
	}
	return false
}

// GetUserByName will retrieve the user struct from Instagram API
func (s *Session) GetUserByName(name string) (*goinsta.User, error) {
	user, err := s.insta.Profiles.ByName(name)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// getLargestCandidate gets the largest picture from an image item
func getLargestCandidate(candidates []goinsta.Candidate) goinsta.Candidate {
	m := candidates[0]
	for _, c := range candidates {
		if c.Width*c.Height > m.Width*m.Height {
			m = c
		}
	}
	return m
}

// GetRecentUserMedias get the recent user media feed
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

// isFollowed checks if user is followed
func (s *Session) isFollowed(name string) bool {
	following := s.GetAllFollowedUsers()
	for _, usr := range following {
		if usr.Username == name {
			return true
		}
	}
	return false
}

// GetProfileInfo composes the profile info
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
	userDetails = ConvertUser(user)
	userDetails.FeaturedPicUrl = img.URL
	return userDetails, nil
}

// SendDirectMessage sends a direct message to Instagram user
func (s *Session) SendDirectMessage(id, message, title string) (string, error) {
	response, err := s.insta.DirectMessage(id, message, title)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println(response)
	return response.Status, nil
}

// Follow follows the user with provided username
func (s *Session) Follow(username string) (proto.User, error) {
	user, err := s.GetUserByName(username)
	if err != nil {
		log.Println(err)
		return proto.User{}, err
	}
	err = user.Follow()
	if err != nil {
		log.Println(err)
		return proto.User{}, err
	}
	return ConvertUser(user), nil
}

// Logout will logout the user from Instagram
func (s *Session) Logout() {
	s.insta.Logout()
}

// GetImageWithMostLikes gets the image with most likes
func GetImageWithMostLikes(images []*model.Media) (*model.Media, error) {
	// use this temporary image in case no suitable background image is found on the user's profile
	mostLiked := &model.Media{URL: "https://akm-img-a-in.tosshub.com/indiatoday/images/story/201603/photostory_647_032216025305.jpg"}
	if len(images) > 0 {
		sort.SliceStable(images, func(i, j int) bool {
			return images[i].LikeCount > images[j].LikeCount
		})
		mostLiked = images[0]

		for _, img := range images {
			// RPC to the facedetect service which runs an OpenCV face detection
			// and returns the number of detected faces in the provided image url
			p, err := client.GetNumberOfFaces(img.URL)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("FOR URL: ", img.URL)
			fmt.Printf("Found %v persons\n", p)
			if img.IsLandscape {
				mostLiked = img
				if img.IsPicOfUser || p > 0 {
					break
				}
			}
		}
	}
	return mostLiked, nil
}

// ConvertUser converts the 'goinsta.User' to 'proto.User'
func ConvertUser(user *goinsta.User) (p proto.User) {
	p.Username = user.Username
	p.ProfilePicUrl = user.ProfilePicURL
	p.Description = user.Biography
	p.FollowerCount = int64(user.FollowerCount)
	p.FullName = user.FullName
	return p
}

// cleanURL removes the token from the url
func cleanURL(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	return u.String(), nil
}
