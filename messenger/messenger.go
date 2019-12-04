// Package messenger transmit a single int value to a topic
package messenger

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
	"time"
)

type KafkaServer struct {
	Protocol, Url, Topic string
	Partition            int
}

func SendValue(k KafkaServer, value int) {
	msg := strconv.Itoa(value)

	conn, err := kafka.DialLeader(context.Background(), k.Protocol, k.Url, k.Topic, k.Partition)
	if err != nil {
		log.Print(err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.WriteMessages(kafka.Message{Value: []byte(msg)})
	conn.Close()
}
