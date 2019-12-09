# Go-tracker

Smol project to learn Golang and Kafka.

## What it does

Tl;dr: displays in your console the total amount of visitors on your websites.

The tracker starts an Echo server and listen to incoming websocket connections.
We add +1 for every websocket opening and -1 for every websocket closing. Nothing else is logged.

This data does NOT constitute PII (Personally Identifiable Information). Therefore, _from what I understand_, it should not fall under the GDPR obligations of consent.
(But maybe go and ask a real lawyer first.)

## Install

### Install Kafka

```bash
# Simple local install with yay package manager (Arch)
sudo yay kafka
```

For other distributions or for a production environment, refer to [Kafka documentation](https://kafka.apache.org/documentation/)

### Setup and open Kafka

```bash
# Start the kafka service
sudo systemctl start kafka.service

# Go to this folder
cd go-tracker

# Setup kafka topic (only once)
sh scripts/create_kafka_topic.sh

# Start the kafka console consumer
sh scripts/start_kafka_consumer.sh
```

Kafka is now open, ready to receive information.

### Install go-tracker

Run the `go-tracker` binary on your favorite server.
You can configure the messenger module and the tracker port in `main.go` (TODO: Create a separate config file)

### Install the client on your websites

Import the `tracker_client.js` file on your website's metadata.
You done !

## TODO

- Increase test coverage (priority)
- Write the JS client
- Create a config file to take configuration out of `main.go`
- Add encrypted token system so we only log from trusted sources
