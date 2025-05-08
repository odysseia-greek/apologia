package meletos

import (
	"context"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/randomizer"
)

type MeletosFixture struct {
	Randomizer randomizer.Random
	ctx        context.Context
	sokrates   Sokrates
}

type Sokrates struct {
	graphql string
	baseUrl string
}

func New() (*MeletosFixture, error) {
	var gqlEndpoint string

	sokratesBaseUrl := config.StringFromEnv("SOKRATES_SERVICE", "")

	if sokratesBaseUrl == "" {
		gqlEndpoint = "http://k3d-odysseia.greek:8080/sokrates/graphql"
	} else {
		gqlEndpoint = sokratesBaseUrl + "/sokrates/graphql"
	}

	randomizer, err := config.CreateNewRandomizer()
	if err != nil {
		return nil, err
	}

	return &MeletosFixture{
		Randomizer: randomizer,
		sokrates: Sokrates{
			graphql: gqlEndpoint,
			baseUrl: sokratesBaseUrl,
		},
		ctx: context.Background(),
	}, nil
}
