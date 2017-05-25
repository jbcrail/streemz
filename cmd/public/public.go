package public

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbcrail/streemz/client"
	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *client.Client, args []string) {
	cmd := flag.NewFlagSet("public", flag.ExitOnError)
	filter := cmd.String("filter", "", "")
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	var stream *twitter.Stream

	if *filter == "" {
		params := &twitter.StreamSampleParams{
			StallWarnings: twitter.Bool(true),
		}
		stream, _ = client.Twitter.Streams.Sample(params)
	} else {
		params := &twitter.StreamFilterParams{
			StallWarnings: twitter.Bool(true),
			Track:         []string{*filter},
		}
		stream, _ = client.Twitter.Streams.Filter(params)
	}

	fmt.Println()

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if *json {
			cmdutil.PrintTweetAsJson(*tweet)
		} else if *full {
			cmdutil.PrintExtendedTweet(*tweet)
		} else {
			cmdutil.PrintTweet(*tweet)
			fmt.Println()
		}
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	stream.Stop()
}
