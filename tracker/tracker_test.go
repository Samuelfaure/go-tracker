package tracker

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
	"os"
	"testing"
)

type SpyMessenger struct {
	Calls int
	Value int
}

func (s *SpyMessenger) SendValue(value int) {
	s.Calls++
	s.Value = value
}

func StartTracker(m *SpyMessenger, quit ChanQuit) {
	tracker := TrackerServer{Port: ":1323", Messenger: m}

	Init(tracker, quit)
}

func TestInit(t *testing.T) {
	m := &SpyMessenger{Calls: 0, Value: 0}
	quit := make(ChanQuit)

	go StartTracker(m, quit)

	assert.Equal(t, 0, m.Value, "Counter should start at 0")

	conn, err := websocket.Dial("ws://localhost:1323/ws/test", "", "ws://localhost/")
	defer conn.Close()

	assert.Equal(t, 1, m.Calls, "Messenger should be called once")
	assert.Equal(t, 1, m.Value, "Messenger should receive the value 1")

	if err != nil {
		t.Fatalf("Error while connecting to websocket: %v", err)
		t.FailNow()
	}

	// Stop server with Ctrl-C
	quit <- os.Interrupt
}
