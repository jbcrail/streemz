package mentions

import (
	"flag"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("mentions", flag.ExitOnError)
	count := cmd.Int("count", 20, "")

	cmd.Parse(args)

	tweets, resp, _ := client.Timelines.MentionTimeline(&twitter.MentionTimelineParams{
		Count: *count,
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		cmdutil.PrintTweet(tweet)
	}
}
