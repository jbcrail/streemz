package recent

import (
	"flag"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("recent", flag.ExitOnError)
	count := cmd.Int("count", 20, "")

	cmd.Parse(args)

	tweets, resp, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: *count,
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		cmdutil.PrintTweet(tweet)
	}
}
