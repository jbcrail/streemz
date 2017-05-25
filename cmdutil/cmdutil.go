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

	likes_icon   = "\u2764"
	retweet_icon = "\u21C4"
)

func ToInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err
}

func IndentTextBlock(s string) string {
	lines := strings.Split(s, "\n")
	if lines[len(lines)-1] == "" {
		return "    " + strings.Join(lines[0:len(lines)-1], "\n    ") + "\n"
	}
	return "    " + strings.Join(lines, "\n    ")
}

func LocalizeTime(s string) string {
	t, _ := time.Parse(time.RubyDate, s)
	return t.Local().String()
}

func PrintUserSummary(user *twitter.User) {
	s := fmt.Sprintf("%v %v : ", blue(user.Name), magenta("@"+user.ScreenName))
	s += fmt.Sprintf("%v tweets  %v likes  ", user.StatusesCount, user.FavouritesCount)
	s += fmt.Sprintf("%v following  %v followers\n", user.FriendsCount, user.FollowersCount)

	if user.ProfileImageURLHttps != "" {
		s += IndentTextBlock(fmt.Sprintf("Profile photo: %v\n", user.ProfileImageURLHttps))
	}

	if user.Location != "" {
		s += IndentTextBlock(fmt.Sprintf("Location: %v\n", user.Location))
	}

	if user.URL != "" {
		s += IndentTextBlock(fmt.Sprintf("URL: %v\n", user.URL))
	}

	s += IndentTextBlock(fmt.Sprintf("Joined on %v\n", user.CreatedAt))

	fmt.Print(IndentTextBlock(s))
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
	s := fmt.Sprintf("%v - @%v - %v\n", bold(user.Name), magenta(user.ScreenName), LocalizeTime(tweet.CreatedAt))
	s += fmt.Sprintf("https://twitter.com/%v/status/%v\n", user.ScreenName, tweet.IDStr)
	s += fmt.Sprintln(IndentTextBlock(tweet.Text))
	retweets := tweet.RetweetCount
	likes := tweet.FavoriteCount
	if tweet.RetweetedStatus != nil {
		retweets = tweet.RetweetedStatus.RetweetCount
		likes = tweet.RetweetedStatus.FavoriteCount
	}
	s += fmt.Sprintf("    %v %v  %v %v", retweet_icon, blue(retweets), likes_icon, green(likes))
	fmt.Println(IndentTextBlock(s))
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
