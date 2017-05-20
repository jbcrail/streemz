package friends

import (
	"flag"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("friends", flag.ExitOnError)

	cmd.Parse(args)

	params := twitter.FriendListParams{
		Cursor: int64(-1),
	}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	}

	for {
		friends, resp, _ := client.Friends.List(&params)

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			cmdutil.PrintUserSummary(&user)
		}

		if friends.NextCursor == 0 {
			break
		}
		params.Cursor = friends.NextCursor
	}
}
