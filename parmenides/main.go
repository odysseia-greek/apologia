package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/apologia/parmenides/aletheia"
	pb "github.com/odysseia-greek/delphi/aristides/proto"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

//go:embed sullego
var sullego embed.FS

func main() {
	logging.System(`
 ____   ____  ____   ___ ___    ___  ____   ____  ___      ___  _____
|    \ /    ||    \ |   |   |  /  _]|    \ |    ||   \    /  _]/ ___/
|  o  )  o  ||  D  )| _   _ | /  [_ |  _  | |  | |    \  /  [_(   \_ 
|   _/|     ||    / |  \_/  ||    _]|  |  | |  | |  D  ||    _]\__  |
|  |  |  _  ||    \ |   |   ||   [_ |  |  | |  | |     ||   [_ /  \ |
|  |  |  |  ||  .  \|   |   ||     ||  |  | |  | |     ||     |\    |
|__|  |__|__||__|\_||___|___||_____||__|__||____||_____||_____| \___|
                                                                     
`)

	logging.System(strings.Repeat("~", 37))
	logging.System("\"τό γάρ αυτο νοειν έστιν τε καί ειναι\"")
	logging.System("\"for it is the same thinking and being\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	handler, err := aletheia.CreateNewConfig()
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		log.Fatal(err)
	}
	err = handler.CreateIndexAtStartup()
	if err != nil {
		log.Fatal(err)
	}

	root := "sullego"

	// **Derive Directory Name from Index**
	quizDirName := stripQuizSuffix(handler.Index)

	// Check if the directory exists
	rootDir, err := sullego.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	var matchingDir string
	for _, dir := range rootDir {
		if dir.IsDir() && dir.Name() == quizDirName {
			matchingDir = dir.Name()
			break
		}
	}

	if matchingDir == "" {
		logging.Warn(fmt.Sprintf("No matching directory found for index: %s", handler.Index))
		return
	}

	// Process only the matching directory
	typePath := path.Join(root, matchingDir)
	typeDir, err := sullego.ReadDir(typePath)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	documentCounts := make(chan int, len(typeDir))
	documents := 0

	for _, quizFile := range typeDir {
		if quizFile.IsDir() {
			continue
		}

		quizPath := path.Join(typePath, quizFile.Name())
		content, err := sullego.ReadFile(quizPath)
		if err != nil {
			logging.Error(fmt.Sprintf("Failed to read file %s: %s", quizPath, err.Error()))
			continue
		}

		logging.Debug(fmt.Sprintf("Processing file: %s for index: %s", quizPath, handler.Index))

		wg.Add(1)
		go func(content []byte, quizPath string) {
			defer wg.Done()

			documentCount := 0
			switch handler.Index {
			case "media-quiz":
				documentCount = processQuizFile[models.MediaQuiz](content, handler, true) // Queue this
			case "dialogue-quiz":
				documentCount = processQuizFile[models.DialogueQuiz](content, handler, false) // No queue for dialogue
			case "author-based-quiz":
				documentCount = processQuizFile[models.AuthorbasedQuiz](content, handler, true) // Queue this
			case "multiple-choice-quiz":
				documentCount = processQuizFile[models.MultipleChoiceQuiz](content, handler, true) // Queue this
			case "grammar-quiz":
				documentCount = processQuizFile[aletheia.GrammarBasedQuiz](content, handler, true)
			case "journey-quiz":
				documentCount = processQuizFile[aletheia.JourneyBasedQuiz](content, handler, false) // No queue for journey mode
			}

			logging.Info(fmt.Sprintf("Documents found in %s: %d", quizPath, documentCount))
			documentCounts <- documentCount
		}(content, quizPath)
	}

	wg.Wait()
	close(documentCounts)
	for documentCount := range documentCounts {
		documents += documentCount
	}
	logging.Info(fmt.Sprintf("Created: %d documents", handler.Created))
	logging.Info(fmt.Sprintf("Words found in sullego: %d", documents))

	logging.Debug("Closing aristides because job is done")
	uuidCode := uuid.New().String()
	_, err = handler.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
	if err != nil {
		logging.Error(err.Error())
	}

	os.Exit(0)
}

// stripQuizSuffix removes '-quiz' and replaces hyphens with an empty string
func stripQuizSuffix(indexName string) string {
	indexName = strings.TrimSuffix(indexName, "-quiz") // Remove '-quiz'
	re := regexp.MustCompile(`-`)                      // Match hyphens
	return re.ReplaceAllString(indexName, "")          // Remove hyphens
}

func processQuizFile[T any](content []byte, handler *aletheia.ParmenidesHandler, useQueue bool) int {
	var quizzes []T
	if err := json.Unmarshal(content, &quizzes); err != nil {
		logging.Error("Failed to unmarshal JSON: " + err.Error())
		return 0
	}

	quizInterfaces := make([]interface{}, len(quizzes))
	for i, q := range quizzes {
		quizInterfaces[i] = q
	}

	if useQueue {
		// Batch all quizzes together using the queue
		if err := handler.AddWithQueue(quizInterfaces); err != nil {
			logging.Error(err.Error())
		}
	} else {
		if err := handler.AddWithoutQueue(quizInterfaces); err != nil {
			logging.Error(err.Error())
		}
	}

	return len(quizzes)
}
