package mentions

import (
	"flag"
	"fmt"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("mentions", flag.ExitOnError)
	count := cmd.Int("count", 20, "")
	maximum := cmd.Int64("max", 0, "")
	minimum := cmd.Int64("min", 0, "")
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.MentionTimelineParams{
		Count: *count,
	}

	if *maximum > 0 {
		// MaxID is inclusive
		params.MaxID = *maximum
	}

	if *minimum > 0 {
		// SinceID is exclusive, so we subtract 1 for consistency with maximum
		params.SinceID = *minimum - 1
	}

	tweets, resp, _ := client.Twitter.Timelines.MentionTimeline(&params)

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for i, tweet := range tweets {
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
