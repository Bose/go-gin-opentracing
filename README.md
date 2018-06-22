# go-gin-opentracing
[![](https://godoc.org/github.com/Bose/go-gin-opentracing?status.svg)](https://godoc.org/github.com/Bose/go-gin-opentracing) 

Gin Web Framework Open Tracing middleware

## Installation

`$ go get github.com/Bose/go-gin-opentracing`

## Usage

```go
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

```

See the [example.go file](https://github.com/github.com/Bose/go-gin-opentracing/blob/master/example/example.go)

