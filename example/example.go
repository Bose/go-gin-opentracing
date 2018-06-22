package main

import (
	"fmt"
	"os"

	"github.com/Bose/go-gin-opentracing"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

func main() {
	r := gin.Default()
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}

	tracer, closer, err := ginopentracing.Config.New(fmt.Sprintf("example.go::%s", hostName))
	if err == nil {
		fmt.Println("Setting global tracer")
		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)
	} else {
		fmt.Println("Can't enable tracing: ", err.Error())
	}

	p := ginopentracing.OpenTracer([]byte("api-request-"))
	r.Use(p)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	r.Run(":29090")
}
