package session

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	Name     string
	sessions []Session
}

func NewSessionStore() *Store {
	return &Store{
		Name: "Session",
	}
}

// Get returns a session if found in the request and not expired yet
func (st *Store) Get(r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie(st.Name)
	if err != nil {
		return &Session{}, false
	}

	id, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return &Session{}, false
	}

	var session Session
	for _, se := range st.sessions {
		if se.Id == string(id) {
			if se.IsExpired() {
				//TODO: Also remove expired session from store
				return &Session{}, false
			} else {
				return &se, true
			}
		}

	}

	return &session, false
}

func (st *Store) New() *Session {
	id := uuid.New().String()

	return &Session{
		Id:        id,
		CreatedAt: time.Now(),
		Values:    map[string]interface{}{},
		Cookie: http.Cookie{
			Name:     st.Name,
			Value:    base64.StdEncoding.EncodeToString([]byte(id)),
			MaxAge:   int(time.Second * 60),
			Path:     "/",
			HttpOnly: true,
		},
	}
}

// Destroy removes a session from the store and also invalidate the \\
// corresponding cookie
func (st *Store) Destroy(s Session, w *http.ResponseWriter) {
	http.SetCookie(*w, &http.Cookie{
		Name:   st.Name,
		Value:  "",
		MaxAge: -1,
	})
}

// Save inserts a session into the store and also set a cookie to the response
func (st *Store) Save(s *Session, w *http.ResponseWriter) {
	st.sessions = append(st.sessions, *s)
	http.SetCookie(*w, &s.Cookie)
}

type Session struct {
	Id        string
	Values    map[string]interface{}
	CreatedAt time.Time
	Cookie    http.Cookie
}

func (s *Session) IsExpired() bool {
	if s.Cookie.MaxAge <= 0 {
		return true
	}

	expiredTime := s.CreatedAt.Add(time.Duration(s.Cookie.MaxAge))

	if time.Now().After(expiredTime) {
		return true
	}

	return false
}
