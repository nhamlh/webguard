package sso

import (
	"fmt"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
	"golang.org/x/oauth2/google"
	"log"
)

type Oauth2Provider struct {
	cfg  *oauth2.Config
	user UserApi
}

func NewOauth2Provider(clientId, clientSecret, redirectURL, provider string) (*Oauth2Provider, error) {
	var user UserApi
	var endpoint oauth2.Endpoint
	switch provider {
	case "github":
		user = GithubUserApi{}
		endpoint = github.Endpoint
	case "gitlab":
		user = GitlabUserApi{}
		endpoint = gitlab.Endpoint
	case "google":
		user = GoogleUserApi{}
		endpoint = google.Endpoint
	default:
		return &Oauth2Provider{}, errors.New("Provider not supported")
	}

	oc := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     endpoint,
		Scopes:       user.Scopes(),
	}

	return &Oauth2Provider{
		cfg:  &oc,
		user: user,
	}, nil
}

//TODO: Implement challenge verifier
func (p *Oauth2Provider) Redirect(w http.ResponseWriter, r *http.Request) {
	//FIXME: Eliminate static state value
	u := p.cfg.AuthCodeURL("foobar")
	http.Redirect(w, r, u, http.StatusFound)
}

func (p *Oauth2Provider) GetToken(r *http.Request) (token *oauth2.Token, err error) {
	r.ParseForm()
	state := r.Form.Get("state")
	code := r.Form.Get("code")

	if state != "foobar" || code == "" {
		err = errors.New("Invalid response from authorization server")
		return
	}

	token, err = p.cfg.Exchange(r.Context(), code)

	return
}

func (p *Oauth2Provider) Email(token oauth2.Token) string {
	client := &http.Client{}

	resp, err := client.Do(p.user.Request(token))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	log.Println(token.AccessToken)
	log.Println(string(rawResp))

	var jsonObj map[string]interface{}
	if err = json.Unmarshal(rawResp, &jsonObj); err != nil {
		return ""
	}

	email := p.user.Parse(jsonObj)

	return email
}

// UserApi is an interface to get user object from various providers
// and extract email field used as our user's identity
type UserApi interface {
	Scopes() []string // Return scopes needed to get access user api
	Request(oauth2.Token) *http.Request
	Parse(map[string]interface{}) string
}

type GithubUserApi struct{}

func (g GithubUserApi) Scopes() []string {
	return []string{"email"}
}

func (g GithubUserApi) Request(token oauth2.Token) (req *http.Request) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token.AccessToken))
	return
}

func (g GithubUserApi) Parse(jsonObj map[string]interface{}) string {
	return jsonObj["email"].(string)
}

type GitlabUserApi struct{}

func (g GitlabUserApi) Scopes() []string {
	return []string{"read_user"}
}

func (g GitlabUserApi) Request(token oauth2.Token) (req *http.Request) {
	req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/user", nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	return
}

func (g GitlabUserApi) Parse(jsonObj map[string]interface{}) string {
	return jsonObj["email"].(string)
}

type GoogleUserApi struct{}

func (g GoogleUserApi) Scopes() []string {
	return []string{"https://www.googleapis.com/auth/userinfo.email"}
}

func (g GoogleUserApi) Request(token oauth2.Token) (req *http.Request) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	return
}

func (g GoogleUserApi) Parse(jsonObj map[string]interface{}) string {
	return jsonObj["email"].(string)
}
