package friends

import (
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func myFriends(client *twitter.Client) {
	cursor := int64(-1)
	for {
		friends, resp, _ := client.Friends.List(&twitter.FriendListParams{
			Cursor: cursor,
		})

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			cmdutil.PrintUserSummary(&user)
		}
		if friends.NextCursor == 0 {
			break
		}
		cursor = friends.NextCursor
	}
}

func friends(client *twitter.Client, name string) {
	cursor := int64(-1)
	for {
		friends, resp, _ := client.Friends.List(&twitter.FriendListParams{
			ScreenName: name,
			Cursor:     cursor,
		})

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			cmdutil.PrintUserSummary(&user)
		}
		if friends.NextCursor == 0 {
			break
		}
		cursor = friends.NextCursor
	}
}

func Run(client *twitter.Client, args []string) {
	if len(args) == 0 {
		myFriends(client)
	} else {
		friends(client, args[0])
	}
}
