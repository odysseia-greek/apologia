package kunismos

type GrammarBasedQuiz struct {
	Theme           string           `json:"theme"`
	Set             int              `json:"set"`
	Segment         string           `json:"segment"`
	Description     string           `json:"description"`
	Difficulty      string           `json:"difficulty"`
	Content         []GrammarContent `json:"content"`
	ContractionRule string           `json:"contractionRule,omitempty"`
}
type GrammarContent struct {
	Greek           string `json:"greek"`
	DictionaryForm  string `json:"dictionaryForm"`
	Translation     string `json:"translation"`
	Stem            string `json:"stem"`
	GrammarQuestion struct {
		Tense         string `json:"tense"`
		Voice         string `json:"voice"`
		Mood          string `json:"mood"`
		Person        string `json:"person"`
		Number        string `json:"number"`
		CorrectAnswer string `json:"correctAnswer"`
	} `json:"grammarQuestion"`
}
