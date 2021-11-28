package oauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	redirectBaseUri        = "http://localhost"
	redirectLandingBaseUri = "https://developers.pinterest.com/apps/"
	oAuthUri               = "https://www.pinterest.com"
	apiUri                 = "https://api.pinterest.com"
)

type OAuth struct {
	apiUri             string
	appId              string
	appSecret          string
	oAuthUri           string
	scope              string
	redirectUri        string
	redirectLandingUri string
	redirectPort       int
	browserPath        string
}

type OAuthConfig struct {
	AppId        string
	AppSecret    string
	Scope        string
	RedirectPort int
	BrowserPath  string
}

func NewOAuth(cfg OAuthConfig) *OAuth {
	return &OAuth{
		appId:              cfg.AppId,
		appSecret:          cfg.AppSecret,
		scope:              cfg.Scope,
		redirectUri:        fmt.Sprintf("%s:%d/", redirectBaseUri, cfg.RedirectPort),
		redirectLandingUri: fmt.Sprintf("%s/%s/", redirectLandingBaseUri, cfg.AppId),
		redirectPort:       cfg.RedirectPort,
		oAuthUri:           oAuthUri,
		apiUri:             apiUri,
		browserPath:        cfg.BrowserPath,
	}
}

func (o *OAuth) CreateAccessToken() (string, error) {

	authCode, err := getAuthCode(o)
	if err != nil {
		return "", err
	}

	accessToken, err := exchangeAuthCode(o, authCode)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func getAuthCode(o *OAuth) (string, error) {

	pinterestOAuthUri := o.oAuthUri
	appId := o.appId
	redirectUri := o.redirectUri
	scope := o.scope
	redirectLandingUri := o.redirectLandingUri
	redirectPort := o.redirectPort
	oauthState := uuid.New().String()

	url := fmt.Sprintf("%s/oauth/?consumer_id=%s&redirect_uri=%s&scope=%s&response_type=code&state=%s", pinterestOAuthUri, appId, redirectUri, scope, oauthState)

	// open browser to grant access
	cmd := exec.Command(
		o.browserPath,
		url,
	)

	if err := cmd.Start(); err != nil {
		return "", errors.New(fmt.Sprintf("unable to open browser. Error %s\n", err.Error()))
	}

	rs := newRedirectServer(oauthState, redirectLandingUri, redirectPort)

	// this will block until access is granted or permitted
	return rs.GetCode(), nil
}

func exchangeAuthCode(o *OAuth, authCode string) (string, error) {

	apiUrl := fmt.Sprintf("%s/%s", o.apiUri, "v5/oauth/token")
	base64Auth := base64Auth(o.appId, o.appSecret)
	c := &http.Client{}

	data := url.Values{}
	data.Set("code", authCode)
	data.Set("redirect_uri", o.redirectUri)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to create http request. Error: %s\n", err.Error()))
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic  %s", base64Auth))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := c.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to send http request. Error: %s\n", err.Error()))
	}

	defer res.Body.Close()
	resDatabytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to read response. Error: %s\n", err.Error()))
	}

	resp := struct {
		AccessToken           string `json:"access_token"`
		ResponseType          string `json:"response_type"`
		TokenType             string `json:"token_type"`
		ExpiresIn             int    `json:"expires_in"`
		RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
		Scope                 string `json:"scope"`
	}{}

	err = json.Unmarshal(resDatabytes, &resp)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to unmarshal response. Error: %s\n", err.Error()))
	}

	return resp.AccessToken, nil
}

func base64Auth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
