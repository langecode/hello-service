package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

const ServiceName = "hello-service"

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	log.Info().Str("TraceID", span.SpanContext().TraceID().String()).Msg("Echo service")

	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	handler := http.HandlerFunc(httpHandler)
	telemetryHandler := otelhttp.NewHandler(handler, ServiceName)
	http.Handle("/hello", telemetryHandler)

	http.HandleFunc("/healthz", healthzHandler)
	initMetrics()
	cleanup := initTracing()
	defer cleanup()

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}

// HandleErr is a generic error handler
func HandleErr(err error, message string) {
	if err != nil {
		log.Fatal().Err(err).Msg(message)
	}
}
