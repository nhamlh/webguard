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
	cfg *oauth2.Config
	pc  *ProviderConfig
}

func NewOauth2Provider(clientId, clientSecret, redirectURL string, pc ProviderConfig) (*Oauth2Provider, error) {
	oc := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     pc.endpoint,
		Scopes:       pc.scopes,
	}

	return &Oauth2Provider{
		cfg: &oc,
		pc:  &pc,
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

	resp, err := client.Do(p.pc.userReq(token))
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

	email := p.pc.parse(jsonObj)

	return email
}

type ProviderConfig struct {
	endpoint oauth2.Endpoint
	scopes   []string
	userReq  func(oauth2.Token) *http.Request
	parse    func(map[string]interface{}) string
}

var GithubProvider = ProviderConfig{
	endpoint: github.Endpoint,
	scopes:   []string{"email"},
	userReq: func(t oauth2.Token) (req *http.Request) {
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Token %s", t.AccessToken))
		return

	},
	parse: func(jsonObj map[string]interface{}) string {
		return jsonObj["email"].(string)
	},
}

var GitlabProvider = ProviderConfig{
	endpoint: gitlab.Endpoint,
	scopes:   []string{"read_user"},
	userReq: func(t oauth2.Token) (req *http.Request) {
		req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/user", nil)
		if err != nil {
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
		return
	},
	parse: func(jsonObj map[string]interface{}) string {
		return jsonObj["email"].(string)
	},
}

var GoogleProvider = ProviderConfig{
	endpoint: google.Endpoint,
	scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
	userReq: func(t oauth2.Token) (req *http.Request) {
		req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
		if err != nil {
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
		return
	},
	parse: func(jsonObj map[string]interface{}) string {
		return jsonObj["email"].(string)
	},
}

func NewOktaProvider(domain string) ProviderConfig {
	return ProviderConfig{
		endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/oauth2/v1/authorize", domain),
			TokenURL: fmt.Sprintf("https://%s/oauth2/v1/token", domain),
		},
		scopes: []string{"openid", "email"},
		userReq: func(t oauth2.Token) (req *http.Request) {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/oauth2/v1/userinfo", domain), nil)
			if err != nil {
				return
			}

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
			return
		},
		parse: func(jsonObj map[string]interface{}) string {
			return jsonObj["email"].(string)
		},
	}
}
