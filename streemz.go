package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/peterh/liner"
)

const (
	defaultPrompt = "twitter> "
	whitespace    = "\t\n\v\f\r "
)

func parse(s string) (string, []string) {
	parser := regexp.MustCompile("[" + whitespace + "]+")
	s = strings.Trim(s, whitespace)
	tokens := parser.Split(s, -1)
	return strings.ToLower(tokens[0]), tokens[1:]
}

func getTweetCount(arg string, initial int) int {
	tweetCount := initial
	i, err := ToInt(arg)
	if len(arg) > 0 && err == nil {
		tweetCount = int(i)
	}
	return tweetCount
}

func homeTimeline(client *twitter.Client, count int) {
	tweets, resp, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: count,
	})

	if IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		PrintTweet(tweet)
	}
}

func mentionTimeline(client *twitter.Client, count int) {
	tweets, resp, _ := client.Timelines.MentionTimeline(&twitter.MentionTimelineParams{
		Count: count,
	})

	if IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		PrintTweet(tweet)
	}
}

func userTimeline(client *twitter.Client, name string, count int) {
	tweets, resp, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName:      name,
		Count:           count,
		IncludeRetweets: twitter.Bool(true),
	})

	if IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		PrintTweet(tweet)
	}
}

func myFollowers(client *twitter.Client) {
	cursor := int64(-1)
	for {
		followers, resp, _ := client.Followers.List(&twitter.FollowerListParams{
			Cursor: cursor,
		})

		if IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			PrintUserSummary(&user)
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

		if IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			PrintUserSummary(&user)
		}
		if followers.NextCursor == 0 {
			break
		}
		cursor = followers.NextCursor
	}
}

func myLikes(client *twitter.Client) {
	tweets, resp, _ := client.Favorites.List(&twitter.FavoriteListParams{})

	if IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		PrintTweet(tweet)
	}
}

func likes(client *twitter.Client, name string) {
	tweets, resp, _ := client.Favorites.List(&twitter.FavoriteListParams{
		ScreenName: name,
	})

	if IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		PrintTweet(tweet)
	}
}

func myFriends(client *twitter.Client) {
	cursor := int64(-1)
	for {
		friends, resp, _ := client.Friends.List(&twitter.FriendListParams{
			Cursor: cursor,
		})

		if IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			PrintUserSummary(&user)
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

		if IsRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			PrintUserSummary(&user)
		}
		if friends.NextCursor == 0 {
			break
		}
		cursor = friends.NextCursor
	}
}

func current(client *twitter.Client) {
	user, resp, _ := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})

	if IsRateLimitExceeded(resp) {
		return
	}

	PrintUser(user)
}

func user(client *twitter.Client, name string) {
	user, resp, _ := client.Users.Show(&twitter.UserShowParams{
		ScreenName: name,
	})

	if IsRateLimitExceeded(resp) {
		return
	}

	PrintUser(user)
}

func public(client *twitter.Client) {
	params := &twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
	}
	stream, _ := client.Streams.Sample(params)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		PrintTweet(*tweet)
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	stream.Stop()
}

func search(client *twitter.Client, keywords []string) {
	search, resp, _ := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: strings.Join(keywords, " "),
	})

	if IsRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range search.Statuses {
		PrintTweet(tweet)
	}
}

func usage() {
	fmt.Println("FAVORITES FOLLOWERS FRIENDS HELP LIKES MENTIONS PUBLIC QUIT RECENT SEARCH TWEETS USER")
}

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")

	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	for {
		cmd, err := line.Prompt(defaultPrompt)
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

		switch command {
		case "favorites":
			fallthrough
		case "likes":
			if len(args) == 0 {
				myLikes(client)
			} else {
				likes(client, args[0])
			}
		case "followers":
			if len(args) == 0 {
				myFollowers(client)
			} else {
				followers(client, args[0])
			}
		case "friends":
			if len(args) == 0 {
				myFriends(client)
			} else {
				friends(client, args[0])
			}
		case "mentions":
			if len(args) == 0 {
				mentionTimeline(client, getTweetCount("", 20))
			} else {
				mentionTimeline(client, getTweetCount(args[0], 20))
			}
		case "public":
			public(client)
		case "recent":
			if len(args) == 0 {
				homeTimeline(client, getTweetCount("", 20))
			} else {
				homeTimeline(client, getTweetCount(args[0], 20))
			}
		case "search":
			search(client, args)
		case "tweets":
			if len(args) > 1 {
				userTimeline(client, args[0], getTweetCount(args[1], 20))
			} else if len(args) == 1 {
				userTimeline(client, args[0], getTweetCount("", 20))
			} else {
				fmt.Println("Usage: tweets NAME [N]")
			}
		case "user":
			if len(args) == 0 {
				current(client)
			} else {
				user(client, args[0])
			}
		case "help":
			usage()
		default:
			fmt.Println("unknown command:", command)
			usage()
		}
	}
}
