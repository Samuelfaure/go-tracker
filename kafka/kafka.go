package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"strconv"
	"time"
)

func SendCount(count int) {
	topic := "visitors"
	partition := 0
	trackerCount := strconv.Itoa(count)

	conn, _ := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.WriteMessages(kafka.Message{Value: []byte(trackerCount)})
	conn.Close()
}
