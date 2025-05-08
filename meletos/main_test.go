package meletos

import (
	"embed"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/odysseia-greek/agora/plato/logging"
	"os"
	"strings"
	"testing"
)

const (
	GrammarQuiz = "grammarQuiz"
	Question    = "question"
	Variables   = "variables"
	Responses   = "responses"
	Progress    = "progress"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

//go:embed features/*.feature
var featureFiles embed.FS

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {

		//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=MELETOS
		logging.System(`
 ___ ___    ___  _        ___ ______   ___   _____
|   |   |  /  _]| |      /  _]      | /   \ / ___/
| _   _ | /  [_ | |     /  [_|      ||     (   \_ 
|  \_/  ||    _]| |___ |    _]_|  |_||  O  |\__  |
|   |   ||   [_ |     ||   [_  |  |  |     |/  \ |
|   |   ||     ||     ||     | |  |  |     |\    |
|___|___||_____||_____||_____| |__|   \___/  \___|
`)
		logging.System("\"πῶς λέγεις, ὦ Μέλητε; οἵδε τοὺς νέους παιδεύειν οἷοί τέ εἰσι καὶ βελτίους ποιοῦσιν;\"")
		logging.System("\"What are you saying, Meletus? Are these gentlemen able to instruct the youth, and do they make them better?\"")
		logging.System("starting test suite setup.....")

		logging.System("getting env variables and creating config")

	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	meletos, err := New()
	if err != nil {
		os.Exit(1)
	}

	//general
	ctx.Step(`^the graphql backend is running$`, meletos.theGraphqlBackendIsRunning)

	//health
	ctx.Step(`^I query the health status$`, meletos.iQueryTheHealthStatus)
	ctx.Step(`^the service should be healthy$`, meletos.theServiceShouldBeHealthy)
	ctx.Step(`^the service "([^"]*)" should be healthy$`, meletos.theServiceShouldBeHealthy)
	ctx.Step(`^basic database health info should be available for "([^"]*)"$`, meletos.basicDatabaseHealthInfoShouldBeAvailableFor)
	ctx.Step(`^the version information should be available for "([^"]*)"$`, meletos.theVersionInformationShouldBeAvailableFor)

	//media
	ctx.Step(`^I query for media quiz options$`, meletos.iQueryForMediaQuizOptions)
	ctx.Step(`^I submit each option once$`, meletos.iSubmitEachOptionOnce)
	ctx.Step(`^I use the media options to create a question$`, meletos.iUseTheMediaOptionsToCreateAQuestion)
	ctx.Step(`^I should have (\d+) incorrect and (\d+) correct answer$`, meletos.iShouldHaveIncorrectAndCorrectAnswer)
	ctx.Step(`^The Progress should be (\d+) incorrect and (\d+) correct answer$`, meletos.theProgressShouldBeIncorrectAndCorrectAnswer)

	//multiple choice
	ctx.Step(`^I query for multiple choice quiz options$`, meletos.iQueryForMultipleChoiceQuizOptions)
	ctx.Step(`^I use the multiple choice options to create a question$`, meletos.iUseTheMultipleChoiceOptionsToCreateAQuestion)
	ctx.Step(`^I submit each multiple choice option once$`, meletos.iSubmitEachMultipleChoiceOptionOnce)

	//grammnar
	ctx.Step(`^I query for grammar quiz options$`, meletos.iQueryForGrammarQuizOptions)
	ctx.Step(`^I use the grammar options to create a question$`, meletos.iUseTheGrammarOptionsToCreateAQuestion)
	ctx.Step(`^I submit each grammar option once$`, meletos.iSubmitEachGrammarOptionOnce)

	//authorbased
	ctx.Step(`^that question has a Greek, English sentence$`, meletos.thatQuestionHasAGreekEnglishSentence)
	ctx.Step(`^that question has a reference to the text module$`, meletos.thatQuestionHasAReferenceToTheTextModule)
	ctx.Step(`^I query for authorbased quiz options$`, meletos.iQueryForAuthorbasedQuizOptions)
	ctx.Step(`^I submit each authorbased option once$`, meletos.iSubmitEachAuthorbasedOptionOnce)
	ctx.Step(`^I use the authorbased options to create a question$`, meletos.iUseTheAuthorbasedOptionsToCreateAQuestion)
	ctx.Step(`^grammar options should be embedded into the quiz for some words$`, meletos.grammarOptionsShouldBeEmbeddedIntoTheQuizForSomeWords)
	ctx.Step(`^I create a quiz that has the name "([^"]*)"$`, meletos.iCreateAQuizThatHasTheName)
	ctx.Step(`^I query the word forms for a segment$`, meletos.iQueryTheWordFormsForASegment)
	ctx.Step(`^the words should be returned as they appear in the text$`, meletos.theWordsShouldBeReturnedAsTheyAppearInTheText)

	//journey
	ctx.Step(`^a new journey is returned with a translation and sentence$`, meletos.aNewJourneyIsReturnedWithATranslationAndSentence)
	ctx.Step(`^a short background on the text should exist$`, meletos.aShortBackgroundOnTheTextShouldExist)
	ctx.Step(`^I query for journey quiz options$`, meletos.iQueryForJourneyQuizOptions)
	ctx.Step(`^I use the journey options to create a question$`, meletos.iUseTheJourneyOptionsToCreateAQuestion)
	ctx.Step(`^the quiz has different types of questions embedded$`, meletos.theQuizHasDifferentTypesOfQuestionsEmbedded)

	//dialogue
	ctx.Step(`^I query for dialogue quiz options$`, meletos.iQueryForDialogueQuizOptions)
	ctx.Step(`^I submit with a perfect input$`, meletos.iSubmitWithAPerfectInput)
	ctx.Step(`^I submit with at least one section wronly placed$`, meletos.iSubmitWithAtLeastOneSectionWronlyPlaced)
	ctx.Step(`^I use the dialogue options to create a question$`, meletos.iUseTheDialogueOptionsToCreateAQuestion)
	ctx.Step(`^the percentage should be (\d+)$`, meletos.thePercentageShouldBe)
	ctx.Step(`^the percentage should be lower than (\d+)$`, meletos.thePercentageShouldBeLowerThan)
	ctx.Step(`^wronglyPlaced should be empty$`, meletos.wronglyPlacedShouldBeEmpty)
	ctx.Step(`^wronglyPlaced should hold a reference to the correct place$`, meletos.wronglyPlacedShouldHoldAReferenceToTheCorrectPlace)
}

func TestMain(m *testing.M) {
	format := "pretty"
	var tag string // Initialize an empty slice to store the tags

	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" {
			format = "progress"
		} else if strings.HasPrefix(arg, "-tags=") {
			tagsString := strings.TrimPrefix(arg, "-tags=")
			tag = strings.Split(tagsString, ",")[0]
		}
	}

	opts := godog.Options{
		Format:          format,
		FeatureContents: getFeatureContents(), // Get the embedded feature files
	}

	if tag != "" {
		opts.Tags = tag
	}

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

func getFeatureContents() []godog.Feature {
	features := []godog.Feature{}
	featureFileNames, _ := featureFiles.ReadDir("features")
	for _, file := range featureFileNames {
		if !file.IsDir() && file.Name() != "README.md" { // Skip directories and README.md if any
			filePath := fmt.Sprintf("features/%s", file.Name())
			fileContent, _ := featureFiles.ReadFile(filePath)
			features = append(features, godog.Feature{Name: file.Name(), Contents: fileContent})
		}
	}
	return features
}
