package followers

import (
	"flag"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("followers", flag.ExitOnError)

	cmd.Parse(args)

	params := twitter.FollowerListParams{
		Cursor: int64(-1),
	}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	}

	for {
		followers, resp, _ := client.Followers.List(&params)

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			cmdutil.PrintUserSummary(&user)
		}

		if followers.NextCursor == 0 {
			break
		}

		params.Cursor = followers.NextCursor
	}
}
