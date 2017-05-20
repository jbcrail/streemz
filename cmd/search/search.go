package search

import (
	"strings"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, keywords []string) {
	search, resp, _ := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: strings.Join(keywords, " "),
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range search.Statuses {
		cmdutil.PrintTweet(tweet)
	}
}
