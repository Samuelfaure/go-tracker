package main

import "github.com/Samuelfaure/go-tracker/tracker"
import "github.com/Samuelfaure/go-tracker/messenger"

func main() {

	c := messenger.Config{
		Protocol:  "tcp",
		URL:       "localhost:9092",
		Topic:     "visitors",
		Partition: 0}

	m := messenger.Messenger{Config: c}

	t := tracker.TrackerServer{Port: ":1323", Messenger: m}

	tracker.Init(t)
}
