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
	aristides "github.com/odysseia-greek/delphi/aristides/diplomat"
	pba "github.com/odysseia-greek/olympia/aristarchos/proto"
	"strings"
	"sync"
)

type ParmenidesHandler struct {
	Index            string
	Created          int
	Elastic          elastic.Client
	Eupalinos        EupalinosClient
	Channel          string
	DutchChannel     string
	ExitCode         string
	PolicyName       string
	Ambassador       *aristides.ClientAmbassador
	Aggregator       pba.Aristarchos_CreateNewEntryClient
	AggregatorCancel context.CancelFunc
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
	indexMapping := quizIndex(p.PolicyName, 1, 0)
	created, err := p.Elastic.Index().Create(p.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", created.Index, created.Acknowledged))

	return nil
}

func (p *ParmenidesHandler) AddWithQueue(quizDocs []interface{}) error {
	var buf bytes.Buffer
	var wg sync.WaitGroup

	// Process each quiz type in a separate goroutine
	for _, doc := range quizDocs {
		switch q := doc.(type) {
		case models.MediaQuiz:
			wg.Add(1)
			go func(q models.MediaQuiz) {
				defer wg.Done()
				p.processMediaQuiz(q)
			}(q)

		case models.AuthorbasedQuiz:
			wg.Add(1)
			go func(q models.AuthorbasedQuiz) {
				defer wg.Done()
				p.processAuthorBasedQuiz(q)
			}(q)

		case models.MultipleChoiceQuiz:
			wg.Add(1)
			go func(q models.MultipleChoiceQuiz) {
				defer wg.Done()
				p.processMultipleChoiceQuiz(q)
			}(q)
		}
	}

	// Collect documents for bulk indexing
	for _, doc := range quizDocs {
		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedQuiz, err := json.Marshal(doc)
		if err != nil {
			logging.Error("Failed to marshal quiz: " + err.Error())
			continue
		}

		buf.Grow(len(meta) + len(jsonifiedQuiz) + 1)
		buf.Write(meta)
		buf.Write(jsonifiedQuiz)
		buf.WriteByte('\n')
	}

	// Wait for queue processing before sending to Elasticsearch
	wg.Wait()

	// Bulk insert into Elasticsearch
	res, err := p.Elastic.Document().Bulk(buf, p.Index)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	p.Created += len(res.Items)
	return nil
}

func (p *ParmenidesHandler) AddWithoutQueue(quizDocs []interface{}) error {
	var buf bytes.Buffer

	for _, doc := range quizDocs {
		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedDoc, err := json.Marshal(doc)
		if err != nil {
			logging.Error("Failed to marshal JSON: " + err.Error())
			continue
		}

		buf.Grow(len(meta) + len(jsonifiedDoc) + 1)
		buf.Write(meta)
		buf.Write(jsonifiedDoc)
		buf.WriteByte('\n') // Ensure newline after each document
	}

	// Send all documents in a single bulk request
	res, err := p.Elastic.Document().Bulk(buf, p.Index)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	p.Created += len(res.Items)
	return nil
}

func (p *ParmenidesHandler) processMediaQuiz(q models.MediaQuiz) {
	for _, word := range q.Content {
		meros := models.Meros{
			Greek:    word.Greek,
			English:  word.Translation,
			Original: word.Greek,
		}

		jsonMeros, _ := meros.Marshal()
		msg := &pb.Epistello{
			Data:    string(jsonMeros),
			Channel: p.Channel,
		}

		err := p.enqueueTask(context.Background(), msg)
		if err != nil {
			logging.Error(err.Error())
		}
	}
}

func (p *ParmenidesHandler) processAuthorBasedQuiz(q models.AuthorbasedQuiz) {
	for _, word := range q.Content {
		meros := models.Meros{
			Greek:    word.Greek,
			English:  word.Translation,
			Original: word.Greek,
		}

		jsonMeros, _ := meros.Marshal()
		msg := &pb.Epistello{
			Data:    string(jsonMeros),
			Channel: p.Channel,
		}

		err := p.enqueueTask(context.Background(), msg)
		if err != nil {
			logging.Error(err.Error())
		}

		if word.HasGrammarQuestions {
			for _, grammarQuestion := range word.GrammarQuestions {
				err = p.sendToAggregator(context.Background(), grammarQuestion, word.Greek, word.Translation)
				if err != nil {
					logging.Error(err.Error())
				}
			}
		}

	}
}

func (p *ParmenidesHandler) processMultipleChoiceQuiz(q models.MultipleChoiceQuiz) {
	for _, word := range q.Content {
		meros := models.Meros{
			Greek:    word.Greek,
			English:  word.Translation,
			Original: word.Greek,
		}

		alternateChannel := false
		if q.QuizMetadata.Language == "Dutch" {
			meros.Dutch = word.Translation
			meros.English = ""
			alternateChannel = true
		}

		jsonMeros, _ := meros.Marshal()
		msg := &pb.Epistello{
			Data:    string(jsonMeros),
			Channel: p.Channel,
		}

		if alternateChannel {
			msg.Channel = p.DutchChannel
		}

		err := p.enqueueTask(context.Background(), msg)
		if err != nil {
			logging.Error(err.Error())
		}
	}
}

func (p *ParmenidesHandler) sendToAggregator(ctx context.Context, grammarQuestion models.GrammarQuestion, greekWord, translation string) error {
	traceID, err := uuid.NewUUID()
	ctx = context.WithValue(ctx, service.HeaderKey, traceID.String())
	if err != nil {
		return err
	}

	// send word to aggregator
	partOfSpeech := pba.PartOfSpeech_VERB
	if grammarQuestion.TypeOfWord == "noun" {
		partOfSpeech = pba.PartOfSpeech_NOUN
	} else if grammarQuestion.TypeOfWord == "misc" {
		partOfSpeech = pba.PartOfSpeech_PARTICIPLE
	} else if grammarQuestion.TypeOfWord == "verb" {
		if strings.Contains(grammarQuestion.CorrectAnswer, "part") {
			partOfSpeech = pba.PartOfSpeech_PARTICLE
		}
	}

	request := &pba.AggregatorCreationRequest{
		Word:         grammarQuestion.WordInText,
		Rule:         grammarQuestion.CorrectAnswer,
		RootWord:     greekWord,
		Translation:  translation,
		PartOfSpeech: partOfSpeech,
		TraceId:      traceID.String(),
	}

	logging.Debug(fmt.Sprintf("sending to aggregator: %s", request.String()))
	if err = p.Aggregator.Send(request); err != nil {
		return err
	}
	return nil
}

// EnqueueTask sends a task to the Eupalinos queue
func (p *ParmenidesHandler) enqueueTask(ctx context.Context, message *pb.Epistello) error {
	traceID, err := uuid.NewUUID()
	ctx = context.WithValue(ctx, service.HeaderKey, traceID.String())
	if err != nil {
		return err
	}

	_, err = p.Eupalinos.EnqueueMessage(ctx, message)
	if err != nil {
		return err
	}
	return err
}
