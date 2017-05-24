package recent

import (
	"flag"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("recent", flag.ExitOnError)
	count := cmd.Int("count", 20, "")
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.HomeTimelineParams{
		Count: *count,
	}

	tweets, resp, _ := client.Twitter.Timelines.HomeTimeline(&params)

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
