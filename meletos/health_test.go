package meletos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/odysseia-greek/apologia/meletos/model/queries/health"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"
)

const (
	ResponseBody = "responseBody"
)

func (m *MeletosFixture) theGraphqlBackendIsRunning() error {
	const maxAttempts = 5
	const delayBetweenAttempts = 500 * time.Millisecond

	url := m.sokrates.baseUrl + "/sokrates/v1/ping"
	client := &http.Client{}

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("failed to call ping endpoint: %w", err)
		} else {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
			} else {
				var pingResponse struct {
					Healthy bool `json:"healthy"`
				}

				err = json.NewDecoder(resp.Body).Decode(&pingResponse)
				if err != nil {
					lastErr = fmt.Errorf("failed to decode ping response: %w", err)
				} else if pingResponse.Healthy {
					// success
					return nil
				} else {
					lastErr = fmt.Errorf("ping response is not healthy")
				}
			}
		}

		// If we failed, wait before retrying
		time.Sleep(delayBetweenAttempts)
	}

	// If all attempts fail, return the last error
	return fmt.Errorf("backend not healthy after %d attempts: %w", maxAttempts, lastErr)
}

func (m *MeletosFixture) iQueryTheHealthStatus() error {
	query := health.Health()
	resp, err := m.ForwardGraphql(query, map[string]interface{}{})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var healthResponse struct {
		Data struct {
			Response model.AggregatedHealthResponse `json:"health"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&healthResponse)

	if err != nil {
		return err
	}
	m.ctx = context.WithValue(m.ctx, ResponseBody, healthResponse.Data.Response)

	return nil
}

func (m *MeletosFixture) basicDatabaseHealthInfoShouldBeAvailableFor(service string) error {
	result := m.ctx.Value(ResponseBody).(model.AggregatedHealthResponse)

	for _, serviceHealth := range result.Services {
		if *serviceHealth.Name == service {

			dbInfo := serviceHealth.DatabaseInfo

			if !*dbInfo.Healthy {
				return fmt.Errorf("database is not healthy")
			}

			if *dbInfo.ClusterName == "" {
				return fmt.Errorf("no valid cluster name returned")
			}

			return nil
		}
	}

	return fmt.Errorf("service %s not found", service)
}

func (m *MeletosFixture) theVersionInformationShouldBeAvailableFor(service string) error {
	result := m.ctx.Value(ResponseBody).(model.AggregatedHealthResponse)

	for _, serviceHealth := range result.Services {
		if *serviceHealth.Name == service {
			if *serviceHealth.Version == "" {
				return fmt.Errorf("no valid version returned")
			}

			return nil
		}

	}

	return fmt.Errorf("service %s not found", service)
}

func (m *MeletosFixture) theServiceShouldBeHealthy(service string) error {
	result := m.ctx.Value(ResponseBody).(model.AggregatedHealthResponse)

	for _, serviceHealth := range result.Services {
		if *serviceHealth.Name == service {
			err := assertTrue(
				assert.True, *serviceHealth.Healthy,
				"expected service %s to be healthy but was: %v", service, *result.Healthy,
			)

			return err
		}
	}

	return fmt.Errorf("service %s not found", service)
}
