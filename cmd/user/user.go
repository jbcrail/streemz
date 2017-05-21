package user

import (
	"flag"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("user", flag.ExitOnError)
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.UserShowParams{}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	} else {
		user := client.AuthorizedUser()
		params.ScreenName = user.ScreenName
	}

	user, resp, _ := client.Twitter.Users.Show(&params)

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	if *json {
		cmdutil.PrintUserAsJson(user)
	} else {
		cmdutil.PrintUser(user)
	}
}
