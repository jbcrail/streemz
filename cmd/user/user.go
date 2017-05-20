package user

import (
	"flag"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
	cmd := flag.NewFlagSet("user", flag.ExitOnError)
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := twitter.UserShowParams{}

	if cmd.NArg() > 0 {
		params.ScreenName = cmd.Arg(0)
	} else {
		user, _, _ := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
		params.ScreenName = user.ScreenName
	}

	user, resp, _ := client.Users.Show(&params)

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	if *json {
		cmdutil.PrintUserAsJson(user)
	} else {
		cmdutil.PrintUser(user)
	}
}
