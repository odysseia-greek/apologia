package strategos

type JourneyBasedQuiz struct {
	Theme       string `json:"theme"`
	Segment     string `json:"segment"`
	Number      int    `json:"number"`
	Location    string `json:"location"`
	Coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"coordinates"`
	FullSentence string `json:"fullSentence,omitempty"`
	Translation  string `json:"translation"`
	FixedSteps   []struct {
		Step    string `json:"step"`
		Type    string `json:"type"`
		Content struct {
			Author     string `json:"author"`
			Work       string `json:"work"`
			Background string `json:"background"`
		} `json:"content"`
	} `json:"fixedSteps"`
	RandomSteps []RandomSteps `json:"randomSteps"`
	FinalStep   struct {
		Type        string   `json:"type"`
		Instruction string   `json:"instruction"`
		Options     []string `json:"options"`
		Answer      string   `json:"answer"`
	} `json:"finalStep"`
	ContextNote struct {
		Text string `json:"text"`
	} `json:"contextNote"`
}

type RandomSteps struct {
	Type        string `json:"type"`
	Instruction string `json:"instruction,omitempty"`
	Pairs       []struct {
		Greek  string `json:"greek"`
		Answer string `json:"answer"`
	} `json:"pairs,omitempty"`
	Question   string   `json:"question,omitempty"`
	Options    []string `json:"options,omitempty"`
	Answer     string   `json:"answer,omitempty"`
	MediaFiles []struct {
		Word   string `json:"word"`
		Answer string `json:"answer"`
	} `json:"mediaFiles,omitempty"`
	Verbs []struct {
		Word   string `json:"word"`
		Answer string `json:"answer"`
	} `json:"verbs,omitempty"`
	Title         string `json:"title,omitempty"`
	Text          string `json:"text,omitempty"`
	NoteOnCorrect string `json:"noteOnCorrect,omitempty"`
	Note          string `json:"note,omitempty"`
}
