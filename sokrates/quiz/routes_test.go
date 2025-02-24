package quiz

import (
	"bytes"
	"encoding/json"
	"fmt"
	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	dialogueModel = `{
  "took": 0,
  "timed_out": false,
  "_shards": {
    "total": 2,
    "successful": 2,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 3,
      "relation": "eq"
    },
    "max_score": 6.2132735,
    "hits": [
      {
        "_index": "quiz",
        "_id": "ZnugVY0BlxdSsJ_QaT5p",
        "_score": 6.2132735,
        "_ignored": [
          "dialogue.introduction.keyword",
          "content.translation.keyword",
          "content.greek.keyword"
        ],
        "_source": {
          "quizMetadata": {
            "language": "English"
          },
          "theme": "Euripides - Medea",
          "quizType": "dialogue",
          "set": 1,
          "dialogue": {
            "introduction": "Euripides' 'Medea,' an ancient Greek tragedy written in 431 BCE, stands as a seminal work in the canon of classical literature. This play delves into the complexities of passion, betrayal, and revenge through its central character, Medea, a woman of formidable intelligence and emotion. Euripides challenges the societal norms of his time by portraying a female protagonist who defies the submissive role traditionally assigned to women in Greek society.  The narrative unfolds in Corinth, where Medea, a foreigner and a sorceress, grapples with the betrayal of her husband, Jason. Having forsaken her homeland and committed unspeakable acts for his sake, Medea is devastated when Jason abandons her and their children to marry Glauce, the daughter of Creon, the Corinthian king. The play examines Medea’s psychological turmoil as she oscillates between love, hatred, and the pursuit of justice.  Euripides' masterful use of dramatic tension, combined with his exploration of themes such as the plight of the outsider, the consequences of hubris, and the complexities of the human psyche, make 'Medea' a timeless tragedy. It not only reflects the mores of ancient Greek society but also resonates with contemporary audiences, inviting reflection on the nature of justice, loyalty, and the human condition.",
            "speakers": [
              {
                "name": "ΚΡΕΩΝ",
                "shorthand": "Κρέων",
                "translation": "Creon"
              },
              {
                "name": "ΜΗΔΕΙΑ",
                "shorthand": "Μήδεια",
                "translation": "Medea"
              }
            ],
            "section": "315-356",
            "linkToPerseus": "https://scaife.perseus.org/reader/urn:cts:greekLit:tlg0006.tlg003.perseus-grc2:300-360?right=perseus-eng2"
          },
          "content": [
            {
              "translation": "Thy words are soft to hear, but much I dread lest thou art devising some mischief in thy heart, and less than ever do I trust thee now; for a cunning woman, and man likewise,is easier to guard against when quick-tempered than when taciturn. Nay, begone at once! speak me no speeches, for this is decreed, nor hast thou any art whereby thou shalt abide amongst us, since thou hatest me.",
              "greek": "λέγεις ἀκοῦσαι μαλθάκ’, ἀλλ’ ἔσω φρενῶν ὀρρωδία μοι μή τι βουλεύσῃς κακόν, τόσῳ δέ γ’ ἧσσον ἢ πάρος πέποιθά σοι· γυνὴ γὰρ ὀξύθυμος, ὡς δ’ αὔτως ἀνήρ, ῥᾴων φυλάσσειν ἢ σιωπηλὸς σοφός.ἀλλ’ ἔξιθ’ ὡς τάχιστα, μὴ λόγους λέγε· ὡς ταῦτ’ ἄραρε, κοὐκ ἔχεις τέχνην ὅπως μενεῖς παρ’ ἡμῖν οὖσα δυσμενὴς ἐμοί.",
              "place": 1,
              "speaker": "Κρέων"
            },
            {
              "translation": "O, say not so! by thy knees and by thy daughter newly-wed, I do implore!",
              "greek": "μή, πρός σε γονάτων τῆς τε νεογάμου κόρης.",
              "place": 2,
              "speaker": "Μήδεια"
            },
            {
              "translation": "Thou wastest words; thou wilt never persuade me.",
              "greek": "λόγους ἀναλοῖς· οὐ γὰρ ἂν πείσαις ποτέ.",
              "place": 3,
              "speaker": "Κρέων"
            },
            {
              "translation": "What, wilt thou banish me, and to my prayers no pity yield?",
              "greek": "ἀλλ’ ἐξελᾷς με κοὐδὲν αἰδέσῃ λιτάς;",
              "place": 4,
              "speaker": "Μήδεια"
            },
            {
              "translation": "I will, for I love not thee above my own family.",
              "greek": "φιλῶ γὰρ οὐ σὲ μᾶλλον ἢ δόμους ἐμούς.",
              "place": 5,
              "speaker": "Κρέων"
            },
            {
              "translation": "O my country! what fond memories I have of thee in this hour!",
              "greek": "ὦ πατρίς, ὥς σου κάρτα νῦν μνείαν ἔχω.",
              "place": 6,
              "speaker": "Μήδεια"
            },
            {
              "translation": "Yea, for I myself love my city best of all things save my children.",
              "greek": "πλὴν γὰρ τέκνων ἔμοιγε φίλτατον πολύ.",
              "place": 7,
              "speaker": "Κρέων"
            },
            {
              "translation": "Ah me! ah me! to mortal man how dread a scourge is love!",
              "greek": "φεῦ φεῦ, βροτοῖς ἔρωτες ὡς κακὸν μέγα.",
              "place": 8,
              "speaker": "Μήδεια"
            },
            {
              "translation": "That, I deem, is according to the turn our fortunes take.",
              "greek": "ὅπως ἄν, οἶμαι, καὶ παραστῶσιν τύχαι.",
              "place": 9,
              "speaker": "Κρέων"
            },
            {
              "translation": "O Zeus! let not the author of these my troubles escape thee.",
              "greek": "Ζεῦ, μὴ λάθοι σε τῶνδ’ ὃς αἴτιος κακῶν.",
              "place": 10,
              "speaker": "Μήδεια"
            },
            {
              "translation": "Begone, thou silly woman, and free me from my toil.",
              "greek": "ἕρπ’, ὦ ματαία, καί μ’ ἀπάλλαξον πόνων.",
              "place": 11,
              "speaker": "Κρέων"
            },
            {
              "translation": "The toil is mine, no lack of it.",
              "greek": "πονοῦμεν ἡμεῖς κοὐ πόνων κεχρήμεθα.",
              "place": 12,
              "speaker": "Μήδεια"
            },
            {
              "translation": "Soon wilt thou be thrust out forcibly by the hand of servants.",
              "greek": "τάχ’ ἐξ ὀπαδῶν χειρὸς ὠσθήσῃ βίᾳ.",
              "place": 13,
              "speaker": "Κρέων"
            },
            {
              "translation": "Not that, not that, I do entreat thee, Creon!",
              "greek": "μὴ δῆτα τοῦτό γ’, ἀλλά σ’ αἰτοῦμαι, Κρέον . . .",
              "place": 14,
              "speaker": "Μήδεια"
            },
            {
              "translation": "Thou wilt cause disturbance yet, it seems.",
              "greek": "ὄχλον παρέξεις, ὡς ἔοικας, ὦ γύναι.",
              "place": 15,
              "speaker": "Κρέων"
            },
            {
              "translation": "I will begone; I ask thee not this boon to grant.",
              "greek": "φευξούμεθ’· οὐ τοῦθ’ ἱκέτευσα σοῦ τυχεῖν.",
              "place": 16,
              "speaker": "Μήδεια"
            },
            {
              "translation": "Why then this violence? why dost thou not depart?",
              "greek": "τί δαὶ βιάζῃ κοὐκ ἀπαλλάσσῃ χερός;",
              "place": 17,
              "speaker": "Κρέων"
            },
            {
              "translation": "Suffer me to abide this single day and devise some plan for the manner of my exile, and means of living for my children, since their father cares not to provide his babes therewith. Then pity them; thou too hast children of thine own; thou needs must have a kindly heart. For my own lot I care naught, though I an exile am, but for those babes I weep, that they should learn what sorrow means.",
              "greek": "μίαν με μεῖναι τήνδ’ ἔασον ἡμέραν καὶ ξυμπερᾶναι φροντίδ’ ᾗ φευξούμεθα, παισίν τ’ ἀφορμὴν τοῖς ἐμοῖς, ἐπεὶ πατὴρ οὐδὲν προτιμᾷ μηχανήσασθαι τέκνοις. οἴκτιρε δ’ αὐτούς· καὶ σύ τοι παίδων πατὴρ πέφυκας· εἰκὸς δ’ ἐστὶν εὔνοιάν σ’ ἔχειν. τοὐμοῦ γὰρ οὔ μοι φροντίς, εἰ φευξούμεθα, κείνους δὲ κλαίω συμφορᾷ κεχρημένους.",
              "place": 18,
              "speaker": "Μήδεια"
            },
            {
              "translation": "Mine is a nature anything but harsh; full oft by showing pity have I suffered shipwreck; and now albeit I clearly see my error, yet shalt thou gain this request, lady; but I do forewarn thee, if to-morrow’s rising sun shall find thee and thy children within the borders of this land, thou diest; my word is spoken and it will not lie. So now, if abide thou must, stay this one day only, for in it thou canst not do any of the fearful deeds I dread.",
              "greek": "ἥκιστα τοὐμὸν λῆμ’ ἔφυ τυραννικόν, αἰδούμενος δὲ πολλὰ δὴ διέφθορα· καὶ νῦν ὁρῶ μὲν ἐξαμαρτάνων, γύναι, ὅμως δὲ τεύξῃ τοῦδε· προυννέπω δέ σοι, εἴ σ’ ἡ ’πιοῦσα λαμπὰς ὄψεται θεοῦ καὶ παῖδας ἐντὸς τῆσδε τερμόνων χθονός, θανῇ· λέλεκται μῦθος ἀψευδὴς ὅδε. νῦν δ’, εἰ μένειν δεῖ, μίμν’ ἐφ’ ἡμέραν μίαν· οὐ γάρ τι δράσεις δεινὸν ὧν φόβος μ’ ἔχει.",
              "place": 19,
              "speaker": "Κρέων"
            }
          ]
        }
      }
    ]
  }
}`
	authorModel = `{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 2,
    "successful": 2,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 1515,
      "relation": "eq"
    },
    "max_score": 0.004683849,
    "hits": [
      {
        "_index": "quiz",
        "_id": "cnugVY0BlxdSsJ_QWzjI",
        "_score": 0.004683849,
        "_source": {
          "quizMetadata": {
            "language": "English"
          },
          "quizType": "authorbased",
          "theme": "Aeschylus - Agamemnon",
          "set": 3,
          "content": [
            {
              "translation": "by which way",
              "greek": "ὅπη"
            },
            {
              "translation": "auris, the ear",
              "greek": "οὖς"
            },
            {
              "translation": "somehow, in some way",
              "greek": "πως"
            },
            {
              "translation": "how? in what way",
              "greek": "πῶς"
            },
            {
              "translation": "a river, stream, flood",
              "greek": "ῥοή"
            },
            {
              "translation": "a moth",
              "greek": "σής"
            },
            {
              "translation": "your",
              "greek": "σός"
            },
            {
              "translation": "safe and sound, alive and well, in good case",
              "greek": "σῶς"
            },
            {
              "translation": "why? wherefore?",
              "greek": "τίη"
            },
            {
              "translation": "any one, any thing, some one, some thing;",
              "greek": "τις"
            },
            {
              "translation": "who? which?",
              "greek": "τίς"
            },
            {
              "translation": "so, in this wise",
              "greek": "τώς"
            },
            {
              "translation": "wood, material",
              "greek": "ὕλη"
            },
            {
              "translation": "a web",
              "greek": "ὑφή"
            },
            {
              "translation": "a man",
              "greek": "φώς"
            },
            {
              "translation": "it is fated, necessary",
              "greek": "χρή"
            },
            {
              "translation": "[sacrificial victim]",
              "greek": "ὥρα"
            },
            {
              "translation": "pollution, expiation",
              "greek": "ἄγος"
            },
            {
              "translation": "a gathering; a contest, a struggle, a trial",
              "greek": "ἀγών"
            },
            {
              "translation": "a blast, gale",
              "greek": "ἄημα"
            }
          ],
          "progress": {
            "timesCorrect": 0,
            "timesIncorrect": 0,
            "averageAccuracy": 0
          }
        }
      }
    ]
  }
}`
	mediaModel = `{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 2,
    "successful": 2,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 2,
      "relation": "eq"
    },
    "max_score": 5.702448,
    "hits": [
      {
        "_index": "quiz",
        "_id": "aHugVY0BlxdSsJ_QaT5t",
        "_score": 5.702448,
        "_source": {
          "quizMetadata": {
            "language": "English"
          },
          "quizType": "media",
          "set": 2,
          "content": [
            {
              "greek": "ἡ ἀγορά",
              "translation": "market place",
              "imageURL": "agora.webp",
              "audioFile": ""
            },
            {
              "greek": "ὁ οἶκος",
              "translation": "house",
              "imageURL": "house.webp",
              "audioFile": ""
            },
            {
              "greek": "ἡ ὁδός",
              "translation": "road",
              "imageURL": "road.webp",
              "audioFile": ""
            },
            {
              "greek": "ἡ πύλη",
              "translation": "gate",
              "imageURL": "gate.webp",
              "audioFile": ""
            },
            {
              "greek": "ὁ δικαστής",
              "translation": "judge",
              "imageURL": "judge.webp",
              "audioFile": ""
            },
            {
              "greek": "ὁ νόμος",
              "translation": "law",
              "imageURL": "law.webp",
              "audioFile": ""
            },
            {
              "greek": "ἡ δημοκρατία",
              "translation": "democracy",
              "imageURL": "democracy.webp",
              "audioFile": ""
            },
            {
              "greek": "ή πόλις",
              "translation": "city",
              "imageURL": "city.webp",
              "audioFile": ""
            },
            {
              "greek": "ἄριστος",
              "translation": "best",
              "imageURL": "best.webp",
              "audioFile": ""
            },
            {
              "greek": "πολύς",
              "translation": "many",
              "imageURL": "many.webp",
              "audioFile": ""
            },
            {
              "greek": "φιλόσοφος",
              "translation": "philosopher",
              "imageURL": "philosopher.webp",
              "audioFile": ""
            },
            {
              "greek": "ὁ ποιητής",
              "translation": "poet",
              "imageURL": "poet.webp",
              "audioFile": ""
            },
            {
              "greek": "ὁ βασιλεύς",
              "translation": "king",
              "imageURL": "king.webp",
              "audioFile": ""
            },
            {
              "greek": "ή βασιλεία",
              "translation": "queen",
              "imageURL": "queen.webp",
              "audioFile": ""
            },
            {
              "greek": "τό τεῖχος",
              "translation": "wall",
              "imageURL": "wall.webp",
              "audioFile": ""
            },
            {
              "greek": "τό ἱερόν",
              "translation": "temple",
              "imageURL": "temple.webp",
              "audioFile": ""
            },
            {
              "greek": "ὁ ῥήτωρ",
              "translation": "orator",
              "imageURL": "orator.webp",
              "audioFile": ""
            },
            {
              "greek": "ὁ τύραννος",
              "translation": "tyrant",
              "imageURL": "tyrant.webp",
              "audioFile": ""
            },
            {
              "greek": "τό θέατρον",
              "translation": "theatre",
              "imageURL": "theatre.webp",
              "audioFile": ""
            },
            {
              "greek": "τό στάδιον",
              "translation": "stadium",
              "imageURL": "stadium.webp",
              "audioFile": ""
            }
          ]
        }
      }
    ]
  }
}`
)

func TestHealthEndPoint(t *testing.T) {
	ticker := time.NewTicker(30 * time.Second)
	quizAttempts := make(chan models.QuizAttempt)
	aggregatedResult := make(map[string]models.QuizAttempt)

	t.Run("Pass", func(t *testing.T) {
		fixtureFile := "info"
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
		}

		router := InitRoutes(&testConfig)
		response := performGetRequest(router, "/sokrates/v1/health")

		var healthModel models.Health
		err = json.NewDecoder(response.Body).Decode(&healthModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, healthModel.Healthy)
	})

	t.Run("Fail", func(t *testing.T) {
		fixtureFile := "infoServiceDown"
		mockCode := 502
		mockElasticClient, err := elastic.NewMockClient(fixtureFile, mockCode)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
		}

		router := InitRoutes(&testConfig)
		response := performGetRequest(router, "/sokrates/v1/health")

		var healthModel models.Health
		err = json.NewDecoder(response.Body).Decode(&healthModel)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadGateway, response.Code)
		assert.False(t, healthModel.Healthy)
	})
}

func TestRandomizerOutput(t *testing.T) {
	randomizer, err := config.CreateNewRandomizer()
	assert.Nil(t, err)
	for i := range 100 {
		randomNumber := randomizer.RandomNumberBaseZero(20)
		fmt.Println(fmt.Sprintf("Random Number %d: %d", i, randomNumber))
	}
}

func TestCreateQuizEndpoint(t *testing.T) {
	ticker := time.NewTicker(1 * time.Hour)
	quizAttempts := make(chan models.QuizAttempt)
	aggregatedResult := make(map[string]models.QuizAttempt)
	randomizer, err := config.CreateNewRandomizer()
	assert.Nil(t, err)

	t.Run("Dialogue", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(dialogueModel)}, mockCode)
		assert.Nil(t, err)

		request := models.CreationRequest{
			Theme:    "sometheme",
			Set:      "1",
			QuizType: models.DIALOGUE,
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/sokrates/v1/quiz/create", bodyInBytes)

		var dialogue models.DialogueQuiz
		err = json.NewDecoder(response.Body).Decode(&dialogue)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, len(dialogue.Content) > 1)
	})

	t.Run("Media", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(mediaModel)}, mockCode)
		assert.Nil(t, err)

		request := models.CreationRequest{
			Theme:    "sometheme",
			Set:      "1",
			QuizType: models.MEDIA,
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
			Randomizer:         randomizer,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/sokrates/v1/quiz/create", bodyInBytes)

		var media models.QuizResponse
		err = json.NewDecoder(response.Body).Decode(&media)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, len(media.Options) > 1)
	})

	t.Run("Author", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(authorModel)}, mockCode)
		assert.Nil(t, err)

		request := models.CreationRequest{
			Theme:    "sometheme",
			Set:      "1",
			QuizType: models.AUTHORBASED,
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
			Randomizer:         randomizer,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/sokrates/v1/quiz/create", bodyInBytes)

		var author models.QuizResponse
		err = json.NewDecoder(response.Body).Decode(&author)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, len(author.Options) > 1)
	})
}

func TestAnswerQuiz(t *testing.T) {
	ticker := time.NewTicker(1 * time.Hour)
	quizAttempts := make(chan models.QuizAttempt)
	aggregatedResult := make(map[string]models.QuizAttempt)
	randomizer, err := config.CreateNewRandomizer()
	assert.Nil(t, err)

	t.Run("Dialogue", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(dialogueModel)}, mockCode)
		assert.Nil(t, err)

		request := models.AnswerRequest{
			Theme:         "sometheme",
			Set:           "1",
			QuizType:      models.DIALOGUE,
			Comprehensive: false,
			Answer:        "hello",
			Dialogue: []models.DialogueContent{
				{
					Translation: "Thy words are soft to hear, but much I dread lest thou art devising some mischief in thy heart, and less than ever do I trust thee now; for a cunning woman, and man likewise,is easier to guard against when quick-tempered than when taciturn. Nay, begone at once! speak me no speeches, for this is decreed, nor hast thou any art whereby thou shalt abide amongst us, since thou hatest me.",
					Greek:       "λέγεις ἀκοῦσαι μαλθάκ’, ἀλλ’ ἔσω φρενῶν ὀρρωδία μοι μή τι βουλεύσῃς κακόν, τόσῳ δέ γ’ ἧσσον ἢ πάρος πέποιθά σοι· γυνὴ γὰρ ὀξύθυμος, ὡς δ’ αὔτως ἀνήρ, ῥᾴων φυλάσσειν ἢ σιωπηλὸς σοφός.ἀλλ’ ἔξιθ’ ὡς τάχιστα, μὴ λόγους λέγε· ὡς ταῦτ’ ἄραρε, κοὐκ ἔχεις τέχνην ὅπως μενεῖς παρ’ ἡμῖν οὖσα δυσμενὴς ἐμοί.",
					Place:       1,
					Speaker:     "Κρέων",
				},
				{
					Translation: "O, say not so! by thy knees and by thy daughter newly-wed, I do implore!",
					Greek:       "μή, πρός σε γονάτων τῆς τε νεογάμου κόρης.",
					Place:       2,
					Speaker:     "Μήδεια",
				},
			},
			QuizWord: "",
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/sokrates/v1/quiz/answer", bodyInBytes)

		var dialogue models.DialogueAnswer
		err = json.NewDecoder(response.Body).Decode(&dialogue)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, dialogue.Percentage, 100.00)
	})

	t.Run("Media", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(mediaModel)}, mockCode)
		assert.Nil(t, err)

		request := models.AnswerRequest{
			Theme:         "sometheme",
			Set:           "1",
			QuizType:      models.MEDIA,
			Comprehensive: false,
			Answer:        "democracy",
			Dialogue:      nil,
			QuizWord:      "ἡ δημοκρατία",
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
			Randomizer:         randomizer,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/sokrates/v1/quiz/answer", bodyInBytes)

		var media models.ComprehensiveResponse
		err = json.NewDecoder(response.Body).Decode(&media)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, media.Correct)
	})

	t.Run("Author", func(t *testing.T) {
		mockCode := 200
		mockElasticClient, err := elastic.NewMockClient([][]byte{[]byte(authorModel)}, mockCode)
		assert.Nil(t, err)

		request := models.AnswerRequest{
			Theme:         "sometheme",
			Set:           "1",
			QuizType:      models.MEDIA,
			Comprehensive: false,
			Answer:        "a man",
			Dialogue:      nil,
			QuizWord:      "φώς",
		}

		jsonBody, err := json.Marshal(request)
		bodyInBytes := bytes.NewReader(jsonBody)
		assert.Nil(t, err)

		testConfig := SokratesHandler{
			Elastic:            mockElasticClient,
			QuizAttempts:       quizAttempts,
			AggregatedAttempts: aggregatedResult,
			Ticker:             ticker,
			Randomizer:         randomizer,
		}

		router := InitRoutes(&testConfig)
		response := performPostRequest(router, "/sokrates/v1/quiz/answer", bodyInBytes)

		var author models.ComprehensiveResponse
		err = json.NewDecoder(response.Body).Decode(&author)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, author.Correct)
	})
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
