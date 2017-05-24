package tweets

import (
	"flag"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("tweets", flag.ExitOnError)
	count := cmd.Int("count", 20, "")
	maximum := cmd.Int64("max", 0, "")
	minimum := cmd.Int64("min", 0, "")
	retweets := cmd.Bool("retweets", true, "")
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.UserTimelineParams{
		Count:           *count,
		IncludeRetweets: twitter.Bool(*retweets),
	}

	if *maximum > 0 {
		// MaxID is inclusive
		params.MaxID = *maximum
	}

	if *minimum > 0 {
		// SinceID is exclusive, so we subtract 1 for consistency with maximum
		params.SinceID = *minimum - 1
	}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	}

	tweets, resp, _ := client.Twitter.Timelines.UserTimeline(&params)

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
