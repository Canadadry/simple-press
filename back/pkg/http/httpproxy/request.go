package httpproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrInternalServerError = fmt.Errorf("internal-server-error")
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
	Client      Client
}

func buildPayloadReader(payload interface{}) (io.Reader, string, error) {
	file, ok := payload.(File)
	if ok {
		return file.Content, file.ContentType, nil
	}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(payload)
	return buf, "application/json", err
}

func doRequest(ctx context.Context, options *requestOption) (*http.Response, error) {
	if options == nil {
		return nil, fmt.Errorf("attempt to resquet with nil option : %w", ErrInternalServerError)
	}
	body, ctIn, err := buildPayloadReader(options.RequestBody)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", err, ErrInternalServerError)
	}
	request, err := http.NewRequestWithContext(ctx, options.Method, options.URL, body)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, fmt.Errorf("%v : %w", err, ErrInternalServerError)
	}

	if ctIn != "" {
		request.Header.Add("Content-Type", ctIn)
	}
	for k, v := range options.Headers {
		request.Header.Add(k, v)
	}

	response, err := options.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", err, ErrInternalServerError)
	}

	return response, nil
}
