package main

import "github.com/Samuelfaure/go-tracker/tracker"
import "github.com/Samuelfaure/go-tracker/messenger"

func main() {

	t := tracker.Server{":1323"}
	m := messenger.Server{"tcp", "localhost:9092", "visitors", 0}

	tracker.Init(t, m)
}
