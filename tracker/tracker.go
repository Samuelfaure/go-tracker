package tracker

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

func handleOpen() {
	fmt.Printf("+1\n")
}

func handleClose() {
	fmt.Printf("-1\n")
}

func connected(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		defer handleClose()

		for {
			err := websocket.Message.Send(ws, "Hello, Client!")
			if err != nil {
				c.Logger().Error(err)
			}

			msg := c.Param("msg")
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			handleOpen()
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "../public")
	e.GET("/ws/:msg", connected)
	e.Logger.Fatal(e.Start(":1323"))
}
