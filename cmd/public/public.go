package public

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbcrail/streemz/cmdutil"

	"github.com/dghubble/go-twitter/twitter"
)

func Run(client *twitter.Client, args []string) {
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
		stream, _ = client.Streams.Sample(params)
	} else {
		params := &twitter.StreamFilterParams{
			StallWarnings: twitter.Bool(true),
			Track:         []string{*filter},
		}
		stream, _ = client.Streams.Filter(params)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if *json {
			cmdutil.PrintTweetAsJson(*tweet)
		} else if *full {
			cmdutil.PrintExtendedTweet(*tweet)
		} else {
			cmdutil.PrintTweet(*tweet)
		}
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	stream.Stop()
}
