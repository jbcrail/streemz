package tweets

import (
	"flag"
	"os"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("tweets", flag.ExitOnError)
	count := cmd.Int("count", 20, "")

	cmd.Parse(args)

	if cmd.NArg() == 0 {
		cmd.PrintDefaults()
		os.Exit(1)
	}

	tweets, resp, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName:      cmd.Arg(0),
		Count:           *count,
		IncludeRetweets: twitter.Bool(true),
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		cmdutil.PrintTweet(tweet)
	}
}
