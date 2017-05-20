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
	full := cmd.Bool("full", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Parse(args)

	params := &twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
	}
	stream, _ := client.Streams.Sample(params)

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
