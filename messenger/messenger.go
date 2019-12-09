// Package messenger transmit a single int value to a topic
package messenger

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
	"time"
)

type Messenger struct {
	Config
}

type Config struct {
	Protocol, URL, Topic string
	Partition            int
}

func (m Messenger) SendValue(value int) {
	msg := strconv.Itoa(value)

	conn, err := kafka.DialLeader(
		context.Background(),
		m.Config.Protocol,
		m.Config.URL,
		m.Config.Topic,
		m.Config.Partition)

	if err != nil {
		log.Print(err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.WriteMessages(kafka.Message{Value: []byte(msg)})
	conn.Close()
}
