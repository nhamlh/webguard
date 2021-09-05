package web

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"fmt"
	"github.com/google/uuid"
)

var sessionStore = NewSessionStore()

type SessionStore struct {
	Name   string
	Expire time.Duration
	list   []Session
}

type Session struct {
	Id     string    `json:"id"`
	Value  string    `json:"value"`
	Expire time.Time `json:"expire"`
}

// New creates new session
func (s *SessionStore) New() (Session, error) {
	session := Session{
		Id:     uuid.New().String(),
		Expire: time.Now().Add(time.Second * s.Expire),
	}

	s.list = append(s.list, session)

	return session, nil
}

func (s *SessionStore) Get(r http.Request) (Session, bool) {
	cookie, err := r.Cookie(s.Name)
	if err != nil {
		fmt.Println(err.Error())
		return Session{}, false
	}

	payload, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		fmt.Println("Unable to decode cookie payload")
		return Session{}, false
	}

	requestSession := Session{}
	err = json.Unmarshal(payload, &requestSession)
	if err != nil {
		fmt.Println("Unable to unmarshal cookie payload")
		return Session{}, false
	}

	for _, session := range s.list {
		if session.Id == requestSession.Id {
			return requestSession, true
		}
	}

	return Session{}, false
}

func (s *SessionStore) Marshal(se Session) (http.Cookie, error) {
	value, err := json.Marshal(se)
	if err != nil {
		return http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:    s.Name,
		Value:   base64.StdEncoding.EncodeToString([]byte(value)),
		Path:    "/",
		Expires: se.Expire,
	}

	return cookie, nil
}

func (s *SessionStore) GetById(id string) (Session, bool) {
	for _, session := range s.list {
		if session.Id == id {
			return session, true
		}
	}

	return Session{}, false
}

func (s *SessionStore) Delete(id string) bool {
	for index, session := range s.list {
		if session.Id == id {
			s.list = append(s.list[:index], s.list[index+1:]...)
			return true
		}
	}

	return false
}

func (s *SessionStore) Print() {
	fmt.Print(s.list)
}

func (s *Session) IsExpired() bool {
	if s.Expire.Before(time.Now()) {
		return true
	}

	return false
}

func NewSessionStore() SessionStore {
	return SessionStore{
		Name:   "Authorization",
		Expire: 960,
	}
}
