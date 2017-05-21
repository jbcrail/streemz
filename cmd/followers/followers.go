package followers

import (
	"flag"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("followers", flag.ExitOnError)
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.FollowerListParams{
		Cursor: int64(-1),
	}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	}

	for {
		followers, resp, _ := client.Twitter.Followers.List(&params)

		if cmdutil.IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			if *json {
				cmdutil.PrintUserAsJson(&user)
			} else if *full {
				cmdutil.PrintUser(&user)
			} else {
				cmdutil.PrintUserSummary(&user)
			}
		}

		if followers.NextCursor == 0 {
			break
		}

		params.Cursor = followers.NextCursor
	}
}
