package followers

import (
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func myFollowers(client *twitter.Client) {
	cursor := int64(-1)
	for {
		followers, resp, _ := client.Followers.List(&twitter.FollowerListParams{
			Cursor: cursor,
		})

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			cmdutil.PrintUserSummary(&user)
		}
		if followers.NextCursor == 0 {
			break
		}
		cursor = followers.NextCursor
	}
}

func followers(client *twitter.Client, name string) {
	cursor := int64(-1)
	for {
		followers, resp, _ := client.Followers.List(&twitter.FollowerListParams{
			ScreenName: name,
			Cursor:     cursor,
		})

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			cmdutil.PrintUserSummary(&user)
		}
		if followers.NextCursor == 0 {
			break
		}
		cursor = followers.NextCursor
	}
}

func Run(client *twitter.Client, args []string) {
	if len(args) == 0 {
		myFollowers(client)
	} else {
		followers(client, args[0])
	}
}
