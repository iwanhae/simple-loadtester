package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var dummyMemory = []byte{}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("WARN: hostname not found")
		hostname = "UNKNOWN"
	}
	startDelay := flag.Int("start-delay", 5, "delay second")
	rps := flag.Int("rate-limit", 0, "allowed requests per second. set zero to disable")
	color := flag.String("color", "aqua", "this value will be added to http response headr with key of 'color'")
	flag.Parse()

	log.Println("startDelay", *startDelay)
	log.Println("rps:", *rps)
	log.Println("color:", *color)

	log.Println("Sleep")
	time.Sleep(time.Duration(*startDelay) * time.Second)
	log.Println("Wake")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) // Allow ALL
	if rps := *rps; rps != 0 {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(rps))))
	}
	e.Use(middleware.Static("./build"))

	e.GET("/delay/:duration", func(c echo.Context) error {
		duration := c.Param("duration")

		if ms, err := strconv.Atoi(duration); err == nil {
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}
		req := c.Request()
		c.Response().Header().Add("color", *color)
		return c.JSONPretty(http.StatusOK, map[string]interface{}{
			"hostname":    hostname,
			"host":        req.Host,
			"method":      req.Method,
			"url":         req.URL.String(),
			"headers":     req.Header,
			"remote_addr": req.RemoteAddr,
		}, "  ")
	})

	e.GET("/cpu/:length", func(c echo.Context) error {
		length, err := strconv.Atoi(c.Param("length"))
		if err != nil {
			return err
		}

		now := time.Now()
		pi := PiMultiThread(8, length)

		c.Response().Header().Add("color", *color)
		c.Response().Header().Add("color", *color)
		return c.JSONPretty(http.StatusOK, map[string]interface{}{
			"takes": time.Since(now).Milliseconds(),
			"pi":    fmt.Sprintf("%v", pi),
		}, "  ")
	})

	e.GET("/memory/:length", func(c echo.Context) error {
		length, err := strconv.Atoi(c.Param("length"))
		if err != nil {
			return err
		}

		now := time.Now()
		dummyMemory = make([]byte, length*1024*1024*1024)
		for i := range dummyMemory {
			dummyMemory[i] = '1'
		}

		c.Response().Header().Add("color", *color)
		return c.JSONPretty(http.StatusOK, map[string]interface{}{
			"takes":    time.Since(now).Milliseconds(),
			"length":   len(dummyMemory),
			"lengthKB": len(dummyMemory) / 1024,
			"lengthMB": len(dummyMemory) / 1024 / 1024,
			"lengthGB": len(dummyMemory) / 1024 / 1024 / 1024,
		}, "  ")
	})

	e.Logger.Fatal(e.Start(":8080"))
}

type Float64 struct {
	value float64
	lock  sync.RWMutex
}

func (f *Float64) Inc(diff float64) float64 {
	f.lock.Lock()
	ret := f.value + diff
	f.value = ret
	f.lock.Unlock()
	return ret
}

func (f *Float64) Get() float64 {
	f.lock.RLock()
	ret := f.value
	f.lock.RUnlock()
	return ret
}
func partialSum(kStart int, kOffset int, amount int) (sum float64) {
	for k := float64(kStart); k < float64(kStart+kOffset*amount); k += float64(kOffset) {
		sum += 1 / math.Pow(16, k) * (4/(8*k+1) - 2/(8*k+4) - 1/(8*k+5) - 1/(8*k+6))
	}

	return
}

func PiMultiThread(workers int, iteration int) float64 {
	wg := sync.WaitGroup{}
	ret := Float64{}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go (func(kStart int) {
			ret.Inc(partialSum(int(kStart), workers, iteration))
			wg.Done()
		})(i)
	}

	wg.Wait()
	return ret.Get()
}
