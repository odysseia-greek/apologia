package rhetorike

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
)

func CreateNewConfig(ctx context.Context) (*GathererServiceImpl, error) {
	start := time.Now()

	cache, err := archytas.CreateBadgerClient()
	if err != nil {
		return nil, err
	}

	version := os.Getenv(config.EnvVersion)

	client, err := config.CreateOdysseiaClient()
	if err != nil {
		return nil, err
	}

	if err != nil {
		logging.Error(err.Error())
	}

	tracer, err := aristophanes.NewClientTracer(aristophanes.DefaultAddress)
	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	streamer, err := tracer.Chorus(ctx)
	alexandrosGraphQLEndpoint := config.StringFromEnv("ALEXANDROS_GATEWAY", "http://alexandros.makedonia.svc:8080/alexandros/graphq")
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second, // “hard cap” safety; ctx can be shorter
	}

	elapsed := time.Since(start)

	logging.System(fmt.Sprintf(`Aspasia Configuration Overview:
- Initialization Time: %s
- Alexandros Service:  %s
`,
		elapsed,
		alexandrosGraphQLEndpoint,
	))

	return &GathererServiceImpl{
		Archytas:          cache,
		Version:           version,
		Client:            client,
		Streamer:          streamer,
		GraphqlClient:     httpClient,
		AlexandrosAddress: alexandrosGraphQLEndpoint,
	}, nil
}
