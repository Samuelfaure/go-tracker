package tracker

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
	"testing"
)

type SpyMessenger struct {
	Calls int
	Value int
}

func (s SpyMessenger) SendValue(value int) {
	s.Calls++
	s.Value = value
}

func StartTracker(m *SpyMessenger) {
	tracker := TrackerServer{Port: ":1323", Messenger: m}
	Init(tracker)
}

func TestInit(t *testing.T) {
	m := &SpyMessenger{Calls: 0, Value: 0}
	StartTracker(m)

	conn, err := websocket.Dial("ws://localhost:1323/ws/test", "", "ws://localhost/")
	defer conn.Close()

	if err != nil {
		t.Fatalf("Error while connecting to websocket: %v", err)
		t.FailNow()
	}

	if err = websocket.JSON.Send(conn, "test"); err != nil {
		t.Fatalf("Error while sending data to websocket: %v", err)
		t.FailNow()
	}

	assert.Equal(t, m.Value, 1, "Messenger should receive the value 1")
}
