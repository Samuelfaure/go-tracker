package tracker

import (
	"golang.org/x/net/websocket"
	"testing"
)

func TestStart(t *testing.T) {
	go Start()
	origin := "ws://localhost/"
	url := "ws://localhost:1323/ws/test"
	conn, err := websocket.Dial(url, "", origin)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer conn.Close()

	if err = websocket.JSON.Send(conn, "test"); err != nil {
		t.Log(err)
		t.FailNow()
	}
}
