package main

import "github.com/Samuelfaure/go-tracker/tracker"
import "github.com/Samuelfaure/go-tracker/messenger"

func main() {

	trackerServer := tracker.Server{":1323"}
	kafkaServer := messenger.KafkaServer{"tcp", "localhost:9092", "visitors", 0}

	tracker.Init(trackerServer, kafkaServer)
}
