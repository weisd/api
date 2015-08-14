package main

import (
	"net/http"
	"strings"

	"./models"
	"./models/user"
	"./modules/api"
	"./modules/setting"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/weisd/cache"
	_ "github.com/weisd/cache/redis"
	"github.com/weisd/echo-statistics"
	"github.com/weisd/jwt"
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

	if setting.Cfg.Debug {
		log.Info("SetDebug on")
		e.SetDebug(setting.Cfg.Debug)
	}

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(cache.EchoCacher(cache.Options{Adapter: "redis", AdapterConfig: `{"Addr":":6379"}`, Section: "test", Interval: 5}))
	e.Use(statistics.Statisticser())
	// 固定返回值
	e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		code := http.StatusInternalServerError
		msg := http.StatusText(code)
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code()
			msg = he.Error()
		}

		if e.Debug() {
			msg = err.Error()
		}

		log.Error(4, "http err  %v", err)

		c.JSON(200, api.ResErr(code, msg))

	})

	// Routes
	e.Get("/", version)
	e.Get("/test/add", func(c *echo.Context) error {
		u := new(user.User)
		u.Name = "weisd"

		act, err := user.Create(u)
		if err != nil {
			return c.JSON(200, api.ResErr(500, err.Error()))
		}

		return c.JSON(200, api.ResOk(act))
	})

	e.Get("/test/stat", func(c *echo.Context) error {
		status := statistics.StatisticsMap.GetMap()
		html := strings.Join(status["Fields"].([]string), " ")
		html += "<br>"
		data := status["Data"].([][]string)
		for i, l := 0, len(data); i < l; i++ {
			html += strings.Join(data[i], " ")
			html += "<br>"
		}

		return c.HTML(200, html)

	})

	e.Get("/test/cache/put", func(c *echo.Context) error {
		err := cache.Store(c).Put("name", "weisd", 10)
		if err != nil {
			return err
		}

		return c.String(200, "store ok")
	})

	e.Get("/test/cache/get", func(c *echo.Context) error {
		name := cache.Store(c).Get("name")

		return c.String(200, "get name %s", name)
	})

	var jwtSigningKeys = map[string]string{"da": "weisd"}

	j := e.Group("/jwt")
	j.Use(jwt.EchoJWTAuther(func(c *echo.Context) (key string, err error) {
		// get the clientId from header
		clientId := c.Request().Header.Get("client-id")
		key, ok := jwtSigningKeys[clientId]
		if !ok {
			return "", echo.NewHTTPError(http.StatusUnauthorized)
		}
		return key, nil
	}))

	j.Get("", func(c *echo.Context) error {
		return c.String(200, "jwt Access ok with claims %v", jwt.Claims(c))
	})

	e.Get("/test/jwt/token", func(c *echo.Context) error {
		claims := map[string]interface{}{"token": "weisd"}
		token, err := jwt.NewToken("weisd", claims)
		if err != nil {
			return err
		}
		// show the token use for test
		return c.String(200, "token : %s", token)
	})

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
