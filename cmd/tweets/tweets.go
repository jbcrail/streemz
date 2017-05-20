package tweets

import (
	"fmt"

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

func userTimeline(client *twitter.Client, name string, count int) {
	tweets, resp, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName:      name,
		Count:           count,
		IncludeRetweets: twitter.Bool(true),
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		cmdutil.PrintTweet(tweet)
	}
}

func Run(client *twitter.Client, args []string) {
	if len(args) > 1 {
		userTimeline(client, args[0], getTweetCount(args[1], 20))
	} else if len(args) == 1 {
		userTimeline(client, args[0], getTweetCount("", 20))
	} else {
		fmt.Println("Usage: tweets NAME [N]")
	}
}
