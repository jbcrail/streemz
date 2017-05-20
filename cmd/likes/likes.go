package likes

import (
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func myLikes(client *twitter.Client) {
	tweets, resp, _ := client.Favorites.List(&twitter.FavoriteListParams{})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		cmdutil.PrintTweet(tweet)
	}
}

func likes(client *twitter.Client, name string) {
	tweets, resp, _ := client.Favorites.List(&twitter.FavoriteListParams{
		ScreenName: name,
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		cmdutil.PrintTweet(tweet)
	}
}

func Run(client *twitter.Client, args []string) {
	if len(args) == 0 {
		myLikes(client)
	} else {
		likes(client, args[0])
	}
}
