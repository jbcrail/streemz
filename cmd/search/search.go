package search

import (
	"flag"
	"strings"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("search", flag.ExitOnError)
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	search, resp, _ := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: strings.Join(cmd.Args(), " "),
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range search.Statuses {
		if *json {
			cmdutil.PrintTweetAsJson(tweet)
		} else if *full {
			cmdutil.PrintExtendedTweet(tweet)
		} else {
			cmdutil.PrintTweet(tweet)
		}
	}
}
