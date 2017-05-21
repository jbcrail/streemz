package client

import (
	"net/http"

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
