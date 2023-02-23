package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("WARN: hostname not found")
		hostname = "UNKNOWN"
	}

	rps := flag.Int("rate-limit", 0, "allowed requests per second. set zero to disable")
	color := flag.String("color", "aqua", "this value will be added to http response headr with key of 'color'")
	flag.Parse()

	log.Println("rps:", *rps)
	log.Println("color:", *color)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) // Allow ALL
	if rps := *rps; rps != 0 {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(rps))))
	}
	e.Use(middleware.Static("./build"))

	e.GET("/delay/:duration", func(c echo.Context) error {
		duration := c.Param("duration")
		if dur, err := time.ParseDuration(duration); err == nil {
			time.Sleep(dur)
		}
		req := c.Request()
		c.Response().Header().Add("color", *color)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"hostname":    hostname,
			"host":        req.Host,
			"method":      req.Method,
			"url":         req.URL.String(),
			"headers":     req.Header,
			"remote_addr": req.RemoteAddr,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
