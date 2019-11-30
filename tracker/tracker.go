package tracker

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

type (
	ChannelChanges chan int

	TrackerContext struct {
		echo.Context
		ChannelChanges
	}
)

func Init() {
	visitors := 0
	changes := make(chan int)

	go count(&visitors, changes)

	// Creating server
	e := echo.New()

	// Middlewares
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tc := &TrackerContext{c, changes}
			return h(tc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ws/:msg", func(c echo.Context) error {
		tc := c.(*TrackerContext)
		tc.StartWebsocket(c)
		return tc.String(200, "OK")
	})

	// Manage closing
	e.Logger.Fatal(e.Start(":1323"))
}

func (tc *TrackerContext) StartWebsocket(c echo.Context) {
	changes := tc.ChannelChanges
	handleOpen(changes)

	websocket.Handler(func(ws *websocket.Conn) {
		// Rules for closing and opening
		defer handleClose(changes)
		defer ws.Close()

		// Websocket
		for {
			err := websocket.Message.Send(ws, "Connection ON")
			if err != nil {
				break
			}

			msg := c.Param("msg") // TODO replace with a token, break if incorrect
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
}

func count(visitors *int, changes chan int) {
	for {
		change := <-changes
		*visitors += change
		// TODO: Send to kafka
	}
}

func handleOpen(changes chan int) {
	changes <- 1
}

func handleClose(changes chan int) {
	changes <- -1
}
