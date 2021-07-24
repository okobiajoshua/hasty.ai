package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hastyai/backend/datastore"
	"github.com/hastyai/backend/handler"
	"github.com/stretchr/testify/suite"
)

type e2eTestSuite struct {
	suite.Suite
}

var jobID string

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) BeforeTest(suiteName, testName string) {
	createJob("12345")
}

func (s *e2eTestSuite) Test_EndToEnd_PostJob() {
	reqStr := `{"object_id":"123457890"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/api/v1/job", 8000), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set("content-type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	var res handler.ResponseDTO
	err = json.NewDecoder(response.Body).Decode(&res)
	s.NoError(err)

	jobID = res.JobID // Set global value for other test to use

	s.NotEmpty(res.JobID, "response should contain a job ID")
	response.Body.Close()
}

func (s *e2eTestSuite) Test_EndToEnd_GetJobStatus() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/api/v1/job/"+jobID, 8000), nil)
	s.NoError(err)

	req.Header.Set("content-type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	var res datastore.Job
	err = json.NewDecoder(response.Body).Decode(&res)
	s.NoError(err)

	s.Equal(res.JobID, jobID)
	s.Equal(res.ObjectID, "12345")
	s.NotEmpty(res.Status)
	response.Body.Close()
}

func createJob(objectID string) {
	reqStr := `{"object_id":"` + objectID + `"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/api/v1/job", 8000), strings.NewReader(reqStr))
	if err != nil {
		panic(err)
	}

	req.Header.Set("content-type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var res handler.ResponseDTO
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		panic(err)
	}

	jobID = res.JobID // Set global value for other test to use
	response.Body.Close()
}
