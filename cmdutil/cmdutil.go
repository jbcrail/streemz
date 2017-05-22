package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/fatih/color"
)

var (
	blue    = color.New(color.FgBlue).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
)

func ToInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err
}

func IndentTextBlock(s string) string {
	return "\t" + strings.Replace(strings.TrimRight(s, "\n"), "\n", "\n\t", -1)
}

func LocalizeTime(s string) string {
	t, _ := time.Parse(time.RubyDate, s)
	return t.Local().String()
}

func PrintUserSummary(user *twitter.User) {
	fmt.Printf("%v: followers=%v following=%v statuses=%v likes=%v\n", magenta(user.ScreenName), user.FollowersCount, user.FriendsCount, user.StatusesCount, user.FavouritesCount)
}

func PrintUser(user *twitter.User) {
	fmt.Println(magenta(user.ScreenName))
	s := reflect.ValueOf(user).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s:%s = %v\n", blue(typeOfT.Field(i).Name), green(f.Type()), f.Interface())
	}
}

func PrintUserAsJson(user *twitter.User) {
	body, err := json.Marshal(user)
	if err != nil {
		fmt.Println(red("failed to serialize user to JSON"))
		return
	}
	fmt.Println(string(body))
}

func PrintTweet(tweet twitter.Tweet) {
	user := tweet.User
	fmt.Printf("%v - @%v - %v\n", bold(user.Name), magenta(user.ScreenName), LocalizeTime(tweet.CreatedAt))
	fmt.Printf("https://twitter.com/%v/status/%v\n", user.ScreenName, tweet.IDStr)
	fmt.Println(IndentTextBlock(tweet.Text))
	retweets := tweet.RetweetCount
	likes := tweet.FavoriteCount
	if tweet.RetweetedStatus != nil {
		retweets = tweet.RetweetedStatus.RetweetCount
		likes = tweet.RetweetedStatus.FavoriteCount
	}
	fmt.Printf("\t⇄ %v  ❤ %v\n", blue(retweets), green(likes))
}

func PrintExtendedTweet(tweet twitter.Tweet) {
	s := reflect.ValueOf(&tweet).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s:%s = %v\n", blue(typeOfT.Field(i).Name), green(f.Type()), f.Interface())
	}
}

func PrintTweetAsJson(tweet twitter.Tweet) {
	body, err := json.Marshal(tweet)
	if err != nil {
		fmt.Println(red("failed to serialize tweet to JSON"))
		return
	}
	fmt.Println(string(body))
}

func IsRateLimitExceeded(resp *http.Response) bool {
	val := resp.Header.Get("X-Rate-Limit-Remaining")
	if val == "" {
		fmt.Println(red("failed to get rate limit"))
		return true
	}
	limit, _ := ToInt(val)
	if limit == 0 {
		fmt.Println(red("rate limit exceeded"))
		return true
	}
	return false
}
