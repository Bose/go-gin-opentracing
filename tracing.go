// Utilities for working with tracing.
// dependencies:
//   go get github.com/opentracing/opentracing-go
//   go get github.com/uber/jaeger-client-go

package ginopentracing

import (
	"net/http"
	"runtime"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Config - the open tracing config singleton
var Config jaegercfg.Configuration

// InitProduction - init a production tracer environment
// example: Create the default tracer and schedule its closing when main returns.
//    func main() {
//	tracing.InitProduction("jaegeragent.svc.cluster.local:6831")
//	tracer, closer, _ := tracing.Config.New("passport-gigya-user-access") // the service name is the param to New()
//	defer closer.Close()
//	opentracing.SetGlobalTracer(tracer)
//
func InitProduction(sampleProbability float64, tracingAgentHostPort []byte) {
	Config = jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: sampleProbability,
		},
		Reporter: reporterConfig(tracingAgentHostPort),
	}
}

// InitDevelopment - init a production tracer environment
// example: Create the default tracer and schedule its closing when main returns.
//    func main() {
//  tracing.InitDevelopment() # defaults to "localhost:6831" for tracing agent
//	tracer, closer, _ := tracing.Config.New("passport-gigya-user-access") // the service name is the param to New()
//	defer closer.Close()
//	opentracing.SetGlobalTracer(tracer)
//
func InitDevelopment(tracingAgentHostPort []byte) {
	Config = jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: reporterConfig(tracingAgentHostPort),
	}
}

// InitMacDocker - init a production tracer environment
// example: Create the default tracer and schedule its closing when main returns.
//    func main() {
//  tracing.InitMacDocker() # defaults to "host.docker.internal:6831 for tracing agent
//	tracer, closer, _ := tracing.Config.New("passport-gigya-user-access") // the service name is the param to New()
//	defer closer.Close()
//	opentracing.SetGlobalTracer(tracer)
//
func InitMacDocker(tracingAgentHostPort []byte) {
	var reporter *jaegercfg.ReporterConfig
	if tracingAgentHostPort != nil {
		reporter = reporterConfig(tracingAgentHostPort)
	} else {
		reporter = reporterConfig([]byte("host.docker.internal:6831"))
	}
	Config = jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: reporter,
	}
}

func reporterConfig(hostPort []byte) *jaegercfg.ReporterConfig {
	agentPort := "localhost:6831"
	if hostPort != nil {
		agentPort = string(hostPort)
	}
	return &jaegercfg.ReporterConfig{
		LogSpans:           true,
		LocalAgentHostPort: agentPort,
	}
}

// StartSpan will start a new span with no parent span.
func StartSpan(operationName, method, path string) opentracing.Span {
	return StartSpanWithParent(nil, operationName, method, path)
}

func StartDBSpanWithParent(parent opentracing.SpanContext, operationName, dbInstance, dbType, dbStatement string) opentracing.Span {
	options := []opentracing.StartSpanOption{opentracing.Tag{Key: ext.SpanKindRPCServer.Key, Value: ext.SpanKindRPCServer.Value}}
	if len(dbInstance) > 0 {
		options = append(options, opentracing.Tag{Key: string(ext.DBInstance), Value: dbInstance})
	}
	if len(dbType) > 0 {
		options = append(options, opentracing.Tag{Key: string(ext.DBType), Value: dbType})
	}

	if parent != nil {
		options = append(options, opentracing.ChildOf(parent))
	}

	return opentracing.StartSpan(operationName, options...)
}

// StartSpanWithParent will start a new span with a parent span.
// example:
//      span:= StartSpanWithParent(c.Get("tracing-context"),
func StartSpanWithParent(parent opentracing.SpanContext, operationName, method, path string) opentracing.Span {
	options := []opentracing.StartSpanOption{
		opentracing.Tag{Key: ext.SpanKindRPCServer.Key, Value: ext.SpanKindRPCServer.Value},
		opentracing.Tag{Key: string(ext.HTTPMethod), Value: method},
		opentracing.Tag{Key: string(ext.HTTPUrl), Value: path},
		opentracing.Tag{Key: "current-goroutines", Value: runtime.NumGoroutine()},
	}

	if parent != nil {
		options = append(options, opentracing.ChildOf(parent))
	}

	return opentracing.StartSpan(operationName, options...)
}

// StartSpanWithHeader will look in the headers to look for a parent span before starting the new span.
// example:
//  func handleGet(c *gin.Context) {
//     span := StartSpanWithHeader(&c.Request.Header, "api-request", method, path)
//     defer span.Finish()
//     c.Set("tracing-context", span) // add the span to the context so it can be used for the duration of the request.
//     bosePersonID := c.Param("bosePersonID")
//     span.SetTag("bosePersonID", bosePersonID)
//
func StartSpanWithHeader(header *http.Header, operationName, method, path string) opentracing.Span {
	var wireContext opentracing.SpanContext
	if header != nil {
		wireContext, _ = opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*header))
	}
	span := StartSpanWithParent(wireContext, operationName, method, path)
	span.SetTag("current-goroutines", runtime.NumGoroutine())
	return span
	// return StartSpanWithParent(wireContext, operationName, method, path)
}

// InjectTraceID injects the span ID into the provided HTTP header object, so that the
// current span will be propogated downstream to the server responding to an HTTP request.
// Specifying the span ID in this way will allow the tracing system to connect spans
// between servers.
//
//  Usage:
//          // resty example
// 	    r := resty.R()
//	    injectTraceID(span, r.Header)
//	    resp, err := r.Get(fmt.Sprintf("http://localhost:8000/users/%s", bosePersonID))
//
//          // galapagos_clients example
//          c := galapagos_clients.GetHTTPClient()
//          req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8000/users/%s", bosePersonID))
//          injectTraceID(span, req.Header)
//          c.Do(req)
func InjectTraceID(ctx opentracing.SpanContext, header http.Header) {
	opentracing.GlobalTracer().Inject(
		ctx,
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header))
}
