package main

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

const (
	VER = "0.1.0.0811"
)

type ApiRes struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ApiOK(data interface{}) ApiRes {
	return ApiRes{Code: 200, Status: "ok", Message: "ok", Data: data}
}

func ApiErr(code int, msg string) ApiRes {
	return ApiRes{Code: code, Status: "err", Message: msg, Data: nil}
}

func version(c *echo.Context) error {
	return c.JSON(200, ApiOK(VER))
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/", version)

	// Start server
	e.Run(":1323")
}
