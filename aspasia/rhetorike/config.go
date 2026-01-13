package rhetorike

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
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

	fuzzyClientAddress := config.StringFromEnv("ANTIGONOS_SERVICE", "antigonos.makedonia.svc.cluster.local:50060")
	fuzzyClient, err := hesiodos.NewGenericGrpcClient[*monophthalmus.FuzzyClient](
		fuzzyClientAddress,
		monophthalmus.NewAntigonosClient,
	)

	if err != nil {
		logging.Error(err.Error())
	}

	fuzzyClientHealthy := false
	if fuzzyClient != nil {
		fuzzyClientHealthy = fuzzyClient.client.WaitForHealthyState()
	}

	elapsed := time.Since(start)

	logging.System(fmt.Sprintf(`Aspasia Configuration Overview:
- Initialization Time: %s
- Antigonos Service:   %v (Address: %s)
`,
		elapsed,
		fuzzyClientHealthy, fuzzyClientAddress,
	))

	return &GathererServiceImpl{
		Archytas:    cache,
		Version:     version,
		Client:      client,
		FuzzyClient: fuzzyClient,
	}, nil
}
