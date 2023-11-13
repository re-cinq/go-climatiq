package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go-climatiq/climatiq"
)

func main() {
	// get a new climatiq http client and pass the auth token
	// option with your climatiq API KEY
	// https://www.climatiq.io/docs/api-reference/authentication
	cli := climatiq.NewClient(
		climatiq.WithAuthToken("API_KEY"),
	)

	// create a SearchRequest instance from the SearchRequest type
	// and add parameters related to your query
	// https://www.climatiq.io/docs/api-reference/search#request
	searchRequest := climatiq.SearchRequest{
		// Note: dataVersion is a required field
		// https://www.climatiq.io/docs/api-reference/data-version-endpoint
		DataVersion: "^5",
		Category:    "cloud computing - cpu",
	}

	// client.Search runs an http GET request with the search request
	// data, checks the response, and unmarshalls the response body
	// into the SearchResponse and the embedded SearchResults structs
	resp, err := cli.Search(context.Background(), &searchRequest)
	if err != nil {
		// do error handling
		fmt.Printf("error with API search get: %s", err)
		return
	}

	// pretty print the JSON response
	val, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		fmt.Printf("error pretty printing json response: %s", err)
		return
	}

	fmt.Println(string(val))
}
