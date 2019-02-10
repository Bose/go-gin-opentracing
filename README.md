# go-gin-opentracing
[![](https://godoc.org/github.com/Bose/go-gin-opentracing?status.svg)](https://godoc.org/github.com/Bose/go-gin-opentracing) 
[![Go Report Card](https://goreportcard.com/badge/github.com/Bose/go-gin-opentracing)](https://goreportcard.com/report/github.com/Bose/go-gin-opentracing)
[![Release](https://img.shields.io/github/release/Bose/go-gin-opentracing.svg?style=flat-square)](https://Bose/go-gin-opentracing/releases) 

Gin Web Framework Open Tracing middleware

## Installation

`$ go get github.com/Bose/go-gin-opentracing`

### Deprecated functions (Feb 2019)
The following functions are deprecated since they use deprecated jaeger client functions and fixing them would require breaking changes.  

- InitDevelopment
- InitProduction
- InitMacDocker

## Usage

```go
package main

import (
	"fmt"
	"os"

	ginopentracing "github.com/Bose/go-gin-opentracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func main() {
	r := gin.Default()
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	tracer, reporter, closer, err := ginopentracing.InitTracing(fmt.Sprintf("go-gin-opentracing-example::%s", hostName), "localhost:5775", ginopentracing.WithEnableInfoLog(true))
	if err != nil {
		panic("unable to init tracing")
	}
	defer closer.Close()
	defer reporter.Close()
	opentracing.SetGlobalTracer(tracer)

	p := ginopentracing.OpenTracer([]byte("api-request-"))
	r.Use(p)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	r.Run(":29090")
}

```

See the [example.go file](https://github.com/github.com/Bose/go-gin-opentracing/blob/master/example/example.go)

