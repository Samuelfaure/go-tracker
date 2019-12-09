// Package tracker creates a websocket server
// and sends total amount of connections to a messenger module
package tracker

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

type (
	ChanChanges chan int

	// Using Messenger interface for injection dependency
	TrackerServer struct {
		Port      string
		Messenger Messenger
	}

	Messenger interface {
		SendValue(value int)
	}

	// Custom context for Echo
	TrackerContext struct {
		echo.Context
		ChanChanges
	}
)

func Init(t TrackerServer) {
	changes := make(chan int)

	go count(t, changes)
	startServer(t, changes)
}

func count(t TrackerServer, changes chan int) {
	visitors := 0

	for {
		change := <-changes
		visitors += change
		t.Messenger.SendValue(visitors)
	}
}

func startServer(t TrackerServer, changes ChanChanges) {
	e := echo.New()

	registerMiddlewares(e, changes)
	registerRoutes(e)
	e.Start(t.Port)
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
			if err := websocket.Message.Send(ws, "Connection ON"); err != nil {
				break
			}

			token := c.Param("token") // TODO: decrypt / check token
			if err := websocket.Message.Receive(ws, &token); err != nil {
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
}
