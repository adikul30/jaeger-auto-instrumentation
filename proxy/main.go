package main

import (
	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	proxyPort   = 8000
	servicePort = 80
	serviceName = "SVC_NAME"
)

var (
	tracer opentracing.Tracer
	cfg jaegercfg.Configuration
)

// Create a structure to define the proxy functionality.
type Proxy struct{}

func contains(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	cfg = jaegercfg.Configuration{
		ServiceName: os.Getenv(serviceName),
		Sampler:     &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter:    &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, _ := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	tracer = opentracing.GlobalTracer()

	// Forward the HTTP request to the destination serviceA.

	//TODO : check for existing jaeger tracer headers. if they exist, forward the request. If they don't, add headers
	// serviceA name has to be unique for diff services
	fmt.Println("********************")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Printf("%v: %v\n", name, h)
		}
	}
	values, ok := req.Header["Svc_name"]
	if ok {
		fmt.Println("service_name: " + os.Getenv(serviceName))
		if contains(values, os.Getenv(serviceName)) {
			fmt.Println("originating from here ...")
			// originating from the host serviceA
			//todo: add jaeger params
			fmt.Println("this should be here: " + os.Getenv(serviceName))
			clientSpan := tracer.StartSpan("svc-A")
			ext.SpanKindRPCClient.Set(clientSpan)
			ext.HTTPUrl.Set(clientSpan, req.URL.String())
			ext.HTTPMethod.Set(clientSpan, "GET")
			defer clientSpan.Finish()
			tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

			res, duration, err := p.performOutboundRequest(req)
			// Notify the client if there was an error while forwarding the request.
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}

			// If the request was forwarded successfully, write the response back to
			// the client.
			p.writeResponse(w, res)

			// Print request and response statistics.
			p.printStats(req, res, duration)

		} else {
			// originating from other services
			fmt.Println("originating from other services ...")
			//todo: extract jaeger params
			spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
			serverSpan := tracer.StartSpan("svc-B", ext.RPCServerOption(spanCtx))
			defer serverSpan.Finish()

			res, duration, err := p.forwardRequest(req)

			// Notify the client if there was an error while forwarding the request.
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}

			// If the request was forwarded successfully, write the response back to
			// the client.
			p.writeResponse(w, res)

			// Print request and response statistics.
			p.printStats(req, res, duration)
		}
	}
}

func (p *Proxy) performOutboundRequest(req *http.Request) (*http.Response, time.Duration, error) {
	httpClient := http.Client{}
	destinationUrl := req.Header.Get("destination-ip")
	newUrl := fmt.Sprintf("http://%s:%d%s", destinationUrl, servicePort, req.RequestURI)
	newRequest, err := http.NewRequest(req.Method, newUrl, req.Body)
	start := time.Now()
	res, err := httpClient.Do(newRequest)
	duration := time.Since(start)
	return res, duration, err
}

func (p *Proxy) forwardRequest(req *http.Request) (*http.Response, time.Duration, error) {
	// Prepare the destination endpoint to forward the request to.
	proxyUrl := fmt.Sprintf("http://127.0.0.1:%d%s", servicePort, req.RequestURI)

	// Print the original URL and the proxied request URL.
	fmt.Printf("Original URL: http://%s:%d%s\n", req.Host, servicePort, req.RequestURI)
	fmt.Printf("Proxy URL: %s\n", proxyUrl)

	// Create an HTTP client and a proxy request based on the original request.
	httpClient := http.Client{}
	proxyReq, err := http.NewRequest(req.Method, proxyUrl, req.Body)

	// Capture the duration while making a request to the destination serviceA.
	start := time.Now()
	res, err := httpClient.Do(proxyReq)
	duration := time.Since(start)

	// Return the response, the request duration, and the error.
	return res, duration, err
}

func (p *Proxy) writeResponse(w http.ResponseWriter, res *http.Response) {
	// Copy all the header values from the response.
	for name, values := range res.Header {
		w.Header()[name] = values
	}

	// Set a special header to notify that the proxy actually serviced the request.
	w.Header().Set("Server", "amazing-proxy")

	// Set the status code returned by the destination serviceA.
	w.WriteHeader(res.StatusCode)

	// Copy the contents from the response body.
	io.Copy(w, res.Body)

	// Finish the request.
	res.Body.Close()
}

func (p *Proxy) printStats(req *http.Request, res *http.Response, duration time.Duration) {
	fmt.Printf("Request Duration: %v\n", duration)
	fmt.Printf("Request Size: %d\n", req.ContentLength)
	fmt.Printf("Response Size: %d\n", res.ContentLength)
	fmt.Printf("Response Status: %d\n\n", res.StatusCode)
}

func main() {
	// Listen on the predefined proxy port.
	fmt.Println("Service started: " + os.Getenv(serviceName))
	//initJaegerStuff()
	http.ListenAndServe(fmt.Sprintf(":%d", proxyPort), &Proxy{})
}

func initJaegerStuff() {
	cfg = jaegercfg.Configuration{
		ServiceName: os.Getenv(serviceName),
		Sampler:     &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter:    &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, _, _ := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	//defer closer.Close()

	tracer = opentracing.GlobalTracer()

}
