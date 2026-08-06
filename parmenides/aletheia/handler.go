package aletheia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	pb "github.com/odysseia-greek/agora/eupalinos/v1"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	arisv1 "github.com/odysseia-greek/alexandreia/aristarchos/gen/go/v1"
	aristides "github.com/odysseia-greek/delphi/aristides/diplomat"
)

const requestTimeout = 30 * time.Second

type queueClient interface {
	EnqueueMessage(context.Context, *pb.Epistello) (*pb.EnqueueResponse, error)
}

type ParmenidesHandler struct {
	Index            string
	Created          int
	Elastic          elastic.Client
	Eupalinos        queueClient
	Channel          string
	DutchChannel     string
	ExitCode         string
	PolicyName       string
	Ambassador       *aristides.ClientAmbassador
	Aggregator       arisv1.Aristarchos_CreateNewEntryClient
	AggregatorCancel context.CancelFunc
}

func (p *ParmenidesHandler) DeleteIndexAtStartUp(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	deleted, err := p.Elastic.Index().DeleteWithContext(ctx, p.Index)
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

func (p *ParmenidesHandler) CreateIndexAtStartup(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	indexMapping := quizIndex(p.PolicyName, 1, 0)
	created, err := p.Elastic.Index().CreateWithContext(ctx, p.Index, indexMapping)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", created.Index, created.Acknowledged))

	return nil
}

func (p *ParmenidesHandler) AddWithQueue(ctx context.Context, quizDocs []interface{}) error {
	var buf bytes.Buffer
	var wg sync.WaitGroup

	// Process each quiz type in a separate goroutine
	for _, doc := range quizDocs {
		switch q := doc.(type) {
		case models.MediaQuiz:
			wg.Add(1)
			go func(q models.MediaQuiz) {
				defer wg.Done()
				p.processMediaQuiz(ctx, q)
			}(q)

		case models.AuthorbasedQuiz:
			wg.Add(1)
			go func(q models.AuthorbasedQuiz) {
				defer wg.Done()
				p.processAuthorBasedQuiz(ctx, q)
			}(q)

		case models.MultipleChoiceQuiz:
			wg.Add(1)
			go func(q models.MultipleChoiceQuiz) {
				defer wg.Done()
				p.processMultipleChoiceQuiz(ctx, q)
			}(q)

		case GrammarBasedQuiz:
			wg.Add(1)
			go func(q GrammarBasedQuiz) {
				defer wg.Done()
				p.processGrammarQuiz(ctx, q)
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
	elasticCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	res, err := p.Elastic.Document().BulkWithContext(elasticCtx, buf, p.Index)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	p.Created += len(res.Items)
	return nil
}

func (p *ParmenidesHandler) AddWithoutQueue(ctx context.Context, quizDocs []interface{}) error {
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
	elasticCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	res, err := p.Elastic.Document().BulkWithContext(elasticCtx, buf, p.Index)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	p.Created += len(res.Items)
	return nil
}

func (p *ParmenidesHandler) processMediaQuiz(ctx context.Context, q models.MediaQuiz) {
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

		err := p.enqueueTask(ctx, msg)
		if err != nil {
			logging.Error(err.Error())
		}
	}
}

func (p *ParmenidesHandler) processAuthorBasedQuiz(ctx context.Context, q models.AuthorbasedQuiz) {
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

		err := p.enqueueTask(ctx, msg)
		if err != nil {
			logging.Error(err.Error())
		}

		if word.HasGrammarQuestions {
			for _, grammarQuestion := range word.GrammarQuestions {
				err = p.sendToAggregator(ctx, grammarQuestion, word.Greek, word.Translation)
				if err != nil {
					logging.Error(err.Error())
					break
				}
			}
		}

	}
}

func (p *ParmenidesHandler) processMultipleChoiceQuiz(ctx context.Context, q models.MultipleChoiceQuiz) {
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

		err := p.enqueueTask(ctx, msg)
		if err != nil {
			logging.Error(err.Error())
		}
	}
}

func (p *ParmenidesHandler) processGrammarQuiz(ctx context.Context, q GrammarBasedQuiz) {
	for _, word := range q.Content {
		grammarQuestion := models.GrammarQuestion{
			CorrectAnswer:    word.GrammarQuestion.CorrectAnswer,
			TypeOfWord:       "verb",
			WordInText:       word.Greek,
			ExtraInformation: "",
		}

		if strings.Contains(q.Theme, "Participle") {
			grammarQuestion.TypeOfWord = "participle"
		}

		err := p.sendToAggregator(ctx, grammarQuestion, word.DictionaryForm, word.Translation)
		if err != nil {
			logging.Error(err.Error())
		}
	}
}

func (p *ParmenidesHandler) sendToAggregator(ctx context.Context, grammarQuestion models.GrammarQuestion, greekWord, translation string) error {
	if p.Aggregator == nil {
		return fmt.Errorf("aggregator is empty")
	}

	traceID, err := uuid.NewUUID()
	ctx = context.WithValue(ctx, service.HeaderKey, traceID.String())
	if err != nil {
		return err
	}

	// send word to aggregator
	partOfSpeech := arisv1.PartOfSpeech_VERB
	if grammarQuestion.TypeOfWord == "noun" {
		partOfSpeech = arisv1.PartOfSpeech_NOUN
	} else if grammarQuestion.TypeOfWord == "misc" {
		partOfSpeech = arisv1.PartOfSpeech_PARTICIPLE
	} else if grammarQuestion.TypeOfWord == "verb" {
		if strings.Contains(grammarQuestion.CorrectAnswer, "part") {
			partOfSpeech = arisv1.PartOfSpeech_PARTICLE
		}
	}

	request := &arisv1.AggregatorCreationRequest{
		Word:         grammarQuestion.WordInText,
		Rule:         grammarQuestion.CorrectAnswer,
		RootWord:     greekWord,
		Translation:  translation,
		PartOfSpeech: partOfSpeech,
		TraceId:      traceID.String(),
	}

	if err = p.Aggregator.Send(request); err != nil {
		logging.Error(err.Error())
		return err
	}

	return nil
}

// EnqueueTask sends a task to the Eupalinos queue
func (p *ParmenidesHandler) enqueueTask(ctx context.Context, message *pb.Epistello) error {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

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
