package crawlworker

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func getHTTPClient(dialcontext func(ctx context.Context, network, address string) (net.Conn, error)) *http.Client {
	transport := otelhttp.NewTransport(&http.Transport{
		Proxy:               nil,
		DialContext:         dialcontext,
		ForceAttemptHTTP2:   false,
		MaxIdleConns:        100, // Defaut
		MaxIdleConnsPerHost: 2,   // Default
		IdleConnTimeout:     90 * time.Second,
	})

	return &http.Client{
		Transport: transport,
	}
}
