package climatiq

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// baseURLPath is a non-empty Client baseURL path to use during tests,
// to ensure relative URLs are used for all endpoints. See go-github
// issue #752. The setupMockClient() functionality has been heavily inspired by
// the go-github library setup() function.
// https://github.com/google/go-github/blob/master/github/github_test.go#L37

const testBaseURL = "/api-v3"

// setupMockClient sets up a test HTTP server along with a climatiq.client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setupMockClient() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(testBaseURL+"/", http.StripPrefix(testBaseURL, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the climatiq client being tested and is
	// configured to use test server.
	client = NewClient(WithBaseURL(server.URL + testBaseURL + "/"))

	return client, mux, server.URL, server.Close
}

// getMockPayload loads one of the json files into a []byte
func getMockPayload(payloadType string) ([]byte, error) {
	fileLocation := fmt.Sprintf("../testdata/%v.json", payloadType)
	file, err := os.ReadFile(fileLocation)
	if err != nil {
		return []byte{}, err
	}

	return file, nil
}

func TestParseSearchRequest(t *testing.T) {
	t.Run("pass: dataVersion is set", func(t *testing.T) {
		_, err := parseSearchRequest(&SearchRequest{DataVersion: "^5"})
		assert.Nil(t, err)
	})

	t.Run("fail: dataVersion not set", func(t *testing.T) {
		_, err := parseSearchRequest(&SearchRequest{})
		expectedErrorMsg := "error: dataVersion must be set"
		assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	})

	t.Run("fail: dataVersion not set multiple params", func(t *testing.T) {
		_, err := parseSearchRequest(&SearchRequest{Category: "cloud compute"})
		expectedErrorMsg := "error: dataVersion must be set"
		assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	})

	t.Run("fail: resultsPerPage exceed maxResultsPerPage", func(t *testing.T) {
		_, err := parseSearchRequest(&SearchRequest{DataVersion: "^0", ResultsPerPage: 105})
		expectedErrorMsg := "error: resultsPerPage cannot be greater than 100"
		assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
	})

	t.Run("pass: parse request with params", func(t *testing.T) {
		sr := SearchRequest{
			DataVersion:    "^5",
			Category:       "cloud computing - cpu",
			ResultsPerPage: 1,
		}

		expectedResult := "category=cloud+computing+-+cpu&data_version=%5E5&results_per_page=1"
		result, err := parseSearchRequest(&sr)
		assert.Nil(t, err)
		assert.Equalf(t, expectedResult, result, "Result should be: %v, got: %v", expectedResult, result)
	})

	t.Run("pass: parse request with multiword query", func(t *testing.T) {
		sr := SearchRequest{
			DataVersion: "^5",
			Query:       "high voltage wind",
		}

		expectedResult := "data_version=%5E5&query=high+voltage+wind"
		result, err := parseSearchRequest(&sr)
		assert.Nil(t, err)
		assert.Equalf(t, expectedResult, result, "Result should be: %v, got: %v", expectedResult, result)
	})

	t.Run("pass: parse request with params and query", func(t *testing.T) {
		sr := SearchRequest{
			DataVersion:    "^5",
			Category:       "cloud computing - cpu",
			ResultsPerPage: 1,
			Query:          "aws",
		}

		expectedResult := "category=cloud+computing+-+cpu&data_version=%5E5&query=aws&results_per_page=1"
		result, err := parseSearchRequest(&sr)
		assert.Nil(t, err)
		assert.Equalf(t, expectedResult, result, "Result should be: %v, got: %v", expectedResult, result)
	})
}

func TestSearchRequest(t *testing.T) {
	t.Run("pass: query with unit types", func(t *testing.T) {
		a := assert.New(t)
		client, mux, _, teardown := setupMockClient()
		defer teardown()

		// The json payload is the response from the following search API query
		// "category=cloud+computing+-+cpu&data_version=%5E5&query=aws&results_per_page=1"
		payload, err := getMockPayload("cloudcomputingCPU_awsquery")
		if err != nil {
			fmt.Printf("error getting mock payload: %s", err)
			return
		}

		mux.HandleFunc(
			"/search",
			func(w http.ResponseWriter, r *http.Request) {
				a.Equal(r.Method, "GET")
				if _, err = w.Write(payload); err != nil {
					fmt.Printf("error writing mock test data: %s", err)
					return
				}
			},
		)

		sr := SearchRequest{
			DataVersion: "^5",
		}

		resp, err := client.Search(context.Background(), &sr)
		a.Nil(err)
		a.Equal(resp.Results[0].Category, "Cloud Computing - CPU")
		a.Equal(resp.Results[0].Factor, float64(0.002196))
		a.Equal(resp.Results[0].Unit, "kg/CPU-hour")
		a.Equal(resp.Results[0].Name, "AWS (af-south-1) CPU")
	})
}
