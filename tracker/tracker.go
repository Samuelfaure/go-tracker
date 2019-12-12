// Package tracker creates a websocket server
// and sends total amount of connections to a messenger module
package tracker

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
	"os"
	"time"
)

type ChanChanges chan int

type ChanQuit chan os.Signal

type TrackerServer struct {
	Port      string
	Messenger Messenger
}

type Messenger interface {
	SendValue(value int)
}

// Custom context for Echo
type TrackerContext struct {
	echo.Context
	ChanChanges
}

func Init(t TrackerServer, quit ChanQuit) {
	changes := make(chan int)

	go count(t, changes)
	startServer(t, changes, quit)
}

func count(t TrackerServer, changes ChanChanges) {
	visitors := 0

	for {
		change, ok := <-changes

		if !ok {
			break
		}

		visitors += change
		t.Messenger.SendValue(visitors)
	}
}

func startServer(t TrackerServer, changes ChanChanges, quit ChanQuit) {
	e := echo.New()

	registerMiddlewares(e, changes)
	registerRoutes(e)

	go e.Start(t.Port)

	handleClose(e, changes, quit)
}

func registerMiddlewares(e *echo.Echo, changes ChanChanges) {
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tc := &TrackerContext{c, changes}
			return h(tc)
		}
	})
	e.Use(middleware.Recover())
}

func registerRoutes(e *echo.Echo) {
	e.GET("/ws/:token", func(c echo.Context) error {
		tc := c.(*TrackerContext)
		tc.Track(c)
		return tc.String(200, "OK")
	})
}

func (tc *TrackerContext) Track(c echo.Context) {
	changes := tc.ChanChanges

	changes <- 1 // Add one visitor
	startWebsocket(c)
	changes <- -1 // Remove one visitor
}

func startWebsocket(c echo.Context) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		for {
			token := c.Param("token") // TODO: decrypt / check token
			if err := websocket.Message.Receive(ws, &token); err != nil {
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
}

func handleClose(e *echo.Echo, changes ChanChanges, quit ChanQuit) {
	<-quit

	// Close count goroutine
	close(changes)

	// Close the server gracefully with 10 sec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("\nEcho server exited gracefully.")
}
