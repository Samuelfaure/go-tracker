package tracker

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
	"os"
	"testing"
	"time"
)

type SpyMessenger struct {
	Value int
}

func (s *SpyMessenger) SendValue(value int) {
	s.Value = value
}

func StartTracker(m *SpyMessenger, quit ChanQuit) {
	tracker := TrackerServer{Port: ":1323", Messenger: m}

	Init(tracker, quit)
}

func TestInit(t *testing.T) {
	m := &SpyMessenger{Value: 0}
	quit := make(ChanQuit)

	go StartTracker(m, quit)

	initialTests(t, m)
	websocketTests(t, m)

	// Stop server
	quit <- os.Interrupt
}

func initialTests(t *testing.T, m *SpyMessenger) {
	assert.Equal(t, 0, m.Value, "Counter should start at 0")
}

func websocketTests(t *testing.T, m *SpyMessenger) {
	oneWebsocketOpenTest(t, m)
	oneWebsocketCloseTest(t, m)
	twoWebsocketOpenTest(t, m)
	twoWebsocketOpenOneCloseTest(t, m)
}

func oneWebsocketOpenTest(t *testing.T, m *SpyMessenger) {
	websocket1 := openNewWebsocket(t)
	defer websocket1.Close()

	assert.Equal(t, 1, m.Value, "Messenger should receive the value 1")
}

// TODO : Make it work without sleep()
func oneWebsocketCloseTest(t *testing.T, m *SpyMessenger) {
	websocket1 := openNewWebsocket(t)

	websocket1.Close()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 0, m.Value, "Messenger should receive the value 0")
}

func twoWebsocketOpenTest(t *testing.T, m *SpyMessenger) {
	websocket1 := openNewWebsocket(t)
	websocket2 := openNewWebsocket(t)
	defer websocket1.Close()
	defer websocket2.Close()

	assert.Equal(t, 2, m.Value, "Messenger should receive the value 2")
}

// TODO : Make it work without sleep()
func twoWebsocketOpenOneCloseTest(t *testing.T, m *SpyMessenger) {
	websocket1 := openNewWebsocket(t)
	websocket2 := openNewWebsocket(t)
	defer websocket1.Close()

	websocket2.Close()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, m.Value, "Messenger should receive the value 1")
}

func openNewWebsocket(t *testing.T) (conn *websocket.Conn) {
	conn, err := websocket.Dial("ws://localhost:1323/ws/test", "", "ws://localhost/")
	if err != nil {
		t.Fatalf("Error while connecting to websocket: %v", err)
		t.FailNow()
	}
	return
}
