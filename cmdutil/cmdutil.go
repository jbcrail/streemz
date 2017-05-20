package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/fatih/color"
)

var (
	blue    = color.New(color.FgBlue).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
)

func ToInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err
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

func PrintTweet(tweet twitter.Tweet) {
	user := tweet.User
	fmt.Printf("[%v] %v\n", magenta(user.ScreenName), tweet.Text)
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
