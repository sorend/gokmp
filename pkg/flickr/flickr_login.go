
package flickr

import (
	"errors"
	"github.com/dghubble/oauth1"
)

const outOfBand = "oob"

var AuthenticateEndpoint = oauth1.Endpoint{
	RequestTokenURL: "https://www.flickr.com/services/oauth/request_token",
	AuthorizeURL: "https://www.flickr.com/services/oauth/authorize",
	AccessTokenURL: "https://www.flickr.com/services/oauth/access_token",
}

type FlickrLogin struct {
	config *oauth1.Config
	requestToken string
	requestSecret string
}

func NewFlickrLogin(ApiKey string, ApiSecret string) (*FlickrLogin, error) {
	if ApiKey == "" || ApiSecret == "" {
		return nil, errors.New("ApiKey or ApiSecret are empty")
	}
	config := oauth1.Config{
		ConsumerKey: ApiKey,
		ConsumerSecret: ApiSecret,
		CallbackURL: outOfBand,
		Endpoint: AuthenticateEndpoint,
	}
	return &FlickrLogin{
		config: &config,
		requestToken: "",
		requestSecret: "",
	}, nil
}

func (f *FlickrLogin) GetLoginURL() (string, error) {
	requestToken, requestSecret, err := f.config.RequestToken()
	if err != nil {
		return "", err
	}
	f.requestToken = requestToken
	f.requestSecret = requestSecret
	url, err := f.config.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}
	return url.String() + "&perms=read", nil
}

func (f *FlickrLogin) GetAccessToken(verifier string) (*oauth1.Token, error) {
	if verifier == "" {
		return nil, errors.New("Verifier PIN is empty")
	}
	accessToken, accessSecret, err := f.config.AccessToken(f.requestToken, f.requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	return oauth1.NewToken(accessToken, accessSecret), nil
}
