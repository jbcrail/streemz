package likes

import (
	"flag"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("likes", flag.ExitOnError)
	count := cmd.Int("count", 20, "")
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.FavoriteListParams{
		Count: *count,
	}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	}

	tweets, resp, _ := client.Favorites.List(&params)

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		if *json {
			cmdutil.PrintTweetAsJson(tweet)
		} else if *full {
			cmdutil.PrintExtendedTweet(tweet)
		} else {
			cmdutil.PrintTweet(tweet)
		}
	}
}
