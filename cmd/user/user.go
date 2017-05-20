package user

import (
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func current(client *twitter.Client) {
	user, resp, _ := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	cmdutil.PrintUser(user)
}

func user(client *twitter.Client, name string) {
	user, resp, _ := client.Users.Show(&twitter.UserShowParams{
		ScreenName: name,
	})

	if cmdutil.IsRateLimitExceeded(resp) {
		return
	}

	cmdutil.PrintUser(user)
}

func Run(client *twitter.Client, args []string) {
	if len(args) == 0 {
		current(client)
	} else {
		user(client, args[0])
	}
}
