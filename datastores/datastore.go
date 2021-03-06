package datastores

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sflewis2970/trivia-service/config"
)

const (
	NBR_OF_ATTEMPTS int = 3
	DELAY_IN_SECS   int = 6
)

// Datastore struct
type DataStore struct {
	cfgData      *config.ConfigData
	serverStatus StatusCode
}

// AddQuestionAndAnswer sends a request to the DataStore server to add a question to the datastore
func (ds *DataStore) Insert(questionID string, dst DataStoreTable) error {
	// Create AddQuestionRequest message
	var aqRequest AddQuestionRequest

	// Build AddQuestionRequest
	aqRequest.QuestionID = questionID
	aqRequest.Question = dst.Question
	aqRequest.Category = dst.Category
	aqRequest.Answer = dst.Answer

	// Build request body from AddQuestionRequest
	requestBody, marshalErr := json.Marshal(aqRequest)
	if marshalErr != nil {
		log.Print("marshaling error: ", marshalErr)
		return marshalErr
	}

	// Post request
	var url string
	if ds.cfgData.Env == config.PRODUCTION {
		url = ds.cfgData.DataStoreName + DS_INSERT_PATH
	} else {
		url = ds.cfgData.DataStoreName + ":" + ds.cfgData.DataStorePort + DS_INSERT_PATH
	}

	response, postErr := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if postErr != nil {
		return postErr
	}
	defer response.Body.Close()

	// Handle add question response
	var aqResponse AddQuestionResponse

	// Read response stream into JSON
	json.NewDecoder(response.Body).Decode(&aqResponse)

	return nil
}

// CheckAnswer sends a request to the DataStore server to determine if the question was answered correctly
func (ds *DataStore) Get(questionID string) (string, *QuestionAndAnswer, error) {
	timestamp := ""
	var caRequest CheckAnswerRequest

	// Build add question request
	caRequest.QuestionID = questionID

	// Convert struct to byte array
	requestBody, marshalErr := json.Marshal(caRequest)
	if marshalErr != nil {
		log.Print("marshaling error: ", marshalErr)
		return "", nil, marshalErr
	}

	// Create a http request
	var url string
	if ds.cfgData.Env == config.PRODUCTION {
		url = ds.cfgData.DataStoreName + DS_GET_PATH
	} else {
		url = ds.cfgData.DataStoreName + ":" + ds.cfgData.DataStorePort + DS_GET_PATH
	}

	response, postErr := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if postErr != nil {
		return "", nil, postErr
	}
	defer response.Body.Close()

	// Handle add question response
	var caResponse CheckAnswerResponse

	// Read response stream into JSON
	json.NewDecoder(response.Body).Decode(&caResponse)

	// Update QuestionAndAnswer struct
	timestamp = caResponse.Timestamp

	var newQA *QuestionAndAnswer
	newQA = new(QuestionAndAnswer)
	newQA.Question = caResponse.Question
	newQA.Category = caResponse.Category
	newQA.Answer = caResponse.Answer
	newQA.Message = caResponse.Message
	newQA.Warning = caResponse.Warning
	newQA.Error = caResponse.Error

	return timestamp, newQA, nil
}

// SendStatusRequest sends a request for the status of the datastore server
func (ds *DataStore) sendStatusRequest() StatusCode {
	var url string
	if ds.cfgData.Env == config.PRODUCTION {
		url = ds.cfgData.DataStoreName + DS_STATUS_PATH
	} else {
		url = ds.cfgData.DataStoreName + ":" + ds.cfgData.DataStorePort + DS_STATUS_PATH
	}

	log.Print("sending status request to ", url)

	// http request
	response, getErr := http.Get(url)
	if getErr != nil {
		log.Print("A response error has occurred...")
		return StatusCode(DS_REQUEST_ERROR)
	}
	defer response.Body.Close()

	// Status (Request) Response
	var sResponse StatusResponse

	// Read JSON from stream
	json.NewDecoder(response.Body).Decode(&sResponse)

	return sResponse.Status
}

// CreateDataStore prepares the datastore component waits for the datastore server before allowing messages to be sent
func New() *DataStore {
	// Create QuestionDataStore object
	log.Print("Creating DataStore object")
	ds := new(DataStore)

	// Update QuestionDataStore fields
	var cfgDataErr error
	ds.cfgData, cfgDataErr = config.Get().GetData()
	if cfgDataErr != nil {
		log.Print("error getting config data...")
		return nil
	}

	// Initialize status code
	ds.serverStatus = StatusCode(DS_NOT_STARTED)

	// Wait for DataStore server to become available
	nbrOfAttempts := NBR_OF_ATTEMPTS
	for ds.serverStatus != StatusCode(DS_RUNNING) || nbrOfAttempts > 0 {
		// Get datastore server status
		ds.serverStatus = ds.sendStatusRequest()

		// Once the datastore is up and running get out!
		if ds.serverStatus == StatusCode(DS_RUNNING) {
			break
		} else {
			log.Print("waiting for Datastore server...")
		}

		// Decrement retry counter
		nbrOfAttempts--

		// Sleep for 3 seconds
		time.Sleep(time.Second * time.Duration(DELAY_IN_SECS))
	}

	return ds
}
