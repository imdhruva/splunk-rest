// Package splunk enabled to perform restful calls to splunk instances
package splunk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// WaitTime to be used for different polling of restAPIs
const (
	WaitTimeSec  = 5
	RetryAttempt = 10
)

// Search returns the results of the search string in Splunk
func (splunkURL URL) Search(query string, user User) (string, error) {
	query = strings.TrimSpace(query)
	// check if the query begins with `search` or | (pipe) keyword
	if !strings.HasPrefix(query, "search") && !strings.HasPrefix(query, "|") {
		query = "search " + query
	}

	// create a payload to be added to the post request
	var payload io.Reader
	val := url.Values{}
	val.Add("output_mode", splunkURL.OutputFormat())
	searchEndPoint, err := splunkURL.SearchEndPoint()
	if err != nil {
		return "", err
	}
	val.Add("search", query)
	if val != nil {
		payload = bytes.NewBufferString(val.Encode())
	}
	responseByte, err := user.HttpCall("POST", searchEndPoint, payload)
	if err != nil {
		return "", err
	}

	// extracting json field sid
	var sid SearchID
	unmarshallErr := json.Unmarshal(responseByte, &sid)
	if unmarshallErr != nil {
		return "", unmarshallErr
	}
	
	return sid.Value, nil
}

// GetJobEndpoint returns the rest endpoint for the respective JobID
func (splunkURL URL) GetJobEndpoint(jobID string) (string, error) {
	searchEndPoint, err := splunkURL.SearchEndPoint()
	if err != nil {
		return "", err
	}
	return (searchEndPoint + "/" + jobID), nil
}

// GetJobStat fetches result set of the search-job as reprented by rest
// endpoint /services/search/jobs/{jobID}
func (user User) GetJobStat(jobID string, splunkURL URL) ([]byte, error) {
	// create a payload
	val := url.Values{}
	var payload io.Reader
	val.Add("output_mode", splunkURL.OutputFormat())
	jobEndPoint, err := splunkURL.GetJobEndpoint(jobID)
	if err != nil {
		return nil, err
	}
	if val != nil {
		payload = bytes.NewBufferString(val.Encode())
	}
	responseByte, err := user.HttpCall("POST", jobEndPoint, payload)
	if err != nil {
		return nil, err
	}

	return responseByte, nil

	// if dispatchState != "DONE" {
	// 	time.Sleep(5 * 1000 * time.Millisecond)
	// 	return user.GetSearchJob(jobID, splunkURL)
	// }

}

// GetJobStatus returns current dispatchState of the current job for the specific sid
func GetJobStatus(responseByte []byte) (string, error) {
	var jobStat SearchJob
	unmarshallErr := json.Unmarshal(responseByte, &jobStat)
	if unmarshallErr != nil {
		return "", unmarshallErr
	}
	dispatchState := jobStat.Entry[0].Content.DispatchState
	return dispatchState, nil
}

// IsJobComplete returns boolean to identify if the JobId's dispatchState is "DONE"
func (splunkURL URL) IsJobComplete(jobID string, user User) (bool, error) {
	//  get stats for jobID
	jobStat, err := user.GetJobStat(jobID, splunkURL)
	if err != nil {
		return false, err
	}

	// get the dispatchState of the jobId
	var dispatchStatus string
	dispatchStatus, err = GetJobStatus(jobStat)
	if err != nil {
		return false, err
	}
	if dispatchStatus == "DONE" {
		return true, nil
	}

	return false, nil
}

// GetSearchResult returns resultset of the search query as represented by rest
// end-point /services/search/jobs/{jobID}/results
func (splunkURL URL) GetSearchResult(jobID string, user User) ([]byte, error) {
	val := url.Values{}
	var payload io.Reader
	val.Add("output_mode", splunkURL.OutputFormat())
	payload = bytes.NewBufferString(val.Encode())
	jobEndPoint, err := splunkURL.GetJobEndpoint(jobID)
	if err != nil {
		return nil, err
	}

	// keep polling the jobID status every 5 seconds to check the status is DONE
	tempCount := 0
	// jobComplete, err := user.IsJobComplete(jobID, splunkURL)
	jobComplete := false

	if err != nil {
		return nil, err
	}
	for ; !(jobComplete) && (tempCount < RetryAttempt); tempCount++ {
		time.Sleep(WaitTimeSec * 1000 * time.Millisecond)
		jobComplete, err = splunkURL.IsJobComplete(jobID, user)
		if err != nil {
			return nil, err
		}

	}
	if tempCount == (RetryAttempt - 1) {
		err = errors.New("Exceeded retry attempts : " + strconv.Itoa(RetryAttempt))
		return nil, err
	}
	// fetch job result
	jobResultResponse, err := user.HttpCall("GET", jobEndPoint+"/results/", payload)
	if err != nil {
		return nil, err
	}

	return jobResultResponse, nil
}

// SearchID is the json representation of json elements in /services/search/jobs
type SearchID struct {
	Value string `json:"sid"`
}

// SearchJob is the struct representation of the json respese received from
// /services/search/jobs/{sid}/
// generated via https://mholt.github.io/json-to-go/
type SearchJob struct {
	Links struct {
	} `json:"links"`
	Origin    string    `json:"origin"`
	Updated   time.Time `json:"updated"`
	Generator struct {
		Build   string `json:"build"`
		Version string `json:"version"`
	} `json:"generator"`
	Entry []struct {
		Name    string    `json:"name"`
		ID      string    `json:"id"`
		Updated time.Time `json:"updated"`
		Links   struct {
			Alternate      string `json:"alternate"`
			SearchLog      string `json:"search.log"`
			Events         string `json:"events"`
			Results        string `json:"results"`
			ResultsPreview string `json:"results_preview"`
			Timeline       string `json:"timeline"`
			Summary        string `json:"summary"`
			Control        string `json:"control"`
		} `json:"links"`
		Published time.Time `json:"published"`
		Author    string    `json:"author"`
		Content   struct {
			CursorTime     time.Time     `json:"cursorTime"`
			DefaultSaveTTL string        `json:"defaultSaveTTL"`
			DefaultTTL     string        `json:"defaultTTL"`
			Delegate       string        `json:"delegate"`
			DiskUsage      int           `json:"diskUsage"`
			DispatchState  string        `json:"dispatchState"`
			DoneProgress   int           `json:"doneProgress"`
			EarliestTime   time.Time     `json:"earliestTime"`
			IsDone         bool          `json:"isDone"`
			IsFailed       bool          `json:"isFailed"`
			IsFinalized    bool          `json:"isFinalized"`
			IsPaused       bool          `json:"isPaused"`
			IsSaved        bool          `json:"isSaved"`
			IsSavedSearch  bool          `json:"isSavedSearch"`
			IsZombie       bool          `json:"isZombie"`
			Label          string        `json:"label"`
			Sid            string        `json:"sid"`
			StatusBuckets  int           `json:"statusBuckets"`
			TTL            int           `json:"ttl"`
			Messages       []interface{} `json:"messages"`
			Request        struct {
				Search string `json:"search"`
			} `json:"request"`
			Runtime struct {
				AutoCancel string `json:"auto_cancel"`
				AutoPause  string `json:"auto_pause"`
			} `json:"runtime"`
			SearchProviders []interface{} `json:"searchProviders"`
		} `json:"content"`
		Acl struct {
			Perms struct {
				Read  []string `json:"read"`
				Write []string `json:"write"`
			} `json:"perms"`
			Owner      string `json:"owner"`
			Modifiable bool   `json:"modifiable"`
			Sharing    string `json:"sharing"`
			App        string `json:"app"`
			CanWrite   bool   `json:"can_write"`
			TTL        string `json:"ttl"`
		} `json:"acl"`
	} `json:"entry"`
	Paging struct {
		Total   int `json:"total"`
		PerPage int `json:"perPage"`
		Offset  int `json:"offset"`
	} `json:"paging"`
}
