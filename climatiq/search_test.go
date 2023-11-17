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

// setupMockClient sets up a test HTTP server along with a climatiq.client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
// The setupMockClient() functionality has been heavily inspired by
// the go-github library setup() function.
func setupMockClient() (client *Client, mux *http.ServeMux, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", mux)

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the climatiq client being tested and is
	// configured to use test server.
	client = NewClient(WithBaseURL(server.URL + "/"))

	return client, mux, server.Close
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
	tests := []struct {
		name     string
		sr       SearchRequest
		hasError bool
		expRes   string
		expErr   string
	}{
		{
			name:     "pass: dataVersion is set",
			sr:       SearchRequest{DataVersion: "^5"},
			hasError: false,
			expRes:   "data_version=%5E5",
			expErr:   "",
		},
		{
			name:     "fail: dataVersion not set",
			sr:       SearchRequest{},
			hasError: true,
			expRes:   "",
			expErr:   "error: dataVersion must be set",
		},
		{
			name: "fail: dataVersion not set multiple params",
			sr: SearchRequest{
				Category: "cloud compute",
				Region:   "Switzerland",
			},
			hasError: true,
			expRes:   "",
			expErr:   "error: dataVersion must be set",
		},
		{
			name: "fail: dataVersion not set multiple params",
			sr: SearchRequest{
				Category: "cloud compute",
				Region:   "Switzerland",
			},
			hasError: true,
			expRes:   "",
			expErr:   "error: dataVersion must be set",
		},
		{
			name: "pass: parse request with params",
			sr: SearchRequest{
				DataVersion:    "^5",
				Category:       "cloud computing - cpu",
				ResultsPerPage: 1,
			},

			hasError: false,
			expRes:   "category=cloud+computing+-+cpu&data_version=%5E5&results_per_page=1",
			expErr:   "",
		},
		{
			name: "pass: parse request with multiword query",
			sr: SearchRequest{
				DataVersion: "^5",
				Query:       "high voltage wind",
			},
			hasError: false,
			expRes:   "data_version=%5E5&query=high+voltage+wind",
			expErr:   "",
		},
		{
			name: "pass: parse request with params and query",
			sr: SearchRequest{
				DataVersion:    "^5",
				Category:       "cloud computing - cpu",
				ResultsPerPage: 1,
				Query:          "aws",
			},
			hasError: false,
			expRes:   "category=cloud+computing+-+cpu&data_version=%5E5&query=aws&results_per_page=1",
			expErr:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := parseSearchRequest(&test.sr)
			assert.Equalf(t, res, test.expRes, "Result should be: %v, got: %v", test.expRes, res)
			if test.hasError {
				assert.EqualErrorf(t, err, test.expErr, "Error should be: %v, got: %v", test.expErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSearchRequest(t *testing.T) {
	t.Run("pass: query with unit types", func(t *testing.T) {
		a := assert.New(t)
		client, mux, teardown := setupMockClient()
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
