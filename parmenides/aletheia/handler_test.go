package aletheia

import (
	"context"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParmenidesHandlerAdd(t *testing.T) {
	index := "test"
	body := models.Logos{Logos: []models.Word{{
		Method:      "",
		Category:    "",
		Greek:       "ἀγγέλλω",
		Translation: "to bear a message ",
		Chapter:     0,
	},
	},
	}

	channel := "testchannel"
	mockClient := &MockEupalinosClient{}

	t.Run("EnqueueTask", func(t *testing.T) {
		testHandler := ParmenidesHandler{
			Index:     index,
			Created:   0,
			Channel:   channel,
			Eupalinos: mockClient,
		}

		message, err := body.Marshal()
		assert.Nil(t, err)
		msg := &pb.Epistello{
			Data:    string(message),
			Channel: channel,
		}

		err = testHandler.EnqueueTask(context.Background(), msg)
		assert.Nil(t, err)
	})

}

func TestHandlerDeleteIndex(t *testing.T) {
	index := "test"

	t.Run("Deleted", func(t *testing.T) {
		file := "deleteIndex"
		status := 201
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := ParmenidesHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("IndexDoesNotExist", func(t *testing.T) {
		file := "deleteIndex404"
		status := 404
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := ParmenidesHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}
		err = testHandler.DeleteIndexAtStartUp()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := ParmenidesHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		err = testHandler.DeleteIndexAtStartUp()
		assert.NotNil(t, err)
	})
}

func TestHandlerCreateIndex(t *testing.T) {
	index := "test"

	t.Run("Created", func(t *testing.T) {
		files := []string{"createIndex", "createIndex"}
		status := 201
		mockElasticClient, err := elastic.NewMockClient(files, status)
		assert.Nil(t, err)

		testHandler := ParmenidesHandler{
			Elastic:    mockElasticClient,
			Index:      index,
			Created:    0,
			PolicyName: fmt.Sprintf("%s-policy", index),
		}
		err = testHandler.CreateIndexAtStartup()
		assert.Nil(t, err)
	})

	t.Run("NotCreated", func(t *testing.T) {
		file := "error"
		status := 502
		mockElasticClient, err := elastic.NewMockClient(file, status)
		assert.Nil(t, err)

		testHandler := ParmenidesHandler{
			Elastic: mockElasticClient,
			Index:   index,
			Created: 0,
		}

		err = testHandler.CreateIndexAtStartup()
		assert.NotNil(t, err)
	})
}
