package middleware

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func Opentracing(spanName string) func(func(*http.Request) (*http.Response, error)) func(*http.Request) (*http.Response, error) {
	return func(f func(*http.Request) (*http.Response, error)) func(*http.Request) (*http.Response, error) {
		return func(req *http.Request) (*http.Response, error) {
			traceCtx, _ := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(req.Header),
			)

			// Create the span referring to the RPC client if available.
			// If traceCtx == nil, a root span will be created.
			serverSpan := opentracing.StartSpan(spanName, ext.RPCServerOption(traceCtx))
			defer serverSpan.Finish()

			ctx := opentracing.ContextWithSpan(req.Context(), serverSpan)
			return f(req.WithContext(ctx))
		}
	}
}
