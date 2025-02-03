package runpod

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

type GQLQueryWrapper struct {
	Query string `json:"query"`
}

type GQLResponse struct {
	Data   map[string]any `json:"data"`
	Errors []GQLError     `json:"errors"`
}

type GQLError struct {
	Message   string           `json:"message"`
	Locations []map[string]any `json:"locations"`
	Path      []any            `json:"path"`
}

func (e *GQLError) Error() string {
	msg := `Runpod GQL query error: %s`
	return fmt.Sprintf(msg, e.Message)
}

var ErrUnauthorized = errors.New("unauthorized request, please check your API key")

type ErrFailedQuery struct {
	Msg string
}

func (e *ErrFailedQuery) Error() string {
	return fmt.Sprintf("query to Runpod failed; %s", e.Msg)
}

var RunpodGQLEndpoint string
var RunpodUserAgent string

func constructUserAgent() string {
	osInfo := fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
	goVersion := runtime.Version()
	integrationMethod := os.Getenv("RUNPOD_UA_INTEGRATION")

	uaComponents := []string{
		fmt.Sprintf("RunPod-Go-SDK/%s", SdKVersion),
		fmt.Sprintf("(%s)", osInfo),
		fmt.Sprintf("Language/Go %s", goVersion),
	}

	if integrationMethod != "" {
		uaComponents = append(uaComponents, fmt.Sprintf("Integration/%s", integrationMethod))
	}

	userAgent := fmt.Sprintf("%s", uaComponents)
	return userAgent
}

func init() {

	RunpodUserAgent = constructUserAgent()

	runpodBaseUrl := os.Getenv("RUNPOD_API_BASE_URL")
	if runpodBaseUrl == "" {
		RunpodGQLEndpoint = "https://api.runpod.io/graphql"
		return
	}
	RunpodGQLEndpoint = fmt.Sprintf("%s/graphql", runpodBaseUrl)
}

func (c *Client) query(GQLQuery string) (*GQLResponse, error) {

	client := &http.Client{Timeout: 30 * time.Second}

	query := GQLQueryWrapper{
		Query: GQLQuery,
	}

	body, err := json.Marshal(query)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s?api_key=%s", RunpodGQLEndpoint, c.ApiKey)

	request, err := http.NewRequest(
		"POST",
		url,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", RunpodUserAgent)

	// this is the most likely error...
	response, err := client.Do(request)
	if err != nil {
		return nil, &ErrFailedQuery{Msg: err.Error()}
	}
	defer response.Body.Close()
	// ...together with this

	// ❗️❗️❗️❗️❗️BUG gql endpoint: unauthorized requests do not return status 401, but 200
	if response.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	respBody, err := io.ReadAll(response.Body)
	// unlikely
	if err != nil {
		return nil, err
	}

	var respWrapper GQLResponse
	// unlikely, unless problems in data formatting
	err = json.Unmarshal(respBody, &respWrapper)
	if err != nil {
		return nil, err
	}

	if respWrapper.Errors != nil && len(respWrapper.Errors) != 0 {
		// here i follow the python sdk implementation
		// that raises the first error in the errors list
		// assuming one error at the time is usually thrown
		//
		// eventually respWrapper can be inspected to find other errors
		return &respWrapper, &respWrapper.Errors[0]
	}

	return &respWrapper, nil
}
