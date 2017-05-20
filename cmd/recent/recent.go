package recent

import (
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func getTweetCount(arg string, initial int) int {
	tweetCount := initial
	i, err := cmdutil.ToInt(arg)
	if len(arg) > 0 && err == nil {
		tweetCount = int(i)
	}
	return tweetCount
}

func homeTimeline(client *twitter.Client, count int) {
	tweets, resp, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: count,
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		cmdutil.PrintTweet(tweet)
	}
}

func Run(client *twitter.Client, args []string) {
	if len(args) == 0 {
		homeTimeline(client, getTweetCount("", 20))
	} else {
		homeTimeline(client, getTweetCount(args[0], 20))
	}
}
