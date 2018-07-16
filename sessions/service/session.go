package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/zale144/goinsta"
	"github.com/zale144/instagram-bot/sessions/model"
)

// Session implements the proto service Session
type Session struct {
	insta *goinsta.Instagram
}

// NewSession creates a new session
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
		msg := fmt.Sprintf("cannot import goinsta config with name: %s", account.Username)
		err := errors.New(msg)
		log.Println(err)
		return nil, err
	}
	// save to cache
	SaveSession(s, account.Username)
	return s, nil
}

// GetInsta is a getter for the *goinsta.Instagram instance
func (s *Session) GetInsta() *goinsta.Instagram {
	return s.insta
}
