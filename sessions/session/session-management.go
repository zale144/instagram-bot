package session

import (
	"log"

	"github.com/zale144/instagram-bot/sessions/model"
)

// read request
type ReadReq struct {
	Key  string
	Resp chan *Session
}

// write request
type WriteReq struct {
	Key  string
	Val  *Session
	Resp chan bool
}

// remove request
type RemoveReq struct {
	Key  string
	Resp chan bool
}

// clear request
type ClearReq struct {
	Resp chan bool
}

// read, write, remove, clear request channels
var (
	reads   = make(chan *ReadReq)
	writes  = make(chan *WriteReq)
	removes = make(chan *RemoveReq)
	clears  = make(chan *ClearReq)
)

// manage multiple sessions connections with a goroutine
func Sessions() {
	sessions := make(map[string]*Session)
	for {
		select {
		case read := <-reads:
			read.Resp <- sessions[read.Key]
		case write := <-writes:
			sessions[write.Key] = write.Val
			write.Resp <- true
		case remove := <-removes:
			delete(sessions, remove.Key)    // delete value for key from map
			_, resp := sessions[remove.Key] // check if key still found in map
			remove.Resp <- !resp
		case clear := <-clears: // clear operation requested
			sessions = make(map[string]*Session) // re-instantiate the  map
			clear.Resp <- len(sessions) == 0     // send true if map is empty now
		}
	}
}

// get sessions by account struct
func GetSession(account *model.Account) (*Session, error) { // return sessions
	read := &ReadReq{ // instantiate a read request
		Key:  account.Username,
		Resp: make(chan *Session),
	}
	reads <- read      // send read request to reads channel
	ret := <-read.Resp // pull return value from the channel
	if ret == nil {
		s, err := NewSession(account)
		if err != nil || s == nil {
			log.Println(err)
			return nil, err
		}
		ret = s
	}
	return ret, nil
}

// add sessions to map
func SaveSession(session *Session, username string) {
	write := &WriteReq{ // instantiate a write request
		Key:  username,
		Val:  session,
		Resp: make(chan bool),
	}
	writes <- write // send write request to writes channel
	<-write.Resp    // pull boolean value from the channel
	log.Printf("Written '%s' to sessions map\n", write.Key)
}

// remove item from sessions map
func Remove(name string) {
	remove := &RemoveReq{ // instantiate a remove request
		Key:  name,
		Resp: make(chan bool),
	}
	removes <- remove // send remove request to removes channel
	<-remove.Resp     // pull boolean value from the channel
	log.Printf("Removed '%s' from sessions map\n", remove.Key)
}

// clear the sessions map
func Clear() {
	clear := &ClearReq{ // instantiate a clear request
		Resp: make(chan bool),
	}
	clears <- clear // send clear request to clears channel
	<-clear.Resp    // pull boolean value from the channel
	log.Println("sessions map cleared")
}
