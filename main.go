package main

import (
	"github.com/Samuelfaure/go-tracker/messenger"
	"github.com/Samuelfaure/go-tracker/tracker"
	"os"
	"os/signal"
)

func main() {

	t := tracker.TrackerServer{Port: ":1323", Messenger: moduleMessenger()}

	//  Channel to close server gracefully
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	tracker.Init(t, quit)
}

func moduleMessenger() messenger.Messenger {
	c := messenger.Config{
		Protocol:  "tcp",
		URL:       "localhost:9092",
		Topic:     "visitors",
		Partition: 0}

	return messenger.Messenger{Config: c}
}
