package client

import (
	"net/http"
	"time"

	"github.com/jbcrail/streemz/cache"

	"github.com/dghubble/go-twitter/twitter"
)

type Client struct {
	Twitter *twitter.Client
	Cache   cache.Cache
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		Twitter: twitter.NewClient(httpClient),
		Cache:   cache.NewInMemoryCache(),
	}
}

func (client *Client) DefaultExpiry() time.Duration {
	return time.Duration(15) * time.Minute
}

func (client *Client) AuthorizedUser() *twitter.User {
	key := "authorized-user"
	if client.Cache.Exists(key) {
		return client.Cache.Get(key).(*twitter.User)
	}
	user, _, _ := client.Twitter.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	client.Cache.Set(key, user, client.DefaultExpiry())
	return user
}
