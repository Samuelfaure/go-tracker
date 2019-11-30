package tracker

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

type (
	ChanChanges chan int

	TrackerContext struct {
		echo.Context
		ChanChanges
	}
)

func Init() {
	changes := make(chan int)

	go count(changes)
	startServer(changes)
}

func startServer(changes ChanChanges) {
	e := echo.New()
	registerMiddlewares(e, changes)

	e.GET("/ws/:token", func(c echo.Context) error {
		tc := c.(*TrackerContext)
		tc.StartWebsocket(c)
		return tc.String(200, "OK")
	})

	e.Start(":1323")
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

func (tc *TrackerContext) StartWebsocket(c echo.Context) {
	changes := tc.ChanChanges
	handleOpen(changes)
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		for {
			if err := websocket.Message.Send(ws, "Connection ON"); err != nil {
				break
			}

			token := c.Param("token") // TODO decrypt / check token
			if err := websocket.Message.Receive(ws, &token); err != nil {
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	handleClose(changes)
}

func count(changes chan int) {
	visitors := 0

	for {
		change := <-changes
		visitors += change
		fmt.Printf("count: %d\n", visitors)
		// TODO: Send to kafka
	}
}

func handleOpen(changes chan int) {
	changes <- 1
}

func handleClose(changes chan int) {
	changes <- -1
}
