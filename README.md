# go-gin-prometheus
[![](https://godoc.org/github.com/BoseCorp/go-gin-opentracing?status.svg)](https://godoc.org/github.com/BoseCorp/go-gin-opentracing) 

Gin Web Framework Open Tracing middleware

## Installation

`$ go get github.com/BoseCorp/go-gin-opentracing`

## Usage

```go
package main

import (
	"github.com/BoseCorp/go-gin-opentracing"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	p := ginopentracing.OpenTracer([]byte("api-request-"))
	r.Use(p)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	r.Run(":29090")
}

```

See the [example.go file](https://github.com/github.com/BoseCorp/go-gin-opentracing/blob/master/example/example.go)

