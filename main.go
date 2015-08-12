package main

import (
	"./models"
	"./modules/api"
	"./modules/setting"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/weisd/log"
)

const (
	VER = "0.1.0.0811"
)

func version(c *echo.Context) error {
	return c.JSON(200, api.ResOk(VER))
}

func main() {
	bootstraps()

	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/", version)

	// Start server
	e.Run(":1323")
}

func bootstraps() {
	setting.InitConfig()
	setting.InitServices()

	models.InitDatabaseConn()

	models.InitRedisPools()
	models.RedisCheckConn()

	log.Debug("%v", setting.Cfg)
}
