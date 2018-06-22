package ginopentracing

import (
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// OpenTracer - middleware that addes opentracing
func OpenTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// all before request is handled
		var span opentracing.Span
		if cspan, ok := c.Get("tracing-context"); ok {
			span = StartSpanWithParent(cspan.(opentracing.Span).Context(), "api-request-"+c.Request.Method, c.Request.Method, c.Request.URL.Path)

		} else {
			span = StartSpanWithHeader(&c.Request.Header, "api-request-"+c.Request.Method, c.Request.Method, c.Request.URL.Path)
		}
		defer span.SetTag(string(ext.HTTPStatusCode), c.Writer.Status()) // this must be before the defer finish to be properly located in the defer stack
		defer span.Finish()                                              // after all the other defers are completed.. finish the span
		c.Set("tracing-context", span)                                   // add the span to the context so it can be used for the duration of the request.
		c.Next()

		// after request is handled...
	}
}
