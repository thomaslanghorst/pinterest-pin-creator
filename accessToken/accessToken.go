package accessToken

import (
	"errors"
	"io/ioutil"
	"os"
	"pin-creator/accessToken/oauth"
)

const (
	scope = "pins:read,pins:write,boards:read,boards:write"
)

type AccessTokenFileHandlerInterface interface {
	Read() (string, error)
	Write(token string) error
}

type AccessTokenFileHandler struct {
	filePath string
}

type AccessTokenCreatorInterface interface {
	NewToken(appId string, appSecret string) (string, error)
}

type AccessTokenCreator struct {
	browserPath  string
	redirectPort int
}

func NewAccessAccessTokenCreator(browserPath string, redirectPort int) *AccessTokenCreator {
	return &AccessTokenCreator{
		browserPath:  browserPath,
		redirectPort: redirectPort,
	}
}

func NewAccessTokenFileHandler(filePath string) *AccessTokenFileHandler {
	return &AccessTokenFileHandler{
		filePath: filePath,
	}
}

func (h *AccessTokenFileHandler) Read() (string, error) {
	f, err := os.Open(h.filePath)
	defer f.Close()

	if err != nil {
		return "", err
	} else {
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}

func (h *AccessTokenFileHandler) Write(accessToken string) error {
	file, err := os.Create(h.filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(accessToken))
	if err != nil {
		return err
	}

	return nil
}

func (c *AccessTokenCreator) NewToken(appId string, appSecret string) (string, error) {
	if appId == "" || appSecret == "" {
		return "", errors.New("no APP_ID and APP_SECRET are provided as env variables")
	}

	oauth := oauth.NewOAuth(oauth.OAuthConfig{
		AppId:        appId,
		AppSecret:    appSecret,
		Scope:        scope,
		RedirectPort: c.redirectPort,
		BrowserPath:  c.browserPath,
	})

	return oauth.CreateAccessToken()
}
