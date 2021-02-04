package main

import (
	"fmt"
	"net/http"
	"os"

	ginopentracing "github.com/Bose/go-gin-opentracing"
	db "github.com/Bose/go-gin-opentracing/example/db"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

var repo *db.Repository

func main() {

	repo = db.NewRepository()
	defer repo.Close()
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
	r.GET("/getBook/:genre", handleGetBook)

	r.Run(":29090")
}

type getBookRequest struct {
	Genre string `uri:"genre"`
}

func handleGetBook(ctx *gin.Context) {

	var span opentracing.Span
	if cspan, ok := ctx.Get("tracing-context"); ok {
		span = ginopentracing.StartSpanWithParent(cspan.(opentracing.Span).Context(), "get-book", ctx.Request.Method, ctx.Request.URL.Path)
	} else {
		span = ginopentracing.StartSpanWithHeader(&ctx.Request.Header, "get-book", ctx.Request.Method, ctx.Request.URL.Path)
	}
	defer span.Finish()
	var req getBookRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	book, err := repo.GetBooks(ctx, req.Genre)

	if err != nil {
		span.SetTag("error", true)
		span.LogFields(otlog.Error(err))
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, book)
}
