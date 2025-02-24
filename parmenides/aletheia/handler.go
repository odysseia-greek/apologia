package aletheia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	ptolemaios "github.com/odysseia-greek/delphi/ptolemaios/diplomat"
	"strings"
)

type ParmenidesHandler struct {
	Index        string
	Created      int
	Elastic      elastic.Client
	Eupalinos    EupalinosClient
	Channel      string
	DutchChannel string
	ExitCode     string
	PolicyName   string
	Ambassador   *ptolemaios.ClientAmbassador
}

func (p *ParmenidesHandler) DeleteIndexAtStartUp() error {
	deleted, err := p.Elastic.Index().Delete(p.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", p.Index, deleted))
	if err != nil {
		if deleted {
			return nil
		}
		if strings.Contains(err.Error(), "index_not_found_exception") {
			logging.Error(err.Error())
			return nil
		}

		return err
	}

	return nil
}

func (p *ParmenidesHandler) createPolicyAtStartup() error {
	policyCreated, err := p.Elastic.Policy().CreateHotPolicy(p.PolicyName)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created policy: %s %v", p.PolicyName, policyCreated.Acknowledged))

	return nil
}

func (p *ParmenidesHandler) CreateIndexAtStartup() error {
	logging.Info(fmt.Sprintf("creating policy: %s", p.PolicyName))
	err := p.createPolicyAtStartup()
	if err != nil {
		return err

	}
	indexMapping := quizIndex(p.PolicyName)
	created, err := p.Elastic.Index().Create(p.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", created.Index, created.Acknowledged))

	return nil
}

func (p *ParmenidesHandler) AddWithQueue(quizzes []models.MultipleChoiceQuiz) error {
	var buf bytes.Buffer

	var currBatch int

	for _, quiz := range quizzes {
		currBatch++
		for _, word := range quiz.Content {
			meros := models.Meros{
				Greek:    word.Greek,
				English:  word.Translation,
				Original: word.Greek,
			}

			alternateChannel := false
			if quiz.QuizMetadata.Language == "Dutch" {
				meros.Dutch = word.Translation
				meros.English = ""
				alternateChannel = true
			}

			jsonsifiedMeros, _ := meros.Marshal()
			msg := &pb.Epistello{
				Data:    string(jsonsifiedMeros),
				Channel: p.Channel,
			}

			if alternateChannel {
				msg.Channel = p.DutchChannel
			}

			err := p.EnqueueTask(context.Background(), msg)
			if err != nil {
				logging.Error(err.Error())
			}
		}

		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedWord, _ := json.Marshal(quiz)
		jsonifiedWord = append(jsonifiedWord, "\n"...)
		buf.Grow(len(meta) + len(jsonifiedWord))
		buf.Write(meta)
		buf.Write(jsonifiedWord)

		if currBatch == len(quizzes) {
			res, err := p.Elastic.Document().Bulk(buf, p.Index)
			if err != nil {
				logging.Error(err.Error())
				return err
			}

			p.Created = p.Created + len(res.Items)
		}
	}
	return nil
}

func (p *ParmenidesHandler) AddWithoutQueue(content []byte) error {
	_, err := p.Elastic.Index().CreateDocument(p.Index, content)
	if err != nil {
		return err
	}

	p.Created += 1

	return nil
}

// EnqueueTask sends a task to the Eupalinos queue
func (p *ParmenidesHandler) EnqueueTask(ctx context.Context, message *pb.Epistello) error {
	traceID, err := uuid.NewUUID()
	ctx = context.WithValue(ctx, service.HeaderKey, traceID.String())

	_, err = p.Eupalinos.EnqueueMessage(ctx, message)
	return err
}
