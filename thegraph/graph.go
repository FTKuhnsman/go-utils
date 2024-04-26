package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/FTKuhnsman/go-utils/common"
)

var apiKey string
var subgraphID string

func init() {
	apiKey = common.GetStringEnvWithDefault("THEGRAPH_API_KEY", "")
	subgraphID = common.GetStringEnvWithDefault("THEGRAPH_SUBGRAPH_ID", "")
}

type graphRequest struct {
	Query string `json:"query"`
}

type graphResponse struct {
	Data map[string]interface{} `json:"data"`
}

func QueryTheGraph(query string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://gateway-arbitrum.network.thegraph.com/api/%s/subgraphs/id/%s", apiKey, subgraphID)

	data := graphRequest{
		Query: query,
	}

	// Marshal the data into a JSON byte slice
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Create a new HTTP request with the JSON data
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Set the content type to application/json
	request.Header.Set("Content-Type", "application/json")

	// Send the request using the http.DefaultClient
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer response.Body.Close()

	// Read and print the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var jsonResponse graphResponse

	err = json.Unmarshal(responseBody, &jsonResponse)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return jsonResponse.Data, nil
}
