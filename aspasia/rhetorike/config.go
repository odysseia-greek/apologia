package rhetorike

import (
	"context"
	"os"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/apologia/diotima/theoria"
)

func CreateNewConfig(ctx context.Context) (*ExtendedServiceImpl, error) {
	theoria.SetStreamer(ctx)

	cache, err := archytas.CreateBadgerClient()
	if err != nil {
		return nil, err
	}

	version := os.Getenv(config.EnvVersion)

	return &ExtendedServiceImpl{
		Archytas: cache,
		Version:  version,
	}, nil
}
