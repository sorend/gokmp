package flickr

import (
	"net/http"
	"context"
	"github.com/gomodule/oauth1/oauth"
)

func NewFlickr(ApiKey string, ApiSecret string, OAuthToken string, OAuthSecret string) *Flickr {
	hc := &http.Client{}
	ctx := context.WithValue(context.Background(), oauth.HTTPClient, hc)
	c := &oauth.Client{}
	return &Flickr{
		ApiKey: ApiKey,
		ApiSecret: ApiSecret,
		OAuthToken: OAuthToken,
		OAuthSecret: OAuthSecret,
		client: c,
		ctx: &ctx,
	}
}


type Flickr struct {
	ApiKey string
	ApiSecret string
	OAuthToken string
	OAuthSecret string
	client *oauth.Client
	ctx *context.Context
}

