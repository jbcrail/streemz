package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmd/followers"
	"github.com/jbcrail/streemz/cmd/friends"
	"github.com/jbcrail/streemz/cmd/home"
	"github.com/jbcrail/streemz/cmd/likes"
	"github.com/jbcrail/streemz/cmd/mentions"
	"github.com/jbcrail/streemz/cmd/public"
	"github.com/jbcrail/streemz/cmd/search"
	"github.com/jbcrail/streemz/cmd/tweets"
	"github.com/jbcrail/streemz/cmd/user"

	"github.com/dghubble/oauth1"
	"github.com/peterh/liner"
)

const (
	whitespace = "\t\n\v\f\r "
)

var commands = []string{"favorites", "followers", "friends", "help", "home", "likes", "mentions", "public", "search", "tweets", "user"}

func parse(s string) (string, []string) {
	parser := regexp.MustCompile("[" + whitespace + "]+")
	s = strings.Trim(s, whitespace)
	tokens := parser.Split(s, -1)
	return strings.ToLower(tokens[0]), tokens[1:]
}

func usage() {
	fmt.Println(strings.Join(commands, " "))
}

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := client.NewClient(httpClient)

	if len(os.Args) == 1 {
		user := client.AuthorizedUser()
		prompt := fmt.Sprintf("[@%v]: ", user.ScreenName)
		RunEvaluatePrintLoop(client, prompt)
	} else {
		RunEvaluatePrint(client, os.Args[1], os.Args[2:])
	}
}

func RunEvaluatePrint(client *client.Client, command string, args []string) {
	switch command {
	case "favorites":
		fallthrough
	case "likes":
		likes.Run(client, args)
	case "followers":
		followers.Run(client, args)
	case "friends":
		friends.Run(client, args)
	case "mentions":
		mentions.Run(client, args)
	case "public":
		public.Run(client, args)
	case "home":
		home.Run(client, args)
	case "search":
		search.Run(client, args)
	case "tweets":
		tweets.Run(client, args)
	case "user":
		user.Run(client, args)
	case "help":
		usage()
	default:
		fmt.Println("unknown command:", command)
		usage()
	}
}

func RunEvaluatePrintLoop(client *client.Client, prompt string) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	for {
		cmd, err := line.Prompt(prompt)
		if err == liner.ErrPromptAborted || err == io.EOF {
			break
		} else if err != nil {
			log.Print("Error reading line: ", err)
			continue
		}

		command, args := parse(cmd)
		line.AppendHistory(cmd)

		if command == "" {
			continue
		} else if command == "quit" {
			break
		}

		RunEvaluatePrint(client, command, args)
	}
}
