package friends

import (
	"flag"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("friends", flag.ExitOnError)
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.FriendListParams{
		Cursor: int64(-1),
	}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	}

	for {
		friends, resp, _ := client.Twitter.Friends.List(&params)

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			if *json {
				cmdutil.PrintUserAsJson(&user)
			} else if *full {
				cmdutil.PrintUser(&user)
			} else {
				cmdutil.PrintUserSummary(&user)
			}
		}

		if friends.NextCursor == 0 {
			break
		}
		params.Cursor = friends.NextCursor
	}
}
