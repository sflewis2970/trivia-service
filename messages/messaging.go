package messages

const MAKE_SELECTION_MSG string = "Make Selection from list..."

const (
	DASH             string = "-"
	COMMA            string = ","
	SPACE            string = " "
	ONE_SET          int    = 1
	FIND_ERROR       int    = -1
	EMPTY_QUESTIONID int    = 0
	EMPTY_QUESTION   int    = 0
	EMPTY_ANSWER     int    = 0
	EMPTY_CATEGORY   int    = 0
	EMPTY_CHOICES    int    = 0
	EMPTY_TIMESTAMP  int    = 0
	EMPTY_WARNING    int    = 0
)

// Client-server communication
type QuestionResponse struct {
	QuestionID string   `json:"questionid"`
	Question   string   `json:"question"`
	Category   string   `json:"category"`
	Choices    []string `json:"choices"`
	Timestamp  string   `json:"timestamp"`
	Warning    string   `json:"warning,omitempty"`
	Error      string   `json:"error,omitempty"`
}

type AnswerRequest struct {
	QuestionID string `json:"questionid"`
	Response   string `json:"response"`
}

type AnswerResponse struct {
	Question  string `json:"question"`
	Timestamp string `json:"timestamp"`
	Category  string `json:"category"`
	Response  string `json:"response"`
	Answer    string `json:"answer"`
	Correct   bool   `json:"correct"`
	Message   string `json:"message,omitempty"`
	Warning   string `json:"warning,omitempty"`
	Error     string `json:"error,omitempty"`
}
