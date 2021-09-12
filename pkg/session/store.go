package session

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
)

//FIXME: This implementation has a drawback:
// Expired sessions might be orphaned and build up
// in the store.
type Store struct {
	Name string
	// store session as a map to its Id for easier management
	sessions map[string]Session
}

func NewSessionStore() *Store {
	return &Store{
		Name:     "Session",
		sessions: map[string]Session{},
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

	s, found := st.sessions[string(id)]
	if found {
		if s.IsExpired() {
			delete(st.sessions, s.Id)
			return &Session{}, false
		} else {
			return &s, true
		}

	}

	return &s, false
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
			MaxAge:   3600 * 24, // Expire in one day
			Path:     "/",
			HttpOnly: true,
		},
	}
}

// Destroy removes a session from the store and also invalidate the \\
// corresponding cookie
func (st *Store) Destroy(s Session, w http.ResponseWriter) {
	delete(st.sessions, s.Id)

	http.SetCookie(w, &http.Cookie{
		Name:   st.Name,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
}

// Save inserts a session into the store and also set a cookie to the response
func (st *Store) Save(s Session, w http.ResponseWriter) {
	st.sessions[s.Id] = s
	http.SetCookie(w, &s.Cookie)
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
