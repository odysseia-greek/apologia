// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AnalyzeResult struct {
	Author        *string `json:"author,omitempty"`
	Book          *string `json:"book,omitempty"`
	Reference     *string `json:"reference,omitempty"`
	ReferenceLink *string `json:"referenceLink,omitempty"`
	Text          *Rhema  `json:"text,omitempty"`
}

type AnalyzeTextResponse struct {
	Conjugations []*ConjugationResponse `json:"conjugations,omitempty"`
	Results      []*AnalyzeResult       `json:"results,omitempty"`
	Rootword     *string                `json:"rootword,omitempty"`
}

type AuthorBasedAnswerInput struct {
	Theme    *string `json:"theme,omitempty"`
	Set      *string `json:"set,omitempty"`
	Segment  *string `json:"segment,omitempty"`
	QuizWord *string `json:"quizWord,omitempty"`
	Answer   *string `json:"answer,omitempty"`
}

type AuthorBasedAnswerResponse struct {
	Correct     *bool     `json:"correct,omitempty"`
	QuizWord    *string   `json:"quizWord,omitempty"`
	WordsInText []*string `json:"wordsInText,omitempty"`
}

type AuthorBasedInput struct {
	ExcludeWords []*string `json:"excludeWords,omitempty"`
	Theme        *string   `json:"theme,omitempty"`
	Set          *string   `json:"set,omitempty"`
	Segment      *string   `json:"segment,omitempty"`
}

type AuthorBasedOptions struct {
	QuizWord *string `json:"quizWord,omitempty"`
}

type AuthorBasedQuiz struct {
	QuizItem      *string               `json:"quizItem,omitempty"`
	NumberOfItems *int32                `json:"numberOfItems,omitempty"`
	Options       []*AuthorBasedOptions `json:"options,omitempty"`
}

type AuthorBasedResponse struct {
	FullSentence *string             `json:"fullSentence,omitempty"`
	Translation  *string             `json:"translation,omitempty"`
	Reference    *string             `json:"reference,omitempty"`
	Quiz         *AuthorBasedQuiz    `json:"quiz,omitempty"`
	GrammarQuiz  []*GrammarQuizAdded `json:"grammarQuiz,omitempty"`
}

type ComprehensiveResponse struct {
	Correct      *bool                `json:"correct,omitempty"`
	FoundInText  *AnalyzeTextResponse `json:"foundInText,omitempty"`
	QuizWord     *string              `json:"quizWord,omitempty"`
	SimilarWords []*Hit               `json:"similarWords,omitempty"`
}

type ConjugationResponse struct {
	Rule *string `json:"rule,omitempty"`
	Word *string `json:"word,omitempty"`
}

type Dialogue struct {
	Introduction  *string    `json:"introduction,omitempty"`
	Speakers      []*Speaker `json:"speakers,omitempty"`
	Section       *string    `json:"section,omitempty"`
	LinkToPerseus *string    `json:"linkToPerseus,omitempty"`
}

type DialogueAnswer struct {
	Percentage    *float64              `json:"percentage,omitempty"`
	Input         []*DialogueContent    `json:"input,omitempty"`
	Answer        []*DialogueContent    `json:"answer,omitempty"`
	WronglyPlaced []*DialogueCorrection `json:"wronglyPlaced,omitempty"`
}

type DialogueAnswerInput struct {
	Theme   *string                 `json:"theme,omitempty"`
	Set     *string                 `json:"set,omitempty"`
	Content []*DialogueInputContent `json:"content,omitempty"`
}

type DialogueContent struct {
	Translation *string `json:"translation,omitempty"`
	Greek       *string `json:"greek,omitempty"`
	Place       *int32  `json:"place,omitempty"`
	Speaker     *string `json:"speaker,omitempty"`
}

type DialogueCorrection struct {
	Translation  *string `json:"translation,omitempty"`
	Greek        *string `json:"greek,omitempty"`
	Place        *int32  `json:"place,omitempty"`
	Speaker      *string `json:"speaker,omitempty"`
	CorrectPlace *int32  `json:"correctPlace,omitempty"`
}

type DialogueInputContent struct {
	Translation *string `json:"translation,omitempty"`
	Greek       *string `json:"greek,omitempty"`
	Place       *int32  `json:"place,omitempty"`
	Speaker     *string `json:"speaker,omitempty"`
}

type DialogueQuizInput struct {
	Theme *string `json:"theme,omitempty"`
	Set   *string `json:"set,omitempty"`
}

type DialogueQuizResponse struct {
	QuizMetadata *QuizMetadata      `json:"quizMetadata,omitempty"`
	Theme        *string            `json:"theme,omitempty"`
	Set          *string            `json:"set,omitempty"`
	Segment      *string            `json:"segment,omitempty"`
	Reference    *string            `json:"reference,omitempty"`
	Dialogue     *Dialogue          `json:"dialogue,omitempty"`
	Content      []*DialogueContent `json:"content,omitempty"`
}

type GrammarQuizAdded struct {
	CorrectAnswer    *string               `json:"correctAnswer,omitempty"`
	WordInText       *string               `json:"wordInText,omitempty"`
	ExtraInformation *string               `json:"extraInformation,omitempty"`
	Options          []*AuthorBasedOptions `json:"options,omitempty"`
}

type Hit struct {
	Dutch      *string `json:"dutch,omitempty"`
	English    *string `json:"english,omitempty"`
	Greek      *string `json:"greek,omitempty"`
	LinkedWord *string `json:"linkedWord,omitempty"`
	Original   *string `json:"original,omitempty"`
}

type MediaAnswerInput struct {
	Theme         *string `json:"theme,omitempty"`
	Set           *string `json:"set,omitempty"`
	Segment       *string `json:"segment,omitempty"`
	QuizWord      *string `json:"quizWord,omitempty"`
	Answer        *string `json:"answer,omitempty"`
	Comprehensive *bool   `json:"comprehensive,omitempty"`
}

type MediaOptions struct {
	AudioURL *string `json:"audioUrl,omitempty"`
	ImageURL *string `json:"imageUrl,omitempty"`
	Option   *string `json:"option,omitempty"`
}

type MediaQuizInput struct {
	ExcludeWords []*string `json:"excludeWords,omitempty"`
	Theme        *string   `json:"theme,omitempty"`
	Set          *string   `json:"set,omitempty"`
	Segment      *string   `json:"segment,omitempty"`
	Order        *string   `json:"order,omitempty"`
}

type MediaQuizResponse struct {
	NumberOfItems *int32          `json:"numberOfItems,omitempty"`
	Options       []*MediaOptions `json:"options,omitempty"`
	QuizItem      *string         `json:"quizItem,omitempty"`
}

type MultipleChoiceAnswerInput struct {
	Theme         *string `json:"theme,omitempty"`
	Set           *string `json:"set,omitempty"`
	QuizWord      *string `json:"quizWord,omitempty"`
	Answer        *string `json:"answer,omitempty"`
	Comprehensive *bool   `json:"comprehensive,omitempty"`
}

type MultipleChoiceResponse struct {
	NumberOfItems *int32     `json:"numberOfItems,omitempty"`
	Options       []*Options `json:"options,omitempty"`
	QuizItem      *string    `json:"quizItem,omitempty"`
}

type MultipleQuizInput struct {
	ExcludeWords []*string `json:"excludeWords,omitempty"`
	Theme        *string   `json:"theme,omitempty"`
	Set          *string   `json:"set,omitempty"`
	Order        *string   `json:"order,omitempty"`
}

type Options struct {
	Option *string `json:"option,omitempty"`
}

type Query struct {
}

type QuizMetadata struct {
	Language *string `json:"language,omitempty"`
}

type Rhema struct {
	Greek        *string   `json:"greek,omitempty"`
	Section      *string   `json:"section,omitempty"`
	Translations []*string `json:"translations,omitempty"`
}

type Speaker struct {
	Name        *string `json:"name,omitempty"`
	Shorthand   *string `json:"shorthand,omitempty"`
	Translation *string `json:"translation,omitempty"`
}
