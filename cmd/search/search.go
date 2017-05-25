package search

import (
	"flag"
	"fmt"
	"strings"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("search", flag.ExitOnError)
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	search, resp, _ := client.Twitter.Search.Tweets(&twitter.SearchTweetParams{
		Query: strings.Join(cmd.Args(), " "),
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for i, tweet := range search.Statuses {
		if *json {
			cmdutil.PrintTweetAsJson(tweet)
		} else if *full {
			cmdutil.PrintExtendedTweet(tweet)
		} else {
			if i == 0 {
				fmt.Println()
			}
			cmdutil.PrintTweet(tweet)
			fmt.Println()
		}
	}
}
