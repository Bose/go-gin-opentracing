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
