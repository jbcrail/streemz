package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/fatih/color"
)

var (
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
		fmt.Printf("%s:%s = %v\n", typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func PrintTweet(tweet twitter.Tweet) {
	user := tweet.User
	fmt.Printf("[%v] %v\n", magenta(user.ScreenName), tweet.Text)
}

func IsRateLimitExceeded(resp *http.Response) bool {
	limit, _ := ToInt(resp.Header["X-Rate-Limit-Remaining"][0])
	if limit == 0 {
		fmt.Println(red("rate limit exceeded"))
		return true
	}
	return false
}
