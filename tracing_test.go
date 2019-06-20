package ginopentracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matryer/is"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/sirupsen/logrus"
)

func TestGeneral(t *testing.T) {
	is := is.New(t)

	tr := &mocktracer.MockTracer{}
	r := gin.New()

	srv := httptest.NewServer(r)
	defer srv.Close()

	opentracing.SetGlobalTracer(tr)

	p := OpenTracer([]byte("api-request-"))
	r.Use(p)

	_, err := http.Get(srv.URL)
	is.NoErr(err)

	spans := tr.FinishedSpans()
	is.True(len(spans) == 1)
	t.Log(spans[0].OperationName)
	is.True(spans[0].OperationName == "api-request-GET")

	logrus.SetLevel(logrus.DebugLevel)
	tracer, _, closer, err := InitTracing("go-gin-opentracing-example::localhost", "localhost:5775", WithEnableInfoLog(true), WithSampleProbability(1.0))
	if err != nil {
		panic("unable to init tracing")
	}
	defer closer.Close()
	s := tracer.StartSpan("dummyspan")
	t.Log(s)

}
