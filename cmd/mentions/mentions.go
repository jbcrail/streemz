package mentions

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

func mentionTimeline(client *twitter.Client, count int) {
	tweets, resp, _ := client.Timelines.MentionTimeline(&twitter.MentionTimelineParams{
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
		mentionTimeline(client, getTweetCount("", 20))
	} else {
		mentionTimeline(client, getTweetCount(args[0], 20))
	}
}
