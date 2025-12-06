package httpcaller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrNotFound            = fmt.Errorf("not-found")
	ErrInternalServerError = fmt.Errorf("internal-server-error")
	ErrUnauthorizedError   = fmt.Errorf("unauthorized-error")
)

type File struct {
	Content     io.Reader
	ContentType string
}

type requestOption struct {
	Method      string
	RequestBody interface{}
	URL         string
	Headers     map[string]string
	ResponesMap map[int]interface{}
	Client      Client
}

func buildPayloadReader(headers map[string]string, payload interface{}) (io.Reader, string, error) {
	file, ok := payload.(File)
	if ok {
		return file.Content, file.ContentType, nil
	}
	byteContent, okByte := payload.(io.Reader)
	if contentType, ok := headers["Content-Type"]; ok && okByte {
		return byteContent, contentType, nil
	}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(payload)
	return buf, "application/json", err
}

func doRequest(ctx context.Context, options *requestOption) (int, *http.Response, error) {
	if options == nil {
		return 0, nil, fmt.Errorf("attempt to resquet with nil option : %w", ErrInternalServerError)
	}
	body, ctIn, err := buildPayloadReader(options.Headers, options.RequestBody)
	if err != nil {
		return 0, nil, fmt.Errorf("%v : %w", err, ErrInternalServerError)
	}
	request, err := http.NewRequestWithContext(ctx, options.Method, options.URL, body)
	if err != nil {
		select {
		case <-ctx.Done():
			return 0, nil, ctx.Err()
		default:
		}
		return 0, nil, fmt.Errorf("%v : %w", err, ErrInternalServerError)
	}

	for k, v := range options.Headers {
		request.Header.Add(k, v)
	}
	if ctIn != "" {
		request.Header.Del("Content-Type")
		request.Header.Add("Content-Type", ctIn)
	}

	response, err := options.Client.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("%v : %w", err, ErrInternalServerError)
	}

	decodedBody, ok := options.ResponesMap[response.StatusCode]
	if !ok {
		return handleStatsCodeNotInResponseMap(options, response, request)
	}
	if options.ResponesMap[response.StatusCode] == nil {
		return response.StatusCode, response, nil
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(decodedBody)
	if err != nil {
		return response.StatusCode, response, fmt.Errorf("%w : while decoding response of %s (%d) %s", err, options.Method, response.StatusCode, options.URL)
	}
	return response.StatusCode, response, nil
}

func handleStatsCodeNotInResponseMap(options *requestOption, response *http.Response, request *http.Request) (int, *http.Response, error) {
	defer response.Body.Close()
	rspBody, _ := ioutil.ReadAll(response.Body)
	expectedStatusCodes := []int{}
	for st := range options.ResponesMap {
		expectedStatusCodes = append(expectedStatusCodes, st)
	}
	err := fmt.Errorf("wrong status code received on (%s) Expected %v, got %d\n%s", request.URL.Path, expectedStatusCodes, response.StatusCode, string(rspBody))
	return response.StatusCode, nil, fmt.Errorf("%v : %w", err, detectErrorType(response.StatusCode))
}

func detectErrorType(statusCode int) error {
	switch statusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorizedError
	default:
		return ErrInternalServerError
	}
}
