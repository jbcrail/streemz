package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/fatih/color"
	"github.com/peterh/liner"
)

const (
	defaultPrompt = "twitter> "
	whitespace    = "\t\n\v\f\r "
)

var (
	magenta = color.New(color.FgMagenta).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
)

func parse(s string) (string, []string) {
	parser := regexp.MustCompile("[" + whitespace + "]+")
	s = strings.Trim(s, whitespace)
	tokens := parser.Split(s, -1)
	return strings.ToLower(tokens[0]), tokens[1:]
}

func toInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err
}

func printUserSummary(user *twitter.User) {
	fmt.Printf("%v: followers=%v following=%v statuses=%v likes=%v\n", magenta(user.ScreenName), user.FollowersCount, user.FriendsCount, user.StatusesCount, user.FavouritesCount)
}

func printUser(user *twitter.User) {
	fmt.Println(magenta(user.ScreenName))
	s := reflect.ValueOf(user).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s:%s = %v\n", typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func isRateLimitExceeded(resp *http.Response) bool {
	limit, _ := toInt(resp.Header["X-Rate-Limit-Remaining"][0])
	if limit == 0 {
		fmt.Println(red("rate limit exceeded"))
		return true
	}
	return false
}

func getTweetCount(arg string, initial int) int {
	tweetCount := initial
	i, err := toInt(arg)
	if len(arg) > 0 && err == nil {
		tweetCount = int(i)
	}
	return tweetCount
}

func homeTimeline(client *twitter.Client, count int) {
	tweets, resp, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: count,
	})

	if isRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		user := tweet.User
		fmt.Printf("[%v] %v\n", magenta(user.ScreenName), tweet.Text)
	}
}

func userTimeline(client *twitter.Client, name string, count int) {
	tweets, resp, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName:      name,
		Count:           count,
		IncludeRetweets: twitter.Bool(true),
	})

	if isRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range tweets {
		user := tweet.User
		fmt.Printf("[%v] %v\n", magenta(user.ScreenName), tweet.Text)
	}
}

func myFollowers(client *twitter.Client) {
	cursor := int64(-1)
	for {
		followers, resp, _ := client.Followers.List(&twitter.FollowerListParams{
			Cursor: cursor,
		})

		if isRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			printUserSummary(&user)
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

		if isRateLimitExceeded(resp) {
			break
		}

		for _, user := range followers.Users {
			printUserSummary(&user)
		}
		if followers.NextCursor == 0 {
			break
		}
		cursor = followers.NextCursor
	}
}

func myFriends(client *twitter.Client) {
	cursor := int64(-1)
	for {
		friends, resp, _ := client.Friends.List(&twitter.FriendListParams{
			Cursor: cursor,
		})

		if isRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			printUserSummary(&user)
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

		if isRateLimitExceeded(resp) {
			break
		}

		for _, user := range friends.Users {
			printUserSummary(&user)
		}
		if friends.NextCursor == 0 {
			break
		}
		cursor = friends.NextCursor
	}
}

func current(client *twitter.Client) {
	user, resp, _ := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})

	if isRateLimitExceeded(resp) {
		return
	}

	printUser(user)
}

func user(client *twitter.Client, name string) {
	user, resp, _ := client.Users.Show(&twitter.UserShowParams{
		ScreenName: name,
	})

	if isRateLimitExceeded(resp) {
		return
	}

	printUser(user)
}

func public(client *twitter.Client) {
	params := &twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
	}
	stream, _ := client.Streams.Sample(params)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		user := tweet.User
		fmt.Printf("[%v] %v\n", magenta(user.ScreenName), tweet.Text)
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

	if isRateLimitExceeded(resp) {
		return
	}

	for _, tweet := range search.Statuses {
		user := tweet.User
		fmt.Printf("[%v] %v\n", magenta(user.ScreenName), tweet.Text)
	}
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
		case "public":
			public(client)
		case "recent":
			if len(args) == 0 {
				homeTimeline(client, getTweetCount("", 20))
			} else {
				homeTimeline(client, getTweetCount(args[1], 20))
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
			fmt.Println("FOLLOWERS FRIENDS HELP PUBLIC QUIT RECENT TWEETS USER")
		default:
			fmt.Println("unknown command:", command)
			fmt.Println("FOLLOWERS FRIENDS HELP PUBLIC QUIT RECENT TWEETS USER")
		}
	}
}
